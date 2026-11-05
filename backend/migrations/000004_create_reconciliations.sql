CREATE TABLE reconciliations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    invoice_id UUID NOT NULL REFERENCES invoices(id),
    payment_id UUID NOT NULL REFERENCES payments(id),
    amount BIGINT NOT NULL CHECK (amount > 0),
    status VARCHAR(20) NOT NULL
        CHECK (status IN ('MATCHED', 'MANUAL_REVIEW', 'REVERSED')),
    reconciled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_reconciliations_payment_id
        UNIQUE (payment_id)
);

CREATE INDEX idx_reconciliations_organization_id
    ON reconciliations(organization_id);

CREATE INDEX idx_reconciliations_invoice_id
    ON reconciliations(invoice_id);
