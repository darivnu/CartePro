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
		ClientID:       client.ID,
		PartnerID:      &partner.ID,
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
