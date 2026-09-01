//
// EPITECH PROJECT, 2026
// CartePro
// File description:
// this file is used to seed the database with a test user
//

package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func seedTestUser() {
	var count int64
	db.Model(&User{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded, skipping test user creation")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash seed password: %v", err)
	}

	user := User{
		Email:        "test@cartepro.dev",
		PasswordHash: string(hash),
		Role:         RoleClient,
	}
	if err := db.Create(&user).Error; err != nil {
		log.Fatalf("failed to seed test user: %v", err)
	}

	log.Printf("seeded test user -> email: test@cartepro.dev password: password123")
}
