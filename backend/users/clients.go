//
// EPITECH PROJECT, 2026
// CartePro
// File description:
// clients
//

package users

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"cartepro/database"
	"cartepro/server"
)

type ClientRegistrationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func HandleClientRegistration(w http.ResponseWriter, r *http.Request) {
	var req ClientRegistrationRequest
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

	//create new user and client
	user := database.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         database.RoleClient,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	client := database.Client{
		UserID: user.ID,
		Name:   req.Name,
	}

	if err := database.DB.Create(&client).Error; err != nil {
		http.Error(w, "Failed to create client", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"client": client,
	})
}

func HandleClientBalance(w http.ResponseWriter, r *http.Request) {

	_, client, err := server.GetClientFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"balance":    client.Balance,
		"updated_at": time.Now().UTC().Format(time.RFC3339), //return the current time in UTC as the updated_at field
	})
}
