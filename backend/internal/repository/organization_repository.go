package repository

import (
	"context"
	"fmt"

	"github.com/imamsyafii-dev/ledgerflow/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrganizationRepository struct {
	db *pgxpool.Pool
}

func NewOrganizationRepository(db *pgxpool.Pool) *OrganizationRepository {
	return &OrganizationRepository{
		db: db,
	}
}

func (r *OrganizationRepository) Create(
	ctx context.Context,
	name string,
	slug string,
) (model.Organization, error) {
	const query = `
		INSERT INTO organizations (name, slug)
		VALUES ($1, $2)
		RETURNING
			id::text,
			name,
			slug,
			created_at,
			updated_at
	`

	var organization model.Organization

	err := r.db.QueryRow(ctx, query, name, slug).Scan(
		&organization.ID,
		&organization.Name,
		&organization.Slug,
		&organization.CreatedAt,
		&organization.UpdatedAt,
	)
	if err != nil {
		return model.Organization{}, fmt.Errorf("create organization: %w", err)
	}

	return organization, nil
}
