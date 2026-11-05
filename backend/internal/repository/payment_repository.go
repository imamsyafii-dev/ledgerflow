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

	const insertReconciliationQuery = `
		INSERT INTO reconciliations (
			organization_id,
			invoice_id,
			payment_id,
			amount,
			status
		)
		VALUES ($1, $2, $3, $4, 'MATCHED')
	`

	_, err = tx.Exec(
		ctx,
		insertReconciliationQuery,
		event.OrganizationID,
		event.InvoiceID,
		payment.ID,
		event.Amount,
	)
	if err != nil {
		return model.Payment{}, false, fmt.Errorf("create reconciliation: %w", err)
	}

	const ensureLedgerAccountsQuery = `
		INSERT INTO ledger_accounts (
			organization_id,
			code,
			name,
			account_type
		)
		VALUES
			($1, '1000', 'Cash', 'ASSET'),
			($1, '1100', 'Accounts Receivable', 'ASSET')
		ON CONFLICT (organization_id, code) DO NOTHING
	`

	_, err = tx.Exec(
		ctx,
		ensureLedgerAccountsQuery,
		event.OrganizationID,
	)
	if err != nil {
		return model.Payment{}, false, fmt.Errorf("ensure ledger accounts: %w", err)
	}

	const findLedgerAccountQuery = `
		SELECT id::text
		FROM ledger_accounts
		WHERE organization_id = $1 AND code = $2
	`

	var cashAccountID string

	err = tx.QueryRow(
		ctx,
		findLedgerAccountQuery,
		event.OrganizationID,
		"1000",
	).Scan(&cashAccountID)
	if err != nil {
		return model.Payment{}, false, fmt.Errorf("find cash account: %w", err)
	}

	var receivableAccountID string

	err = tx.QueryRow(
		ctx,
		findLedgerAccountQuery,
		event.OrganizationID,
		"1100",
	).Scan(&receivableAccountID)
	if err != nil {
		return model.Payment{}, false, fmt.Errorf("find accounts receivable account: %w", err)
	}

	const insertJournalEntryQuery = `
		INSERT INTO journal_entries (
			organization_id,
			reference_type,
			reference_id,
			description
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id::text
	`

	var journalEntryID string

	err = tx.QueryRow(
		ctx,
		insertJournalEntryQuery,
		event.OrganizationID,
		"PAYMENT",
		payment.ID,
		fmt.Sprintf("Payment received for invoice %s", event.InvoiceID),
	).Scan(&journalEntryID)
	if err != nil {
		return model.Payment{}, false, fmt.Errorf("create journal entry: %w", err)
	}

	const insertJournalLinesQuery = `
		INSERT INTO journal_lines (
			journal_entry_id,
			ledger_account_id,
			debit,
			credit,
			currency
		)
		VALUES
			($1, $2, $3, 0, $4),
			($1, $5, 0, $3, $4)
	`

	_, err = tx.Exec(
		ctx,
		insertJournalLinesQuery,
		journalEntryID,
		cashAccountID,
		event.Amount,
		event.Currency,
		receivableAccountID,
	)
	if err != nil {
		return model.Payment{}, false, fmt.Errorf("create journal lines: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Payment{}, false, fmt.Errorf("commit payment transaction: %w", err)
	}

	return payment, false, nil
}
