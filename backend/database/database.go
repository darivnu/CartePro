//
// EPITECH PROJECT, 2026
// CartePro
// File description:
// database
//

package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

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

type Session struct {
	ID        uint   `gorm:"primaryKey"`
	Token     string `gorm:"unique;not null"`
	UserID    uint
	User      User `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	ExpiresAt time.Time
}

type Admin struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	User      User `gorm:"foreignKey:UserID"`
	Name      string
	CreatedAt time.Time
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
	Balance      int64         `gorm:"default:0"` //stored in cents
	Status       PartnerStatus `gorm:"type:varchar(20);not null;default:'pending';check:status IN ('pending','approved','rejected')"`
	RejectReason *string       //pointer to string to allow null value
	CreatedAt    time.Time
}

type Client struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint
	User       User `gorm:"foreignKey:UserID"`
	EmployerID *uint
	Employer   *Employer `gorm:"foreignKey:EmployerID"`
	Name       string
	Balance    int64 `gorm:"default:0"` //stored in cents
	CreatedAt  time.Time
}

// the employer is not a real user rn, and is just a placeholder. the admin manages all employers.
// it is however the employer that tops up the client balance (through the admin panel).
// therefore the transaction table will have an optional foreign key to the employer should it be a client topup transaction.
type Employer struct {
	ID        uint     `gorm:"primaryKey"`
	Name      string   `gorm:"unique"`
	Clients   []Client `gorm:"foreignKey:EmployerId"`
	CreatedAt time.Time
}
type QrToken struct {
	ID        uint `gorm:"primaryKey"`
	ClientID  uint
	Client    Client `gorm:"foreignKey:ClientID"`
	Token     string `gorm:"unique"`
	CreatedAt time.Time
	ExpiresAt time.Time
	UsedAt    *time.Time // pointer to time.Time to allow null value
}

type TransactionType string

const (
	TransactionTypeDebit TransactionType = "debit"
	TransactionTypeTopup TransactionType = "topup"
)

type Transaction struct {
	ID         uint `gorm:"primaryKey"`
	ClientID   uint
	Client     Client `gorm:"foreignKey:ClientID"`
	PartnerID  *uint
	Partner    *Partner `gorm:"foreignKey:PartnerID"`
	EmployerID *uint
	Employer   *Employer `gorm:"foreignKey:EmployerID"`
	QrTokenID  *uint
	QrToken    *QrToken `gorm:"foreignKey:QrTokenID"`
	Amount     int64
	CreatedAt  time.Time
	Type       TransactionType `gorm:"type:varchar(20);not null;check:type IN ('debit','topup')"`
}

func InitDatabase() {
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
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Println("Database connection established")

	if err := DB.AutoMigrate(&User{}, &Session{}, &Partner{}, &Client{}, &QrToken{}, &Transaction{}, &Employer{}, &Admin{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	log.Println("Database migration completed")
}
