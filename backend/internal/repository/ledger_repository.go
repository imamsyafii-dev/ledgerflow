package repository

import (
	"context"
	"fmt"

	"github.com/imamsyafii-dev/ledgerflow/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LedgerRepository struct {
	db *pgxpool.Pool
}

func NewLedgerRepository(db *pgxpool.Pool) *LedgerRepository {
	return &LedgerRepository{
		db: db,
	}
}

func (r *LedgerRepository) ListJournalEntries(
	ctx context.Context,
	organizationID string,
) ([]model.JournalEntry, error) {
	const query = `
		SELECT
			je.id::text,
			je.organization_id::text,
			je.reference_type,
			je.reference_id::text,
			je.description,
			je.entry_date,
			je.created_at,
			la.code,
			la.name,
			jl.debit,
			jl.credit,
			jl.currency
		FROM journal_entries je
		JOIN journal_lines jl
			ON jl.journal_entry_id = je.id
		JOIN ledger_accounts la
			ON la.id = jl.ledger_account_id
		WHERE je.organization_id = $1
		ORDER BY je.created_at DESC, jl.created_at ASC
	`

	rows, err := r.db.Query(ctx, query, organizationID)
	if err != nil {
		return nil, fmt.Errorf("list journal entries: %w", err)
	}
	defer rows.Close()

	entriesByID := make(map[string]*model.JournalEntry)
	entryIDs := make([]string, 0)

	for rows.Next() {
		var (
			entry   model.JournalEntry
			line    model.JournalLine
			entryID string
		)

		err := rows.Scan(
			&entryID,
			&entry.OrganizationID,
			&entry.ReferenceType,
			&entry.ReferenceID,
			&entry.Description,
			&entry.EntryDate,
			&entry.CreatedAt,
			&line.AccountCode,
			&line.AccountName,
			&line.Debit,
			&line.Credit,
			&line.Currency,
		)
		if err != nil {
			return nil, fmt.Errorf("scan journal entry: %w", err)
		}

		existingEntry, exists := entriesByID[entryID]
		if !exists {
			entry.ID = entryID
			entry.Lines = make([]model.JournalLine, 0)

			entriesByID[entryID] = &entry
			entryIDs = append(entryIDs, entryID)
			existingEntry = &entry
		}

		existingEntry.Lines = append(existingEntry.Lines, line)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate journal entries: %w", err)
	}

	entries := make([]model.JournalEntry, 0, len(entryIDs))

	for _, entryID := range entryIDs {
		entries = append(entries, *entriesByID[entryID])
	}

	return entries, nil
}
