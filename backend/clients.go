//
// EPITECH PROJECT, 2026
// CartePro
// File description:
// clients
//

package main

import (
	"encoding/json"
	"net/http"
	"golang.org/x/crypto/bcrypt"
)

type ClientRegistrationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name	 string `json:"name"`
	EmployerID uint   `json:"employer_id"`
}

func handleClientRegistration(w http.ResponseWriter, r *http.Request) {
	var req ClientRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if the email already exists in the database
	var existingUser User
	db.Where("email = ?", req.Email).First(&existingUser)
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
	user := User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         RoleClient,
	}
	if err := db.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	client := Client{
		UserID:     user.ID,
		Name:       req.Name,
		EmployerID: req.EmployerID,
	}
	if err := db.Create(&client).Error; err != nil {
		http.Error(w, "Failed to create client", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"client": client,
	})
}