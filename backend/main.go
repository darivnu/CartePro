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
)

func main() {
	log.Println("Starting server")

	initDatabase()

	mux := newRouter()
	seedTestUser()

	addr := ":4242"
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
