//
// EPITECH PROJECT, 2026
// CartePro
// File description:
// admin
//

package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"cartepro/database"
	"cartepro/server"
)

// maxTopupAmount caps a single admin top-up at 10,000.00 EUR (balances are stored in cents).
const maxTopupAmount = 1_000_000

type AdminRegistrationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func HandleAdminRegistration(w http.ResponseWriter, r *http.Request) {

	var req AdminRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if the email already exists in the database
	var existingUser database.User
	database.DB.Where("email = ?", req.Email).First(&existingUser)
	if existingUser.ID != 0 {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user := database.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         database.RoleAdmin,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	admin := database.Admin{
		Name: req.Name,
		User: user,
	}

	if err := database.DB.Create(&admin).Error; err != nil {
		http.Error(w, "Failed to create admin", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"admin": admin,
	})

}

func HandleApprovePartner(w http.ResponseWriter, r *http.Request) {
	user, _, err := server.GetAdminFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if user.Role != database.RoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)

	if err != nil {
		http.Error(w, "Invalid partner ID", http.StatusBadRequest)
		return
	}

	_, partner, err := getPartnerByID(uint(id))
	if err != nil {
		http.Error(w, "Partner not found", http.StatusNotFound)
		return
	}

	partner.Status = database.StatusApproved
	if err := database.DB.Save(partner).Error; err != nil {
		http.Error(w, "Failed to approve partner", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"partner": partner,
	})
}

type AdminTopupRequest struct {
	ClientID uint   `json:"client_id"`
	Amount   int64  `json:"amount"`
	Comment  string `json:"comment"`
}

type CancelTransactionRequest struct {
	Reason string `json:"reason"`
}

func debitcancelTransaction(tx *gorm.DB, originalTx database.Transaction, admin *database.Admin, reason string) (database.Transaction, error) {

	reversalTx := database.Transaction{
		ReceiverUserID:        originalTx.SenderUserID,
		SenderUserID:          originalTx.ReceiverUserID,
		QrTokenID:             nil,
		Amount:                originalTx.Amount,
		Type:                  database.TransactionTypeReversal,
		Comment:               &reason,
		OriginalTransactionID: &originalTx.ID,
	}

	// the unique constraint on original_transaction_id is what actually blocks a
	// double cancellation: a raced or repeated cancel request hits this insert
	// and fails here, inside the transaction, instead of a separate check-then-act query
	if err := tx.Create(&reversalTx).Error; err != nil {
		return database.Transaction{}, err
	}

	// credit the client who originally paid
	if err := tx.Model(&database.Client{}).
		Where("user_id = ?", originalTx.SenderUserID).
		Update("balance", gorm.Expr("balance + ?", originalTx.Amount)).Error; err != nil {
		return database.Transaction{}, err
	}

	// debit the partner who originally received
	if err := tx.Model(&database.Partner{}).
		Where("user_id = ?", originalTx.ReceiverUserID).
		Update("balance", gorm.Expr("balance - ?", originalTx.Amount)).Error; err != nil {
		return database.Transaction{}, err
	}

	//update the original transaction to mark it as cancelled
	if err := tx.Model(&database.Transaction{}).
		Where("id = ?", originalTx.ID).
		Update("cancelled", true).Error; err != nil {
		return database.Transaction{}, err
	}

	return reversalTx, nil
}

func HandleCancelTransaction(w http.ResponseWriter, r *http.Request) {
	user, _, err := server.GetAdminFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if user.Role != database.RoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	var req CancelTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var originalTx database.Transaction
	if err := database.DB.First(&originalTx, id).Error; err != nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	if originalTx.Type != database.TransactionTypeDebit {
		http.Error(w, "Only debit transactions can be cancelled", http.StatusBadRequest)
		return
	}

	var reversalTx database.Transaction
	txErr := database.DB.Transaction(func(tx *gorm.DB) error {
		var err error
		reversalTx, err = debitcancelTransaction(tx, originalTx, nil, req.Reason)
		return err
	})

	switch {
	case txErr == nil:
	case isUniqueViolation(txErr):
		http.Error(w, "Transaction already cancelled", http.StatusConflict)
		return
	default:
		http.Error(w, "Failed to cancel transaction", http.StatusInternalServerError)
		return
	}

	response, err := toTransactionResponse(reversalTx)
	if err != nil {
		http.Error(w, "Failed to cancel transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"transaction": response,
	})
}

