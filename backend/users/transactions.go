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
)

type QrValidationRequest struct {
	QrPayload      string `json:"qr_payload"`
	Amount         int64  `json:"amount"`
	IdempotencyKey string `json:"idempotency_key"`
}

func HandleQrCodeValidation(w http.ResponseWriter, r *http.Request) {
	var req QrValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, expiresAt, err := verifyQrPayload(req.QrPayload)
	defer database.DB.Where("token = ?", token).Delete(&database.QrToken{}) // delete the token after processing, regardless of success or failure

	if err != nil {
		http.Error(w, "Invalid QR payload", http.StatusBadRequest)
		return
	}

	if time.Now().After(expiresAt) {
		http.Error(w, "QR code has expired", http.StatusUnauthorized)
		return
	}

	//process the validation logic here

}
