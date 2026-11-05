package repository

import (
	"context"
	"fmt"

	"github.com/imamsyafii-dev/ledgerflow/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerRepository struct {
	db *pgxpool.Pool
}

func NewCustomerRepository(db *pgxpool.Pool) *CustomerRepository {
	return &CustomerRepository{
		db: db,
	}
}

func (r *CustomerRepository) Create(
	ctx context.Context,
	organizationID string,
	name string,
	email string,
) (model.Customer, error) {
	const query = `
		INSERT INTO customers (organization_id, name, email)
		VALUES ($1, $2, NULLIF($3, ''))
		RETURNING
			id::text,
			organization_id::text,
			name,
			email,
			created_at,
			updated_at
	`

	var customer model.Customer

	err := r.db.QueryRow(
		ctx,
		query,
		organizationID,
		name,
		email,
	).Scan(
		&customer.ID,
		&customer.OrganizationID,
		&customer.Name,
		&customer.Email,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)
	if err != nil {
		return model.Customer{}, fmt.Errorf("create customer: %w", err)
	}

	return customer, nil
}