func HandleAdminTopups(w http.ResponseWriter, r *http.Request) {
	user, admin, err := server.GetAdminFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if user.Role != database.RoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	var req AdminTopupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount <= 0 || req.Amount > maxTopupAmount {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// balance move and ledger insert run in one DB transaction (see
	// creditClientTopup) so a failure partway through (client vanishes, insert
	// fails, ...) can't leave the balance bumped with no matching transaction row,
	// or vice versa
	var transaction database.Transaction
	txErr := database.DB.Transaction(func(tx *gorm.DB) error {
		var err error
		transaction, err = creditClientTopup(tx, req, admin)
		return err
	})

	switch {
	case txErr == nil:
		response, err := toTransactionResponse(transaction)
		if err != nil {
			http.Error(w, "Failed to process transaction", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"transaction": response})
	case errors.Is(txErr, gorm.ErrRecordNotFound):
		http.Error(w, "Client not found", http.StatusNotFound)
	default:
		http.Error(w, "Failed to process transaction", http.StatusInternalServerError)
	}
}

// creditClientTopup bumps the client's balance and inserts the ledger row, both
// against tx so the caller can run it inside a single DB transaction.
func creditClientTopup(tx *gorm.DB, req AdminTopupRequest, admin *database.Admin) (database.Transaction, error) {
	var client database.Client
	if err := tx.First(&client, req.ClientID).Error; err != nil {
		return database.Transaction{}, err
	}

	// atomic increment: no read-modify-write window for a concurrent top-up
	// on the same client to race and clobber
	if err := tx.Model(&database.Client{}).
		Where("id = ?", client.ID).
		Update("balance", gorm.Expr("balance + ?", req.Amount)).Error; err != nil {
		return database.Transaction{}, err
	}

	transaction := database.Transaction{
		ReceiverUserID: client.UserID,
		SenderUserID:   admin.UserID,
		QrTokenID:      nil,
		Amount:         req.Amount,
		Type:           database.TransactionTypeTopup,
		Comment:        &req.Comment,
	}
	if err := tx.Create(&transaction).Error; err != nil {
		return database.Transaction{}, err
	}
	return transaction, nil
}

func HandleGetAdminPartners(w http.ResponseWriter, r *http.Request) {
	user, _, err := server.GetAdminFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if user.Role != database.RoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	query := r.URL.Query()
	status := query.Get("status")
	if status != "" && status != string(database.StatusPending) && status != string(database.StatusApproved) && status != string(database.StatusRejected) {
		http.Error(w, "Invalid status filter", http.StatusBadRequest)
		return
	}
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

	var partners []database.Partner
	dbQuery := database.DB.Model(&database.Partner{}).Order("id")
	if status != "" {
		dbQuery = dbQuery.Where("status = ?", status)
	}
	if err := dbQuery.Offset((page - 1) * limit).Limit(limit).Find(&partners).Error; err != nil {
		http.Error(w, "Failed to fetch partners", http.StatusInternalServerError)
		return
	}

	total := len(partners)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": partners,
		"meta": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})

}

func HandleGetAdminClients(w http.ResponseWriter, r *http.Request) {
	user, _, err := server.GetAdminFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if user.Role != database.RoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
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

	var clients []database.Client
	dbQuery := database.DB.Model(&database.Client{}).Preload("User").Order("id")
	if err := dbQuery.Offset((page - 1) * limit).Limit(limit).Find(&clients).Error; err != nil {
		http.Error(w, "Failed to fetch clients", http.StatusInternalServerError)
		return
	}

	total := len(clients)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": clients,
		"meta": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})

}

