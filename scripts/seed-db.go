package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Stu076/SpendSense/config"
	"github.com/Stu076/SpendSense/internal/models"
	"github.com/Stu076/SpendSense/internal/repositories"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

// SeedUsers inserts test users into the database
func SeedUsers(db *bun.DB) error {
	users := []models.User{
		{Username: "admin", Email: "admin@example.com", Password: "adminpass", Role: "admin"},
		{Username: "user1", Email: "user1@example.com", Password: "userpass", Role: "user"},
		{Username: "user2", Email: "user2@example.com", Password: "userpass", Role: "user"},
	}

	for _, user := range users {
		hashedPass, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		user.Password = string(hashedPass)

		_, err := db.NewInsert().Model(&user).Exec(context.Background())
		if err != nil {
			return fmt.Errorf("failed to insert user %s: %w", user.Username, err)
		}
	}

	log.Println("✅ Seeding complete: Users added.")
	return nil
}

func main() {
	confStatus := config.Init()
	if !confStatus {
		panic("Error initializing config")
	}
	db := repositories.MustSetupDatabase()
	if err := SeedUsers(db); err != nil {
		log.Fatal("Seeding failed:", err)
	}
}
