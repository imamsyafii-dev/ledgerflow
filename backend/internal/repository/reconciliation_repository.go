package repository

import (
	"context"
	"fmt"

	"github.com/imamsyafii-dev/ledgerflow/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReconciliationRepository struct {
	db *pgxpool.Pool
}

func NewReconciliationRepository(
	db *pgxpool.Pool,
) *ReconciliationRepository {
	return &ReconciliationRepository{
		db: db,
	}
}

func (r *ReconciliationRepository) List(
	ctx context.Context,
	organizationID string,
) ([]model.Reconciliation, error) {
	const query = `
		SELECT
			r.id::text,
			r.organization_id::text,
			r.invoice_id::text,
			i.number,
			c.name,
			r.payment_id::text,
			p.provider_payment_id,
			r.amount,
			r.status,
			r.reconciled_at,
			r.created_at
		FROM reconciliations r
		JOIN invoices i
			ON i.id = r.invoice_id
		JOIN customers c
			ON c.id = i.customer_id
		JOIN payments p
			ON p.id = r.payment_id
		WHERE r.organization_id = $1
		ORDER BY r.reconciled_at DESC
	`

	rows, err := r.db.Query(ctx, query, organizationID)
	if err != nil {
		return nil, fmt.Errorf("list reconciliations: %w", err)
	}
	defer rows.Close()

	reconciliations := make([]model.Reconciliation, 0)

	for rows.Next() {
		var reconciliation model.Reconciliation

		err := rows.Scan(
			&reconciliation.ID,
			&reconciliation.OrganizationID,
			&reconciliation.InvoiceID,
			&reconciliation.InvoiceNumber,
			&reconciliation.CustomerName,
			&reconciliation.PaymentID,
			&reconciliation.ProviderPaymentID,
			&reconciliation.Amount,
			&reconciliation.Status,
			&reconciliation.ReconciledAt,
			&reconciliation.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan reconciliation: %w", err)
		}

		reconciliations = append(reconciliations, reconciliation)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate reconciliations: %w", err)
	}

	return reconciliations, nil
}
