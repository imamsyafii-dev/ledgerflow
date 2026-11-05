ALTER TABLE refund_requests
    DROP CONSTRAINT refund_requests_status_check;

ALTER TABLE refund_requests
    ADD CONSTRAINT refund_requests_status_check
    CHECK (
        status IN (
            'PENDING',
            'APPROVED',
            'REJECTED',
            'COMPLETED',
            'FAILED'
        )
    );

ALTER TABLE refund_requests
    ADD COLUMN executed_by UUID REFERENCES users(id),
    ADD COLUMN executed_at TIMESTAMPTZ,
    ADD COLUMN provider_refund_id VARCHAR(150),
    ADD COLUMN completed_at TIMESTAMPTZ,
    ADD COLUMN failure_reason TEXT;