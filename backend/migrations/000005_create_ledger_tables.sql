CREATE TABLE ledger_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    code VARCHAR(20) NOT NULL,
    name VARCHAR(150) NOT NULL,
    account_type VARCHAR(20) NOT NULL
        CHECK (account_type IN (
            'ASSET',
            'LIABILITY',
            'EQUITY',
            'REVENUE',
            'EXPENSE'
        )),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_ledger_accounts_organization_code
        UNIQUE (organization_id, code)
);

CREATE TABLE journal_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    reference_type VARCHAR(50) NOT NULL,
    reference_id UUID NOT NULL,
    description VARCHAR(255) NOT NULL,
    entry_date DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_journal_entries_reference
        UNIQUE (organization_id, reference_type, reference_id)
);

CREATE TABLE journal_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    journal_entry_id UUID NOT NULL
        REFERENCES journal_entries(id) ON DELETE CASCADE,
    ledger_account_id UUID NOT NULL
        REFERENCES ledger_accounts(id),
    debit BIGINT NOT NULL DEFAULT 0 CHECK (debit >= 0),
    credit BIGINT NOT NULL DEFAULT 0 CHECK (credit >= 0),
    currency CHAR(3) NOT NULL DEFAULT 'IDR',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_journal_lines_single_side
        CHECK (
            (debit > 0 AND credit = 0)
            OR
            (credit > 0 AND debit = 0)
        )
);

CREATE INDEX idx_ledger_accounts_organization_id
    ON ledger_accounts(organization_id);

CREATE INDEX idx_journal_entries_organization_id
    ON journal_entries(organization_id);

CREATE INDEX idx_journal_entries_reference
    ON journal_entries(reference_type, reference_id);

CREATE INDEX idx_journal_lines_journal_entry_id
    ON journal_lines(journal_entry_id);
