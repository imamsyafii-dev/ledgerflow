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

func (r *InvoiceRepository) List(
	ctx context.Context,
	organizationID string,
) ([]model.Invoice, error) {
	const query = `
		SELECT
			i.id::text,
			i.organization_id::text,
			i.customer_id::text,
			c.name,
			i.number,
			i.amount,
			i.currency,
			i.status,
			i.due_date,
			i.paid_at,
			i.created_at,
			i.updated_at
		FROM invoices i
		JOIN customers c
			ON c.id = i.customer_id
			AND c.organization_id = i.organization_id
		WHERE i.organization_id = $1
		ORDER BY i.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, organizationID)
	if err != nil {
		return nil, fmt.Errorf("list invoices: %w", err)
	}
	defer rows.Close()

	invoices := make([]model.Invoice, 0)

	for rows.Next() {
		var invoice model.Invoice

		err := rows.Scan(
			&invoice.ID,
			&invoice.OrganizationID,
			&invoice.CustomerID,
			&invoice.CustomerName,
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
			return nil, fmt.Errorf("scan invoice: %w", err)
		}

		invoices = append(invoices, invoice)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate invoices: %w", err)
	}

	return invoices, nil
}

func (r *InvoiceRepository) Issue(
	ctx context.Context,
	organizationID string,
	invoiceID string,
) (model.Invoice, error) {
	const query = `
		UPDATE invoices
		SET
			status = 'ISSUED',
			updated_at = NOW()
		WHERE
			id = $1
			AND organization_id = $2
			AND status = 'DRAFT'
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
		invoiceID,
		organizationID,
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
		return model.Invoice{}, fmt.Errorf("issue invoice: %w", err)
	}

	return invoice, nil
}
