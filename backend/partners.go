//
// EPITECH PROJECT, 2026
// CartePro
// File description:
// This handles all the different types of partners
//

package main

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type PartnerRegistrationRequest struct {
	BusinessName string `json:"business_name"`
	Siret        string `json:"siret"`
	Category     string `json:"category"`
	Address      string `json:"address"`
	Region       string `json:"region"`
	ContactEmail string `json:"contact_email"`
	Password     string `json:"password"`
}

func handlePartnerRegistration(w http.ResponseWriter, r *http.Request) {
	var req PartnerRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	//we gonna check first that the email does not already exist in the database
	var existingUser User
	db.Where("email = ?", req.ContactEmail).First(&existingUser)
	if existingUser.ID != 0 {
		http.Error(w, "Email already exists", http.StatusConflict)
		log.Println("Existing User email:", existingUser.Email)
		return
	}

	//then we gonna check that the siret or buisness name  does not already exist in the database
	var existingPartner Partner
	db.Where("siret = ? OR business_name = ?", req.Siret, req.BusinessName).First(&existingPartner)
	if existingPartner.ID != 0 {
		http.Error(w, "Siret or Business Name already exists", http.StatusConflict)
		log.Println("Existing Partner siret:", existingPartner.Siret)
		log.Println("Existing Partner business name:", existingPartner.BusinessName)
		return
	}

	//hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	//create the new partner
	user := User{
		Email:        req.ContactEmail,
		PasswordHash: string(hashedPassword),
		Role:         RolePartner,
	}

	if err := db.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	partner := Partner{
		UserID:       user.ID,
		BusinessName: req.BusinessName,
		Siret:        req.Siret,
		Category:     req.Category,
		Address:      req.Address,
		Region:       req.Region,
	}

	if err := db.Create(&partner).Error; err != nil {
		http.Error(w, "Failed to create partner", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}
