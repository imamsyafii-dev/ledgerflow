package repository

import (
	"context"
	"fmt"

	"github.com/imamsyafii-dev/ledgerflow/backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

func (r *PaymentRepository) ProcessSucceededWebhook(
	ctx context.Context,
	event model.PaymentSucceededWebhook,
	payload []byte,
) (model.Payment, bool, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return model.Payment{}, false, fmt.Errorf("begin payment transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	const insertEventQuery = `
		INSERT INTO payment_webhook_events (
			provider,
			event_id,
			event_type,
			payload,
			processed_at
		)
		VALUES ($1, $2, $3, $4::jsonb, NOW())
		ON CONFLICT (provider, event_id) DO NOTHING
		RETURNING id::text
	`

	var webhookEventID string

	err = tx.QueryRow(
		ctx,
		insertEventQuery,
		"mock_gateway",
		event.EventID,
		event.EventType,
		string(payload),
	).Scan(&webhookEventID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.Payment{}, true, nil
		}

		return model.Payment{}, false, fmt.Errorf("store webhook event: %w", err)
	}

	const markInvoicePaidQuery = `
		UPDATE invoices
		SET
			status = 'PAID',
			paid_at = NOW(),
			updated_at = NOW()
		WHERE
			id = $1
			AND organization_id = $2
			AND status = 'ISSUED'
			AND amount = $3
			AND currency = $4
		RETURNING id::text
	`

	var invoiceID string

	err = tx.QueryRow(
		ctx,
		markInvoicePaidQuery,
		event.InvoiceID,
		event.OrganizationID,
		event.Amount,
		event.Currency,
	).Scan(&invoiceID)
	if err != nil {
		return model.Payment{}, false, fmt.Errorf("mark invoice as paid: %w", err)
	}

	const insertPaymentQuery = `
		INSERT INTO payments (
			organization_id,
			invoice_id,
			provider,
			provider_payment_id,
			amount,
			currency,
			status,
			paid_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		RETURNING
			id::text,
			organization_id::text,
			invoice_id::text,
			provider,
			provider_payment_id,
			amount,
			currency,
			status,
			paid_at,
			created_at,
			updated_at
	`

	var payment model.Payment

	err = tx.QueryRow(
		ctx,
		insertPaymentQuery,
		event.OrganizationID,
		event.InvoiceID,
		"mock_gateway",
		event.ProviderPaymentID,
		event.Amount,
		event.Currency,
		model.PaymentStatusSucceeded,
	).Scan(
		&payment.ID,
		&payment.OrganizationID,
		&payment.InvoiceID,
		&payment.Provider,
		&payment.ProviderPaymentID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.PaidAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		return model.Payment{}, false, fmt.Errorf("create payment: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Payment{}, false, fmt.Errorf("commit payment transaction: %w", err)
	}

	return payment, false, nil
}
