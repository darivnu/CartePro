//
// EPITECH PROJECT, 2026
// CartePro
// File description:
// admin
//

package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"cartepro/database"
	"cartepro/server"
)

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
	ClientID       uint   `json:"client_id"`
	Amount         int64  `json:"amount"`
	Comment        string `json:"comment"`
}

type CancelTransactionRequest struct {
	Reason string `json:"reason"`
}

func HandleCancelTransaction(w http.ResponseWriter, r *http.Request) {
	user, admin, err := server.GetAdminFromSession(r)
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
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	var req CancelTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var originalTx database.Transaction
	if err := database.DB.Preload("Client").First(&originalTx, id).Error; err != nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	if originalTx.Type != database.TransactionTypeDebit {
		http.Error(w, "Only debit transactions can be cancelled", http.StatusBadRequest)
		return
	}

	reversalTx := database.Transaction{
		ClientID:       originalTx.ClientID,
		PartnerID:      originalTx.PartnerID,
		QrTokenID:      nil,
		Amount:         originalTx.Amount,
		Type:           database.TransactionTypeReversal,
		Comment:        &req.Reason,
		AdminID:        &admin.ID,
		IdempotencyKey: nil,
	}

	if err := database.DB.Create(&reversalTx).Error; err != nil {
		http.Error(w, "Failed to create reversal transaction", http.StatusInternalServerError)
		return
	}

	originalTx.Client.Balance += originalTx.Amount
	if err := database.DB.Save(&originalTx.Client).Error; err != nil {
		http.Error(w, "Failed to update client balance", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"transaction": reversalTx,
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var client database.Client
	if err := database.DB.First(&client, req.ClientID).Error; err != nil {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	client.Balance += req.Amount
	database.DB.Save(&client)

	transaction := database.Transaction{
		ClientID:       client.ID,
		PartnerID:      nil,
		QrTokenID:      nil,
		Amount:         req.Amount,
		Type:           database.TransactionTypeTopup,
		Comment:        &req.Comment,
		AdminID:        &admin.ID,
	}
	database.DB.Create(&transaction)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"transaction": transaction})
}
