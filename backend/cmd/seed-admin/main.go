package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/imamsyafii-dev/ledgerflow/backend/internal/auth"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/config"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/database"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/model"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/repository"
)

func main() {
	organizationID := strings.TrimSpace(os.Getenv("ADMIN_ORGANIZATION_ID"))
	name := strings.TrimSpace(os.Getenv("ADMIN_NAME"))
	email := strings.TrimSpace(os.Getenv("ADMIN_EMAIL"))
	password := os.Getenv("ADMIN_PASSWORD")

	if organizationID == "" || name == "" || email == "" || password == "" {
		log.Fatal(
			"ADMIN_ORGANIZATION_ID, ADMIN_NAME, ADMIN_EMAIL, and ADMIN_PASSWORD are required",
		)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	db, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		log.Fatal(fmt.Errorf("hash password: %w", err))
	}

	userRepository := repository.NewUserRepository(db)

	user, err := userRepository.Create(
		ctx,
		organizationID,
		name,
		strings.ToLower(email),
		passwordHash,
		model.UserRoleFinanceAdmin,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(
		"admin created: id=%s email=%s role=%s",
		user.ID,
		user.Email,
		user.Role,
	)
}
