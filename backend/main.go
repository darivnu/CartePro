//
// EPITECH PROJECT, 2026
// CartePro
// File description:
// main
//

package main

import (
	"log"
	"net/http"

	"github.com/darivnu/cartepro/database"
	"github.com/darivnu/cartepro/server"
)

func main() {
	log.Println("Starting server")

	database.InitDatabase()

	mux := server.NewRouter()
	database.SeedTestUser()

	addr := ":4242"
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, server.CorsMiddleware(mux)); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
