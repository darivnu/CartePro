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

func SeedTestPartners() {
	var count int64
	DB.Model(&Partner{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded, skipping test partner creation")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash seed password: %v", err)
	}

	partners := []struct {
		Email        string
		BusinessName string
		Siret        string
		Category     string
		Address      string
		Region       string
	}{
		{
			Email:        "testpartner@cartepro.dev",
			BusinessName: "Test Partner",
			Siret:        "12345678901234",
			Category:     "Test Category",
			Address:      "123 Test St, Test City",
			Region:       "Test Region",
		},
		{
			Email:        "contact@lecroissantdore.fr",
			BusinessName: "Le Croissant Doré",
			Siret:        "23456789012345",
			Category:     "Bakery",
			Address:      "12 Rue de la Paix, Lyon",
			Region:       "Auvergne-Rhône-Alpes",
		},
		{
			Email:        "hello@pixelburgerco.com",
			BusinessName: "Pixel Burger Co.",
			Siret:        "34567890123456",
			Category:     "Restaurant",
			Address:      "45 Avenue des Frites, Bordeaux",
			Region:       "Nouvelle-Aquitaine",
		},
		{
			Email:        "info@fitzonegymclub.fr",
			BusinessName: "FitZone Gym Club",
			Siret:        "45678901234567",
			Category:     "Sports & Fitness",
			Address:      "78 Boulevard du Muscle, Marseille",
			Region:       "Provence-Alpes-Côte d'Azur",
		},
	}

	for _, p := range partners {
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
