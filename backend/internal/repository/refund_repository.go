package repository

import (
	"context"
	"fmt"

	"github.com/imamsyafii-dev/ledgerflow/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefundRepository struct {
	db *pgxpool.Pool
}

func NewRefundRepository(db *pgxpool.Pool) *RefundRepository {
	return &RefundRepository{
		db: db,
	}
}

func (r *RefundRepository) Create(
	ctx context.Context,
	organizationID string,
	paymentID string,
	requestedBy string,
	reason string,
) (model.RefundRequest, error) {
	const query = `
		INSERT INTO refund_requests (
			organization_id,
			payment_id,
			requested_by,
			reason
		)
		SELECT
			$1,
			id,
			$3,
			$4
		FROM payments
		WHERE id = $2
			AND organization_id = $1
			AND status = 'SUCCEEDED'
		RETURNING
			id::text,
			organization_id::text,
			payment_id::text,
			requested_by::text,
			reason,
			status,
			created_at
	`

	var refund model.RefundRequest

	err := r.db.QueryRow(
		ctx,
		query,
		organizationID,
		paymentID,
		requestedBy,
		reason,
	).Scan(
		&refund.ID,
		&refund.OrganizationID,
		&refund.PaymentID,
		&refund.RequestedBy,
		&refund.Reason,
		&refund.Status,
		&refund.CreatedAt,
	)

	return refund, err
}

