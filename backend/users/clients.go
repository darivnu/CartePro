//
// EPITECH PROJECT, 2026
// CartePro
// File description:
// clients
//

package users

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"cartepro/database"
	"cartepro/server"
)

type ClientRegistrationRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	Name       string `json:"name"`
	EmployerID uint   `json:"employer_id"`
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

	//check if employer exists
	var employer database.Employer
	var employer_exists bool = true
	if err := database.DB.First(&employer, req.EmployerID).Error; err != nil {
		employer_exists = false
		log.Printf("Employer with ID %d does not exist, proceeding without employer association", req.EmployerID)
	}

	//if employer does not exist, check if there the one they sent is 0 (so its a placeholder meant to be updated later). if it isn't, the error
	if req.EmployerID != 0 && !employer_exists {
		http.Error(w, "Employer not found", http.StatusBadRequest)
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

	var client database.Client
	if employer_exists {
		client = database.Client{
			UserID:     user.ID,
			Name:       req.Name,
			EmployerID: &employer.ID,
		}
	} else {
		client = database.Client{
			UserID:     user.ID,
			Name:       req.Name,
			EmployerID: nil,
		}
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

	session, err := server.GetSessionFromRequest(r)
	if err != nil {
		http.Error(w, "No active session", http.StatusUnauthorized)
		return
	}

	var user database.User
	if err := database.DB.First(&user, session.UserID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
}
