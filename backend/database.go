package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleClient  Role = "client"
	RolePartner Role = "partner"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"unique"`
	PasswordHash string
	Role         Role `gorm:"type:varchar(20);not null;default:'client';check:role IN ('admin','client','partner')"`
}

type PartnerStatus string

const (
	StatusPending  PartnerStatus = "pending"
	StatusApproved PartnerStatus = "approved"
	StatusRejected PartnerStatus = "rejected"
)

type Partner struct {
	ID           uint   `gorm:"primaryKey"`
	BusinessName string `gorm:"unique"`
	UserID       uint
	User         User   `gorm:"foreignKey:UserID"`
	Siret        string `gorm:"unique"`
	Category     string
	Address      string
	Region       string
	Status       PartnerStatus `gorm:"type:varchar(20);not null;default:'pending';check:status IN ('pending','approved','rejected')"`
	RejectReason *string       //pointer to string to allow null value
	CreatedAt    time.Time
}

type Client struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	User      User `gorm:"foreignKey:UserID"`
	Name string
	Balance float64
	CreatedAt time.Time
}

func database_init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Println("Database connection established")

	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}
