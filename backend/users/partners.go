//
// EPITECH PROJECT, 2026
// CartePro
// File description:
// This handles all the different types of partners
//

package users

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"cartepro/database"
	"cartepro/server"
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

func HandlePartnerRegistration(w http.ResponseWriter, r *http.Request) {
	var req PartnerRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	//we gonna check first that the email does not already exist in the database
	var existingUser database.User
	database.DB.Where("email = ?", req.ContactEmail).First(&existingUser)
	if existingUser.ID != 0 {
		http.Error(w, "Email already exists", http.StatusConflict)
		log.Println("Existing User email:", existingUser.Email)
		return
	}

	//then we gonna check that the siret or buisness name  does not already exist in the database
	var existingPartner database.Partner
	database.DB.Where("siret = ? OR business_name = ?", req.Siret, req.BusinessName).First(&existingPartner)
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
	user := database.User{
		Email:        req.ContactEmail,
		PasswordHash: string(hashedPassword),
		Role:         database.RolePartner,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	partner := database.Partner{
		UserID:       user.ID,
		BusinessName: req.BusinessName,
		Siret:        req.Siret,
		Category:     req.Category,
		Address:      req.Address,
		Region:       req.Region,
	}

	if err := database.DB.Create(&partner).Error; err != nil {
		http.Error(w, "Failed to create partner", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func getPartnerByUserID(userID uint) (*database.User, *database.Partner, error) {
	var user database.User
	var partner database.Partner
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, nil, err
	}
	if err := database.DB.Where("user_id = ?", userID).First(&partner).Error; err != nil {
		return nil, nil, err
	}
	return &user, &partner, nil
}

func getPartnerByID(partnerID uint) (*database.User, *database.Partner, error) {
	var partner database.Partner
	var user database.User
	if err := database.DB.First(&partner, partnerID).Error; err != nil {
		return nil, nil, err
	}
	if err := database.DB.First(&user, partner.UserID).Error; err != nil {
		return nil, nil, err
	}
	return &user, &partner, nil
}

func HandleGetSpecificPartner(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid partner ID", http.StatusBadRequest)
		return
	}

	user, partner, err := getPartnerByID(uint(id))
	if err != nil {
		http.Error(w, "Partner not found", http.StatusNotFound)
		return
	}

	//we want to return only specific fields of the partner, not all of them, so we create a new struct to hold only the fields we want to return

	type PartnerResponse struct {
		ID           uint   `json:"id"`
		BusinessName string `json:"business_name"`
		Siret        string `json:"siret"`
		Category     string `json:"category"`
		Address      string `json:"address"`
		Region       string `json:"region"`
		Email        string `json:"email"`
		Status       string `json:"status"`
		MinisterPick bool   `json:"minister_pick"`
	}
	var response PartnerResponse
	response.ID = partner.ID
	response.BusinessName = partner.BusinessName
	response.Siret = partner.Siret
	response.Category = partner.Category
	response.Address = partner.Address
	response.Region = partner.Region
	response.Email = user.Email
	response.Status = string(partner.Status)
	response.MinisterPick = partner.MinisterPick

	//now we encode response to json like : {"partner": { ... }}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"partner": response})
}

func HandleGetOwnPartnerInfo(w http.ResponseWriter, r *http.Request) {
	user, partner, err := server.GetPartnerFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	partner.User = *user //we add the user to the partner so that we can return the email in the response

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"partner": partner})
}
