//
// EPITECH PROJECT, 2026
// CartePro
// File description:
// router
//

package server

import (
	"net/http"

	"github.com/darivnu/cartepro/users"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	//all the routhing things go here
	//authentification
	mux.HandleFunc("POST /auth/login", handleLogin)
	mux.HandleFunc("POST /auth/logout", handleLogout)
	mux.HandleFunc("GET /auth/me", handleUserInfo)

	//partner
	mux.HandleFunc("POST /partners/register", users.HandlePartnerRegistration)

	//client
	mux.HandleFunc("POST /clients/register", users.HandleClientRegistration)

	//admin
	mux.HandleFunc("POST /admin/register", users.HandleAdminRegistration)

	return mux
}
