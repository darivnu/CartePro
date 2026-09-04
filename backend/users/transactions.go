//
// EPITECH PROJECT, 2026
// CartePro
// File description:
// transactions
//

package users

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"cartepro/database"
	"cartepro/server"
)

type QrValidationRequest struct {
	QrPayload      string `json:"qr_payload"`
	Amount         int64  `json:"amount"`
	IdempotencyKey string `json:"idempotency_key"`
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

	//idempotency: a retried request with the same key replays the original
	//result instead of being processed twice
	//this is because the client may retry a request if they don't get a response, and we don't want to charge them twice
	var existing database.Transaction
	if err := database.DB.Where("idempotency_key = ?", req.IdempotencyKey).First(&existing).Error; err == nil {
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"transaction": existing})
		return
	}

	token, expiresAt, err := verifyQrPayload(req.QrPayload)
	if err != nil || time.Now().After(expiresAt) {
		http.Error(w, "Invalid or expired QR code", http.StatusBadRequest)
		return
	}

	// atomically claim the token: this single conditional UPDATE is the whole
	// concurrency guard, only one concurrent request can match "used_at IS NULL"
	// for a given row, so a code can never be redeemed twice
	now := time.Now()
	claim := database.DB.Model(&database.QrToken{}).
		Where("token = ? AND used_at IS NULL AND expires_at > ?", token, now).
		Update("used_at", now)
	if claim.Error != nil {
		http.Error(w, "Failed to process QR code", http.StatusInternalServerError)
		return
	}
	//rowsaffected come from the update, if no rows were updated, it means the token was already used or expired
	if claim.RowsAffected == 0 {
		http.Error(w, "QR code not found, expired, or already used", http.StatusConflict)
		return
	}

	var qrToken database.QrToken
	if err := database.DB.Where("token = ?", token).First(&qrToken).Error; err != nil {
		http.Error(w, "QR code not found", http.StatusNotFound)
		return
	}

	var client database.Client
	if err := database.DB.First(&client, qrToken.ClientID).Error; err != nil {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}
	if client.Balance < req.Amount {
		http.Error(w, "Insufficient balance", http.StatusPaymentRequired)
		return
	}

	client.Balance -= req.Amount
	partner.Balance += req.Amount
	database.DB.Save(&client)
	database.DB.Save(partner)

	transaction := database.Transaction{
		SenderUserID:   client.UserID,
		ReceiverUserID: partner.UserID,
		QrTokenID:      &qrToken.ID,
		Amount:         req.Amount,
		Type:           database.TransactionTypeDebit,
		IdempotencyKey: &req.IdempotencyKey,
	}
	database.DB.Create(&transaction)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"transaction": transaction})
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
