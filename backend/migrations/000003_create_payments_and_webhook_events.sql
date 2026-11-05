CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    invoice_id UUID NOT NULL REFERENCES invoices(id),
    provider VARCHAR(50) NOT NULL,
    provider_payment_id VARCHAR(150) NOT NULL,
    amount BIGINT NOT NULL CHECK (amount > 0),
    currency CHAR(3) NOT NULL DEFAULT 'IDR',
    status VARCHAR(20) NOT NULL
        CHECK (status IN ('PENDING', 'SUCCEEDED', 'FAILED', 'REFUNDED')),
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_payments_provider_payment_id
        UNIQUE (provider, provider_payment_id)
);

CREATE INDEX idx_payments_organization_id
    ON payments(organization_id);

CREATE INDEX idx_payments_invoice_id
    ON payments(invoice_id);

CREATE INDEX idx_payments_status
    ON payments(status);

CREATE TABLE payment_webhook_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider VARCHAR(50) NOT NULL,
    event_id VARCHAR(150) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_payment_webhook_events_provider_event_id
        UNIQUE (provider, event_id)
);

CREATE INDEX idx_payment_webhook_events_created_at
    ON payment_webhook_events(created_at);
