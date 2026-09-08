//
// EPITECH PROJECT, 2026
// CartePro
// File description:
// external
//

package users

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"
)

func checkHRISApiKey(r *http.Request) bool {
	provided := r.Header.Get("X-Api-Key") // Go canonicalizes the header namyou
	expected := os.Getenv("HRIS_API_KEY")
	if expected == "" || provided == "" {
		return false
	}
	return provided == expected
}

func HandleExternalIntegrationBalanceClientGetter(w http.ResponseWriter, r *http.Request) {

	if !checkHRISApiKey(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)

	if err != nil || id <= 0 {
		http.Error(w, "{\"error\": \"Invalid client ID\"}", http.StatusBadRequest)
		return
	}

	_, client, err := getClientByID(id)
	if err != nil {
		http.Error(w, "{\"error\": \"Client not found\"}", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"client_id": client.ID,
		"user_id":   client.UserID,
		"balance":   client.Balance,
		"currency":  "euro", // Assuming the balance is in euros; adjust as necessary
		"as_of":     time.Now().UTC().Format(time.RFC3339),
	})
}
