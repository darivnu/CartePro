//
// EPITECH PROJECT, 2026
// CartePro
// File description:
// qrcode
//

package users

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"cartepro/database"
	"cartepro/server"

	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func generateQrToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil // stored as QrToken.Token
}

func signQrPayload(token string, expiresAt time.Time) string {
	payload := fmt.Sprintf("%s.%d", token, expiresAt.Unix())
	mac := hmac.New(sha256.New, []byte(os.Getenv("QR_SIGNING_SECRET")))
	mac.Write([]byte(payload))
	sig := hex.EncodeToString(mac.Sum(nil))
	// base64 first so the embedded "." in payload doesn't collide with our separator
	return base64.RawURLEncoding.EncodeToString([]byte(payload)) + "." + sig
}

func verifyQrPayload(qrPayload string) (token string, expiresAt time.Time, err error) {
	parts := strings.SplitN(qrPayload, ".", 2)
	if len(parts) != 2 {
		return "", time.Time{}, errors.New("malformed qr payload")
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", time.Time{}, errors.New("malformed qr payload")
	}

	mac := hmac.New(sha256.New, []byte(os.Getenv("QR_SIGNING_SECRET")))
	mac.Write(payloadBytes)
	sigBytes, err := hex.DecodeString(parts[1])
	if err != nil || !hmac.Equal(mac.Sum(nil), sigBytes) { // constant-time compare
		return "", time.Time{}, errors.New("invalid signature")
	}

	payloadParts := strings.SplitN(string(payloadBytes), ".", 2)
	expiresUnix, err := strconv.ParseInt(payloadParts[1], 10, 64)
	if err != nil {
		return "", time.Time{}, errors.New("malformed qr payload")
	}
	return payloadParts[0], time.Unix(expiresUnix, 0), nil
}

func HandleQrCodeGeneration(w http.ResponseWriter, r *http.Request) {
	_, client, err := server.GetClientFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	hash, err := generateQrToken()
	if err != nil {
		http.Error(w, "Failed to generate QR token", http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(5 * time.Minute) // QR code expires in 5 minutes
	qrPayload := signQrPayload(hash, expiresAt)

	qrToken := database.QrToken{
		Token:     hash,
		ClientID:  client.ID,
		ExpiresAt: expiresAt,
	}

	if err := database.DB.Create(&qrToken).Error; err != nil {
		http.Error(w, "Failed to store QR token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":      hash,
		"qr_payload": qrPayload,
		"expires_at": expiresAt.Format(time.RFC3339),
	})

}
