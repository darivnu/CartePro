package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"
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

func createSession(userID uint) (*Session, error) {
	token, err := generateSessionToken()
	if err != nil {
		return nil, err
	}

	session := &Session{
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(sessionDuration),
	}
	if err := db.Create(session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func setSessionCookie(w http.ResponseWriter, session *Session) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    session.Token,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// Secure: true, // enable once the API is served over HTTPS
	})
}