func HandleGetAdminClientDetail(w http.ResponseWriter, r *http.Request) {
	user, _, err := server.GetAdminFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if user.Role != database.RoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	var client database.Client
	if err := database.DB.Preload("User").First(&client, id).Error; err != nil {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	var transactions []database.Transaction
	if err := database.DB.Where("receiver_user_id = ? OR sender_user_id = ?", client.UserID, client.UserID).Order("created_at DESC").Find(&transactions).Error; err != nil {
		http.Error(w, "Failed to fetch top-ups", http.StatusInternalServerError)
		return
	}

	responses, err := toTransactionResponses(transactions)
	if err != nil {
		http.Error(w, "Failed to fetch top-ups", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"client":       client,
		"transactions": responses,
	})
}

type MinisterPickRequest struct {
	Reason bool `json:"minister_pick"`
}

func HandleMinisterPick(w http.ResponseWriter, r *http.Request) {
	_, _, err := server.GetAdminFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid partner ID", http.StatusBadRequest)
		return
	}
	var req MinisterPickRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	var partner database.Partner
	if err := database.DB.Preload("User").First(&partner, id).Error; err != nil {
		http.Error(w, "Partner not found", http.StatusNotFound)
		return
	}
	partner.MinisterPick = req.Reason
	if err := database.DB.Save(&partner).Error; err != nil {
		http.Error(w, "Failed to pick partner", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"partner": partner,
	})
}

type AdminRejectPartnerRequest struct {
	Reason string `json:"reason"`
}

func HandleAdminRejectPartner(w http.ResponseWriter, r *http.Request) {
	user, _, err := server.GetAdminFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if user.Role != database.RoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)

	if err != nil {
		http.Error(w, "Invalid partner ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Reason == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	var partner database.Partner
	if err := database.DB.Preload("User").First(&partner, id).Error; err != nil {
		http.Error(w, "Partner not found", http.StatusNotFound)
		return
	}
	partner.Status = database.StatusRejected
	if (req.Reason) != "" {
		partner.RejectReason = &req.Reason
	}
	if err := database.DB.Save(&partner).Error; err != nil {
		http.Error(w, "Failed to reject partner", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"partner": partner,
	})
}

type UpdatePartnerStatusRequest struct {
	Status string `json:"status"`
}

func HandleUpdatePartnerStatus(w http.ResponseWriter, r *http.Request) {
	_, _, err := server.GetAdminFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid partner ID", http.StatusBadRequest)
		return
	}
	var req UpdatePartnerStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || (req.Status != string(database.StatusPending) && req.Status != string(database.StatusApproved) && req.Status != string(database.StatusRejected)) {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	var partner database.Partner
	if err := database.DB.Preload("User").First(&partner, id).Error; err != nil {
		http.Error(w, "Partner not found", http.StatusNotFound)
		return
	}
	partner.Status = database.PartnerStatus(req.Status)
	if err := database.DB.Save(&partner).Error; err != nil {
		http.Error(w, "Failed to update partner status", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"partner": partner,
	})
}

func HandleGetAdminDashboard(w http.ResponseWriter, r *http.Request) {
	_, _, err := server.GetAdminFromSession(r)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var active_partners []database.Partner
	if err := database.DB.Model(&database.Partner{}).Preload("User").Where("status = ?", database.StatusApproved).Find(&active_partners).Error; err != nil {
		http.Error(w, "Failed to fetch active partners", http.StatusInternalServerError)
		return
	}

	var total_clients int64
	if err := database.DB.Model(&database.Client{}).Count(&total_clients).Error; err != nil {
		http.Error(w, "Failed to fetch total clients count", http.StatusInternalServerError)
		return
	}

	var list_of_regions []string
	if err := database.DB.Model(&database.Partner{}).Distinct("Region").Pluck("Region", &list_of_regions).Error; err != nil {
		http.Error(w, "Failed to fetch list of locations", http.StatusInternalServerError)
		return
	}
	// sum debit transaction volume, optionally restricted to a from/to created_at range
	query := r.URL.Query()
	from := query.Get("from")
	to := query.Get("to")

	db := database.DB.Model(&database.Transaction{})
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

	var total_transaction_volume int64
	if err := db.Where("type = ?", database.TransactionTypeDebit).Select("COALESCE(SUM(amount), 0)").Scan(&total_transaction_volume).Error; err != nil {
		http.Error(w, "Failed to fetch total transaction volume", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"active_partners":          active_partners,
		"total_transaction_volume": total_transaction_volume,
		"total_clients":            total_clients,
		"list_of_regions":          list_of_regions,
	})

}
