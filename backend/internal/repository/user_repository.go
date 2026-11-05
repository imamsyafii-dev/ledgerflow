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

func (r *UserRepository) ListByOrganization(
	ctx context.Context,
	organizationID string,
) ([]model.User, error) {
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
		WHERE organization_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, organizationID)
	if err != nil {
		return nil, fmt.Errorf("list users by organization: %w", err)
	}
	defer rows.Close()

	users := make([]model.User, 0)

	for rows.Next() {
		var user model.User

		if err := rows.Scan(
			&user.ID,
			&user.OrganizationID,
			&user.Name,
			&user.Email,
			&user.PasswordHash,
			&user.Role,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return users, nil
}

func (r *UserRepository) FindByIDAndOrganization(
	ctx context.Context,
	userID string,
	organizationID string,
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
		WHERE id = $1 AND organization_id = $2
	`

	var user model.User

	err := r.db.QueryRow(ctx, query, userID, organizationID).Scan(
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
		return model.User{}, fmt.Errorf("find user by organization: %w", err)
	}

	return user, nil
}

func (r *UserRepository) UpdateRole(
	ctx context.Context,
	userID string,
	organizationID string,
	role model.UserRole,
) (model.User, error) {
	const query = `
		UPDATE users
		SET role = $3, updated_at = NOW()
		WHERE id = $1 AND organization_id = $2
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

	return r.scanUser(ctx, query, userID, organizationID, role)
}

func (r *UserRepository) UpdateStatus(
	ctx context.Context,
	userID string,
	organizationID string,
	isActive bool,
) (model.User, error) {
	const query = `
		UPDATE users
		SET is_active = $3, updated_at = NOW()
		WHERE id = $1 AND organization_id = $2
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

	return r.scanUser(ctx, query, userID, organizationID, isActive)
}

func (r *UserRepository) CountActiveFinanceAdmins(
	ctx context.Context,
	organizationID string,
) (int, error) {
	const query = `
		SELECT COUNT(*)
		FROM users
		WHERE organization_id = $1
			AND role = $2
			AND is_active = TRUE
	`

	var count int

	err := r.db.QueryRow(
		ctx,
		query,
		organizationID,
		model.UserRoleFinanceAdmin,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count active finance admins: %w", err)
	}

	return count, nil
}

func (r *UserRepository) scanUser(
	ctx context.Context,
	query string,
	args ...any,
) (model.User, error) {
	var user model.User

	err := r.db.QueryRow(ctx, query, args...).Scan(
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
		return model.User{}, fmt.Errorf("update user: %w", err)
	}

	return user, nil
}
