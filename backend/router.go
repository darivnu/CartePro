//
// EPITECH PROJECT, 2026
// CartePro
// File description:
// router
//

package main

import (
	"net/http"
)

func newRouter() *http.ServeMux {
	mux := http.NewServeMux()

	//all the routhing things go here
	//authentification
	mux.HandleFunc("POST /auth/login", handleLogin)
	mux.HandleFunc("POST /auth/logout", handleLogout)
	mux.HandleFunc("GET /auth/me", handleUserInfo)

	//partner
	mux.HandleFunc("POST /partners/register", handlePartnerRegistration)

	return mux
}
