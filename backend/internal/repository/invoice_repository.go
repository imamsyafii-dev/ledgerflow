package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/imamsyafii-dev/ledgerflow/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InvoiceRepository struct {
	db *pgxpool.Pool
}

func NewInvoiceRepository(db *pgxpool.Pool) *InvoiceRepository {
	return &InvoiceRepository{
		db: db,
	}
}

func (r *InvoiceRepository) Create(
	ctx context.Context,
	organizationID string,
	customerID string,
	number string,
	amount int64,
	currency string,
	dueDate time.Time,
) (model.Invoice, error) {
	const query = `
		INSERT INTO invoices (
			organization_id,
			customer_id,
			number,
			amount,
			currency,
			status,
			due_date
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING
			id::text,
			organization_id::text,
			customer_id::text,
			number,
			amount,
			currency,
			status,
			due_date,
			paid_at,
			created_at,
			updated_at
	`

	var invoice model.Invoice

	err := r.db.QueryRow(
		ctx,
		query,
		organizationID,
		customerID,
		number,
		amount,
		currency,
		model.InvoiceStatusDraft,
		dueDate,
	).Scan(
		&invoice.ID,
		&invoice.OrganizationID,
		&invoice.CustomerID,
		&invoice.Number,
		&invoice.Amount,
		&invoice.Currency,
		&invoice.Status,
		&invoice.DueDate,
		&invoice.PaidAt,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
	)
	if err != nil {
		return model.Invoice{}, fmt.Errorf("create invoice: %w", err)
	}

	return invoice, nil
}