func (r *RefundRepository) List(
	ctx context.Context,
	organizationID string,
) ([]model.RefundRequest, error) {
	const query = `
		SELECT
			r.id::text,
			r.organization_id::text,
			r.payment_id::text,
			i.number,
			c.name,
			p.amount,
			p.currency,
			r.requested_by::text,
			u.name,
			r.reason,
			r.status,
			r.decided_by::text,
			r.decision_note,
			r.decided_at,
			r.created_at
		FROM refund_requests r
		JOIN payments p ON p.id = r.payment_id
		JOIN invoices i ON i.id = p.invoice_id
		JOIN customers c ON c.id = i.customer_id
		JOIN users u ON u.id = r.requested_by
		WHERE r.organization_id = $1
		ORDER BY r.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, organizationID)
	if err != nil {
		return nil, fmt.Errorf("list refund requests: %w", err)
	}
	defer rows.Close()

	refunds := make([]model.RefundRequest, 0)

	for rows.Next() {
		var refund model.RefundRequest

		if err := rows.Scan(
			&refund.ID,
			&refund.OrganizationID,
			&refund.PaymentID,
			&refund.InvoiceNumber,
			&refund.CustomerName,
			&refund.Amount,
			&refund.Currency,
			&refund.RequestedBy,
			&refund.RequestedByName,
			&refund.Reason,
			&refund.Status,
			&refund.DecidedBy,
			&refund.DecisionNote,
			&refund.DecidedAt,
			&refund.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan refund request: %w", err)
		}

		refunds = append(refunds, refund)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate refund requests: %w", err)
	}

	return refunds, nil
}

func (r *RefundRepository) Decide(
	ctx context.Context,
	organizationID string,
	refundID string,
	actorID string,
	status model.RefundStatus,
	note string,
) (model.RefundRequest, error) {
	const query = `
		UPDATE refund_requests
		SET
			status = $4,
			decided_by = $3,
			decision_note = $5,
			decided_at = NOW(),
			updated_at = NOW()
		WHERE id = $2
			AND organization_id = $1
			AND status = 'PENDING'
			AND requested_by <> $3
		RETURNING
			id::text,
			organization_id::text,
			payment_id::text,
			requested_by::text,
			reason,
			status,
			decided_by::text,
			decision_note,
			decided_at,
			created_at
	`

	var refund model.RefundRequest

	err := r.db.QueryRow(
		ctx,
		query,
		organizationID,
		refundID,
		actorID,
		status,
		note,
	).Scan(
		&refund.ID,
		&refund.OrganizationID,
		&refund.PaymentID,
		&refund.RequestedBy,
		&refund.Reason,
		&refund.Status,
		&refund.DecidedBy,
		&refund.DecisionNote,
		&refund.DecidedAt,
		&refund.CreatedAt,
	)

	return refund, err
}

func (r *RefundRepository) Execute(
	ctx context.Context,
	organizationID string,
	refundID string,
	executedBy string,
) (model.RefundRequest, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return model.RefundRequest{}, fmt.Errorf(
			"begin refund execution transaction: %w",
			err,
		)
	}
	defer tx.Rollback(ctx)

	const lockRefundQuery = `
		SELECT
			payment_id::text,
			reason
		FROM refund_requests
		WHERE id = $1
			AND organization_id = $2
			AND status = 'APPROVED'
		FOR UPDATE
	`

	var paymentID string
	var reason string

	err = tx.QueryRow(
		ctx,
		lockRefundQuery,
		refundID,
		organizationID,
	).Scan(
		&paymentID,
		&reason,
	)
	if err != nil {
		return model.RefundRequest{}, fmt.Errorf(
			"lock approved refund request: %w",
			err,
		)
	}

	const updatePaymentQuery = `
		UPDATE payments
		SET
			status = 'REFUNDED',
			updated_at = NOW()
		WHERE id = $1
			AND organization_id = $2
			AND status = 'SUCCEEDED'
		RETURNING
			amount,
			currency
	`

	var amount int64
	var currency string

	err = tx.QueryRow(
		ctx,
		updatePaymentQuery,
		paymentID,
		organizationID,
	).Scan(
		&amount,
		&currency,
	)
	if err != nil {
		return model.RefundRequest{}, fmt.Errorf(
			"mark payment as refunded: %w",
			err,
		)
	}

	const findLedgerAccountQuery = `
		SELECT id::text
		FROM ledger_accounts
		WHERE organization_id = $1
			AND code = $2
	`

	var cashAccountID string
	var receivableAccountID string

	err = tx.QueryRow(
		ctx,
		findLedgerAccountQuery,
		organizationID,
		"1000",
	).Scan(&cashAccountID)
	if err != nil {
		return model.RefundRequest{}, fmt.Errorf(
			"find cash account: %w",
			err,
		)
	}

	err = tx.QueryRow(
		ctx,
		findLedgerAccountQuery,
		organizationID,
		"1100",
	).Scan(&receivableAccountID)
	if err != nil {
		return model.RefundRequest{}, fmt.Errorf(
			"find accounts receivable account: %w",
			err,
		)
	}

	const insertJournalEntryQuery = `
		INSERT INTO journal_entries (
			organization_id,
			reference_type,
			reference_id,
			description
		)
		VALUES ($1, 'REFUND', $2, $3)
		RETURNING id::text
	`

	var journalEntryID string

	err = tx.QueryRow(
		ctx,
		insertJournalEntryQuery,
		organizationID,
		refundID,
		fmt.Sprintf("Refund executed: %s", reason),
	).Scan(&journalEntryID)
	if err != nil {
		return model.RefundRequest{}, fmt.Errorf(
			"create refund journal entry: %w",
			err,
		)
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
		receivableAccountID,
		amount,
		currency,
		cashAccountID,
	)
	if err != nil {
		return model.RefundRequest{}, fmt.Errorf(
			"create refund journal lines: %w",
			err,
		)
	}

	const updateRefundQuery = `
		UPDATE refund_requests
		SET
			status = 'COMPLETED',
			executed_by = $3,
			executed_at = NOW(),
			provider_refund_id = CONCAT(
				'refund_mock_',
				REPLACE(gen_random_uuid()::text, '-', '')
			),
			completed_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
			AND organization_id = $2
		RETURNING
			id::text,
			organization_id::text,
			payment_id::text,
			requested_by::text,
			reason,
			status,
			decided_by::text,
			decision_note,
			decided_at,
			executed_by::text,
			executed_at,
			provider_refund_id,
			completed_at,
			failure_reason,
			created_at
	`

	var refund model.RefundRequest

	err = tx.QueryRow(
		ctx,
		updateRefundQuery,
		refundID,
		organizationID,
		executedBy,
	).Scan(
		&refund.ID,
		&refund.OrganizationID,
		&refund.PaymentID,
		&refund.RequestedBy,
		&refund.Reason,
		&refund.Status,
		&refund.DecidedBy,
		&refund.DecisionNote,
		&refund.DecidedAt,
		&refund.ExecutedBy,
		&refund.ExecutedAt,
		&refund.ProviderRefundID,
		&refund.CompletedAt,
		&refund.FailureReason,
		&refund.CreatedAt,
	)
	if err != nil {
		return model.RefundRequest{}, fmt.Errorf(
			"complete refund request: %w",
			err,
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return model.RefundRequest{}, fmt.Errorf(
			"commit refund execution transaction: %w",
			err,
		)
	}

	return refund, nil
}
