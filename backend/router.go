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
	mux.HandleFunc("POST /auth/login", handleLogin)

	return mux
}
