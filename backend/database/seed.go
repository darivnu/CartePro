//
// EPITECH PROJECT, 2026
// CartePro
// File description:
// this file is used to seed the database with a test user
//

package database

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func SeedTestClient() {
	var count int64
	DB.Model(&Client{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded, skipping test client creation")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash seed password: %v", err)
	}

	user := User{
		Email:        "testclient@cartepro.dev",
		PasswordHash: string(hash),
		Role:         RoleClient,
	}
	if err := DB.Create(&user).Error; err != nil {
		log.Fatalf("failed to seed test user: %v", err)
	}

	client := Client{
		UserID:  user.ID,
		Name:    "jack",
		Balance: 1000, // Set an initial balance for the test client (in cents, so 10 euros)
	}
	if err := DB.Create(&client).Error; err != nil {
		log.Fatalf("failed to seed test client: %v", err)
	}

	log.Printf("seeded test client -> email: testclient@cartepro.dev password: password123")
}

func SeedTestAdmin() {
	var count int64
	DB.Model(&Admin{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded, skipping test admin creation")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash seed password: %v", err)
	}

	user := User{
		Email:        "testadmin@cartepro.dev",
		PasswordHash: string(hash),
		Role:         RoleAdmin,
	}
	if err := DB.Create(&user).Error; err != nil {
		log.Fatalf("failed to seed test user: %v", err)
	}

	admin := Admin{
		Name: "John Doe",
		User: user,
	}
	if err := DB.Create(&admin).Error; err != nil {
		log.Fatalf("failed to seed test admin: %v", err)
	}

	log.Printf("seeded test admin -> email: testadmin@cartepro.dev password: password123")
}

func SeedTestPartners() {
	var count int64
	DB.Model(&Partner{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded, skipping test partner creation")
		return
	}

	partners := []struct {
		Email        string
		BusinessName string
		Siret        int64
		Category     string
		Address      string
		Region       string
	}{
		{
			Email:        "testpartner@cartepro.dev",
			BusinessName: "Test Partner",
			Siret:        12345678901234,
			Category:     "Test Category",
			Address:      "123 Test St, Test City",
			Region:       "Test Region",
		},
		{
			Email:        "contact@poneydream78.fr",
			BusinessName: "Poney Dream 78",
			Siret:        23456789012345,
			Category:     "Pony club",
			Address:      "Route des Écuries, Rambouillet",
			Region:       "Île-de-France",
		},
		{
			Email:        "contact@kostumparty.fr",
			BusinessName: "KostumParty",
			Siret:        34567890123456,
			Category:     "Costume shop",
			Address:      "25 Rue de la Roquette, Paris 11e",
			Region:       "Île-de-France",
		},
		{
			Email:        "contact@glaces-correze.fr",
			BusinessName: "Glaces Artisanales Corrèze",
			Siret:        45678901234567,
			Category:     "Ice cream maker (online + click & collect)",
			Address:      "Click & collect, Brive-la-Gaillarde",
			Region:       "Nouvelle-Aquitaine",
		},
		{
			Email:        "contact@chapelierfontaine.fr",
			BusinessName: "Chapelier Fontaine",
			Siret:        56789012345678,
			Category:     "Felt hats",
			Address:      "8 Rue Saint-Rome, Toulouse",
			Region:       "Occitanie",
		},
	}

	for _, p := range partners {
		hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("failed to hash seed password: %v", err)
		}
		user := User{
			Email:        p.Email,
			PasswordHash: string(hash),
			Role:         RolePartner,
		}
		if err := DB.Create(&user).Error; err != nil {
			log.Fatalf("failed to seed test user: %v", err)
		}

		partner := Partner{
			UserID:       user.ID,
			BusinessName: p.BusinessName,
			Siret:        p.Siret,
			Category:     p.Category,
			Address:      p.Address,
			Region:       p.Region,
			Balance:      100000,
			Status:       StatusApproved,
		}
		if err := DB.Create(&partner).Error; err != nil {
			log.Fatalf("failed to seed test partner: %v", err)
		}

		log.Printf("seeded test partner -> email: %s password: password123", p.Email)
	}
}
