//
// EPITECH PROJECT, 2026
// CartePro
// File description:
// session
//

package server

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"cartepro/database"
)

const sessionCookieName = "session_token"
const sessionDuration = 7 * 24 * time.Hour

func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func createSession(userID uint) (*database.Session, error) {
	token, err := generateSessionToken()
	if err != nil {
		return nil, err
	}

	session := &database.Session{
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(sessionDuration),
	}
	if err := database.DB.Create(session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func setSessionCookie(w http.ResponseWriter, session *database.Session) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    session.Token,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, //samesitelax mode is a security feature that prevents cross-site request forgery
		// Secure: true, // enable once the API is served over HTTPS
	})
}

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func deleteSession(token string) error {
	return database.DB.Where("token = ?", token).Delete(&database.Session{}).Error
}

// this function retrieves the session from the request cookie and checks if it's valid
func GetSessionFromRequest(r *http.Request) (*database.Session, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil, err
	}

	var session database.Session
	if err := database.DB.Where("token = ?", cookie.Value).First(&session).Error; err != nil {
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		deleteSession(session.Token)
		return nil, http.ErrNoCookie
	}

	return &session, nil
}
