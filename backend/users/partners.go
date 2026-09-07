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
	"time"

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
		ID           uint   `json:"ID"`
		BusinessName string `json:"BusinessName"`
		Siret        string `json:"Siret"`
		Category     string `json:"Category"`
		Address      string `json:"Address"`
		Region       string `json:"Region"`
		Email        string `json:"Email"`
		Status       string `json:"Status"`
		MinisterPick bool   `json:"MinisterPick"`
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

// withoutPasswordHash returns a copy of user with PasswordHash cleared, so it's safe to serialize.
func withoutPasswordHash(user database.User) database.User {
	user.PasswordHash = ""
	return user
}

func HandleGetOwnPartnerInfo(w http.ResponseWriter, r *http.Request) {
	user, partner, err := server.GetPartnerFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	partner.User = withoutPasswordHash(*user) //we add the user to the partner so that we can return the email in the response, without leaking the password hash

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"partner": partner})
}

type PublicPartnerResponse struct {
	ID           uint   `json:"ID"`
	BusinessName string `json:"BusinessName"`
	Category     string `json:"Category"`
	Address      string `json:"Address"`
	Region       string `json:"Region"`
	MinisterPick bool   `json:"MinisterPick"`
}

type PartnersMeta struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

func HandleGetMultiplePartners(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	search := query.Get("search")
	category := query.Get("category")

	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 1 {
		page = 1 //the page is which chunk we want (defined by limit), so if we want the first chunk, we set page to 1)
	}

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	db := database.DB.Model(&database.Partner{}).Where("status = ?", database.StatusApproved)
	if search != "" {
		db = db.Where("business_name LIKE ?", "%"+search+"%")
	}
	if category != "" {
		db = db.Where("category = ?", category)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		http.Error(w, "Failed to count partners", http.StatusInternalServerError)
		return
	}

	var partners []database.Partner
	if err := db.Order("id").Offset((page - 1) * limit).Limit(limit).Find(&partners).Error; err != nil {
		http.Error(w, "Failed to fetch partners", http.StatusInternalServerError)
		return
	}

	data := make([]PublicPartnerResponse, len(partners))
	for i, partner := range partners {
		data[i] = PublicPartnerResponse{
			ID:           partner.ID,
			BusinessName: partner.BusinessName,
			Category:     partner.Category,
			Address:      partner.Address,
			Region:       partner.Region,
			MinisterPick: partner.MinisterPick,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": data,
		"meta": PartnersMeta{Page: page, Limit: limit, Total: total},
	})
}

// this func looks incredibly similar to HandleGetClientOwnTransactions
// TODO: we should refactor this to avoid code duplication, but for now we will keep it like this
func HandleGetOwnTransactions(w http.ResponseWriter, r *http.Request) {
	_, partner, err := server.GetPartnerFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	query := r.URL.Query()

	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 1 {
		page = 1 //the page is which chunk we want (defined by limit), so if we want the first chunk, we set page to 1)
	}

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	from := query.Get("from")
	to := query.Get("to")

	db := database.DB.Model(&database.Transaction{}).Where("sender_user_id = ?", partner.UserID).Or("receiver_user_id = ?", partner.UserID).Order("created_at DESC")
	if from != "" {
		fromTime, err := time.Parse(time.RFC3339, from)
		if err != nil {
			http.Error(w, "Invalid 'from' date format", http.StatusBadRequest)
			return
		}
		db = db.Where("created_at >= ?", fromTime)
	}
	if to != "" {
		toTime, err := time.Parse(time.RFC3339, to)
		if err != nil {
			http.Error(w, "Invalid 'to' date format", http.StatusBadRequest)
			return
		}
		db = db.Where("created_at <= ?", toTime)
	}

	var total int64
	db.Count(&total)

	var transactions []database.Transaction
	if err := db.Offset((page - 1) * limit).Limit(limit).Find(&transactions).Error; err != nil {
		http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": transactions,
		"meta": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
			"from":        from,
			"to":          to,
		},
	})
}
