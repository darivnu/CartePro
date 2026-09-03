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
	"os"

	"cartepro/database"
	"cartepro/routing"
)

func main() {
	log.Println("Starting server")

	database.InitDatabase()

	mux := routing.NewRouter()
	database.SeedTestClient()
	database.SeedTestPartners()

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "4242"
	}
	addr := ":" + port
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, routing.CorsMiddleware(mux)); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
