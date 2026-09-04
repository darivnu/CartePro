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
	mux.HandleFunc("POST /partners/me/validate", users.HandleQrCodeValidation)
	mux.HandleFunc("GET /partners", users.HandleGetMultiplePartners)
	mux.HandleFunc("GET /partners/{id}", users.HandleGetSpecificPartner)
	mux.HandleFunc("GET /partners/me", users.HandleGetOwnPartnerInfo)

	//client
	mux.HandleFunc("POST /clients/register", users.HandleClientRegistration)
	mux.HandleFunc("GET /clients/me/balance", users.HandleClientBalance)
	mux.HandleFunc("POST /clients/me/qrcode", users.HandleQrCodeGeneration)
	mux.HandleFunc("GET /clients/me/transactions", users.HandleGetClientOwnTransactions)

	//admin
	mux.HandleFunc("POST /admin/register", users.HandleAdminRegistration)
	mux.HandleFunc("POST /admin/partners/{id}/approve", users.HandleApprovePartner)
	mux.HandleFunc("POST /admin/topups", users.HandleAdminTopups)
	mux.HandleFunc("POST /admin/transactions/{id}/cancel", users.HandleCancelTransaction)

	return mux
}
