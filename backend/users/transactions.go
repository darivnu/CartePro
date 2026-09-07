//
// EPITECH PROJECT, 2026
// CartePro
// File description:
// transactions
//

package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"cartepro/database"
	"cartepro/server"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type QrValidationRequest struct {
	QrPayload      string `json:"qr_payload"`
	Amount         int64  `json:"amount"`
	IdempotencyKey string `json:"idempotency_key"`
}

var (
	errQrTokenUnavailable = errors.New("qr token not found, expired, or already used")
	errInsufficientBalance = errors.New("insufficient balance")
)

// isUniqueViolation reports whether err is a Postgres unique-constraint failure.
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505" //code for unique violation
}

func HandleQrCodeValidation(w http.ResponseWriter, r *http.Request) {
	_, partner, err := server.GetPartnerFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if partner.Status != database.StatusApproved {
		http.Error(w, "Partner is not approved", http.StatusForbidden)
		return
	}

	var req QrValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount <= 0 || req.IdempotencyKey == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, expiresAt, err := verifyQrPayload(req.QrPayload)
	if err != nil || time.Now().After(expiresAt) {
		http.Error(w, "Invalid or expired QR code", http.StatusBadRequest)
		return
	}

	// everything in debitClientForQr runs in one DB transaction: the QR-token claim,
	// the balance move, and the ledger insert either all land or all roll back
	// together, so a failure partway through (insufficient balance, a raced
	// idempotency key, ...) can never leave the token claimed but the money
	// unmoved, or vice versa
	var transaction database.Transaction
	txErr := database.DB.Transaction(func(tx *gorm.DB) error {
		var err error
		transaction, err = debitClientForQr(tx, token, req, partner)
		return err
	})

	switch {
	case txErr == nil:
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"transaction": transaction})
	case errors.Is(txErr, errQrTokenUnavailable):
		http.Error(w, "QR code not found, expired, or already used", http.StatusConflict)
	case errors.Is(txErr, errInsufficientBalance):
		http.Error(w, "Insufficient balance", http.StatusPaymentRequired)
	case isUniqueViolation(txErr):
		// a retried request with the same idempotency key replays the original
		// result instead of being processed twice
		var existing database.Transaction
		if err := database.DB.Where("idempotency_key = ?", req.IdempotencyKey).First(&existing).Error; err != nil {
			http.Error(w, "Failed to process transaction", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"transaction": existing})
	default:
		http.Error(w, "Failed to process transaction", http.StatusInternalServerError)
	}
}

// debitClientForQr claims the QR token, moves the balance from client to partner,
// and inserts the ledger row, all against tx so the caller can run it inside a
// single DB transaction.
func debitClientForQr(tx *gorm.DB, token string, req QrValidationRequest, partner *database.Partner) (database.Transaction, error) {
	// atomically claim the token: this single conditional UPDATE is the whole
	// concurrency guard, only one concurrent request can match "used_at IS NULL"
	// for a given row, so a code can never be redeemed twice
	now := time.Now()
	claim := tx.Model(&database.QrToken{}).
		Where("token = ? AND used_at IS NULL AND expires_at > ?", token, now).
		Update("used_at", now)
	if claim.Error != nil {
		return database.Transaction{}, claim.Error
	}
	if claim.RowsAffected == 0 {
		return database.Transaction{}, errQrTokenUnavailable
	}

	var qrToken database.QrToken
	if err := tx.Where("token = ?", token).First(&qrToken).Error; err != nil {
		return database.Transaction{}, err
	}

	var client database.Client
	if err := tx.First(&client, qrToken.ClientID).Error; err != nil {
		return database.Transaction{}, err
	}

	// atomic conditional debit: the WHERE clause re-checks the balance at write
	// time, so two concurrent requests can't both read a stale balance and both
	// succeed (the same pattern as the QR-token claim above)
	debit := tx.Model(&database.Client{}).
		Where("id = ? AND balance >= ?", client.ID, req.Amount).
		Update("balance", gorm.Expr("balance - ?", req.Amount))
	if debit.Error != nil {
		return database.Transaction{}, debit.Error
	}
	if debit.RowsAffected == 0 {
		return database.Transaction{}, errInsufficientBalance
	}

	if err := tx.Model(&database.Partner{}).
		Where("id = ?", partner.ID).
		Update("balance", gorm.Expr("balance + ?", req.Amount)).Error; err != nil {
		return database.Transaction{}, err
	}

	transaction := database.Transaction{
		SenderUserID:   client.UserID,
		ReceiverUserID: partner.UserID,
		QrTokenID:      &qrToken.ID,
		Amount:         req.Amount,
		Type:           database.TransactionTypeDebit,
		IdempotencyKey: &req.IdempotencyKey,
	}
	// the unique constraint on idempotency_key is what actually makes this
	// atomic: a raced retry with the same key hits this insert and fails here,
	// inside the transaction, instead of a separate check-then-act query
	if err := tx.Create(&transaction).Error; err != nil {
		return database.Transaction{}, err
	}
	return transaction, nil
}

func HandleGetClientOwnTransactions(w http.ResponseWriter, r *http.Request) {
	_, client, err := server.GetClientFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	query := r.URL.Query()

	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 1 {
		page = 1 //the page is which chunk we want (defined by limit), so if we want the first chunk, we set page to 1)
	}

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	from := query.Get("from")
	to := query.Get("to")

	db := database.DB.Model(&database.Transaction{}).Where("client_id = ?", client.ID)
	if from != "" {
		fromTime, err := time.Parse(time.RFC3339, from)
		if err != nil {
			http.Error(w, "Invalid 'from' date format", http.StatusBadRequest)
			return
		}
		db = db.Where("created_at >= ?", fromTime)
	}
	if to != "" {
		toTime, err := time.Parse(time.RFC3339, to)
		if err != nil {
			http.Error(w, "Invalid 'to' date format", http.StatusBadRequest)
			return
		}
		db = db.Where("created_at <= ?", toTime)
	}

	var total int64
	db.Count(&total)

	var transactions []database.Transaction
	if err := db.Offset((page - 1) * limit).Limit(limit).Find(&transactions).Error; err != nil {
		http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": transactions,
		"meta": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
			"from":        from,
			"to":          to,
		},
	})
}
