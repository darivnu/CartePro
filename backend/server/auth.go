//
// EPITECH PROJECT, 2026
// CartePro
// File description:
// auth
//

package server

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"cartepro/database"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type LoginResponse struct {
	User UserResponse `json:"user"`
}

type UserInfoResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func check_valid_user(req AuthRequest) (database.User, int) {
	var user database.User

	database.DB.Where("email = ?", req.Email).First(&user) //returns first record that matches the condition, or an error if no record is found

	//should we return 404 for not found here?
	if user.ID == 0 {
		//return 401 (Unauthorized) if user not found
		return user, http.StatusUnauthorized
	}

	//now check pass
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))

	if err != nil {
		//return 401 (Unauthorized) if password does not match
		return user, http.StatusUnauthorized
	}

	return user, http.StatusOK

}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, status := check_valid_user(req)

	if status != 200 {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	session, err := CreateSession(user.ID)
	if err != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	SetSessionCookie(w, session)

	// Send a successful login response
	response := LoginResponse{
		User: UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Role:  string(user.Role),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := GetSessionFromRequest(r)
	if err != nil {
		http.Error(w, "No active session", http.StatusUnauthorized)
		return
	}

	if err := deleteSession(session.Token); err != nil {
		http.Error(w, "Error deleting session", http.StatusInternalServerError)
		return
	}

	// Clear the session cookie
	clearSessionCookie(w)

	//return 204 No Content
	w.WriteHeader(http.StatusNoContent)

}

func HandleUserInfo(w http.ResponseWriter, r *http.Request) {
	session, err := GetSessionFromRequest(r)
	if err != nil {
		http.Error(w, "No active session", http.StatusUnauthorized)
		return
	}

	var user database.User
	if err := database.DB.First(&user, session.UserID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := UserInfoResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  string(user.Role),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
