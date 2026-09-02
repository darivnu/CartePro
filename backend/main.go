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

	"cartepro/database"
	"cartepro/routing"
)

func main() {
	log.Println("Starting server")

	database.InitDatabase()

	mux := routing.NewRouter()
	database.SeedTestUser()

	addr := ":4242"
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, routing.CorsMiddleware(mux)); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
