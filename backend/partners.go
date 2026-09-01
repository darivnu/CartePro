//
// EPITECH PROJECT, 2026
// CartePro
// File description:
// This handles all the different types of partners
//

package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type PartnerRegistrationRequest struct {
	businessName string `json:"business_name"`
	siret        string `json:"siret"`
	category     string `json:"category"`
	address      string `json:"address"`
	region       string `json:"region"`
	contactEmail string `json:"contact_email"`
	password     string `json:"password"`
}

func handlePartnerRegistration(w http.ResponseWriter, r *http.Request) {
	var req PartnerRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	//we gonna check first that the email does not already exist in the database
	var existingUser User
	db.Where("email = ?", req.contactEmail).First(&existingUser)
	if existingUser.ID != 0 {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	//hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	//create the new partner
	user := User{
		Email:        req.contactEmail,
		PasswordHash: string(hashedPassword),
		Role:         RolePartner,
	}

	if err := db.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	partner := Partner{
		UserID:       user.ID,
		BusinessName: req.businessName,
		Siret:        req.siret,
		Category:     req.category,
		Address:      req.address,
		Region:       req.region,
	}

	if err := db.Create(&partner).Error; err != nil {
		http.Error(w, "Failed to create partner", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}
