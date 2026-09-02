//
// EPITECH PROJECT, 2026
// CartePro
// File description:
// router
//

package routing

import (
	"net/http"

	"cartepro/server"
	"cartepro/users"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	//all the routhing things go here
	//authentification
	mux.HandleFunc("POST /auth/login", server.HandleLogin)
	mux.HandleFunc("POST /auth/logout", server.HandleLogout)
	mux.HandleFunc("GET /auth/me", server.HandleUserInfo)

	//partner
	mux.HandleFunc("POST /partners/register", users.HandlePartnerRegistration)

	//client
	mux.HandleFunc("POST /clients/register", users.HandleClientRegistration)
	mux.HandleFunc("GET /clients/me/balance", users.HandleClientBalance)
	mux.HandleFunc("POST /clients/me/qrcode", users.HandleQrCodeGeneration)

	//admin
	mux.HandleFunc("POST /admin/register", users.HandleAdminRegistration)

	return mux
}
