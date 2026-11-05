package repository

import (
	"context"
	"fmt"

	"github.com/imamsyafii-dev/ledgerflow/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(
	ctx context.Context,
	organizationID string,
	name string,
	email string,
	passwordHash string,
	role model.UserRole,
) (model.User, error) {
	const query = `
		INSERT INTO users (
			organization_id,
			name,
			email,
			password_hash,
			role
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			id::text,
			organization_id::text,
			name,
			email,
			password_hash,
			role,
			is_active,
			created_at,
			updated_at
	`

	var user model.User

	err := r.db.QueryRow(
		ctx,
		query,
		organizationID,
		name,
		email,
		passwordHash,
		role,
	).Scan(
		&user.ID,
		&user.OrganizationID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return model.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (model.User, error) {
	const query = `
		SELECT
			id::text,
			organization_id::text,
			name,
			email,
			password_hash,
			role,
			is_active,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
	`

	var user model.User

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.OrganizationID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return model.User{}, fmt.Errorf("find user by email: %w", err)
	}

	return user, nil
}

func (r *UserRepository) FindByID(
	ctx context.Context,
	userID string,
) (model.User, error) {
	const query = `
		SELECT
			id::text,
			organization_id::text,
			name,
			email,
			password_hash,
			role,
			is_active,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	var user model.User

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.OrganizationID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return model.User{}, fmt.Errorf("find user by ID: %w", err)
	}

	return user, nil
}
