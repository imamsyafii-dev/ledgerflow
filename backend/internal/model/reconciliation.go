package model

import "time"

type Reconciliation struct {
	ID                string    `json:"id"`
	OrganizationID    string    `json:"organization_id"`
	InvoiceID         string    `json:"invoice_id"`
	InvoiceNumber     string    `json:"invoice_number"`
	CustomerName      string    `json:"customer_name"`
	PaymentID         string    `json:"payment_id"`
	ProviderPaymentID string    `json:"provider_payment_id"`
	Amount            int64     `json:"amount"`
	Status            string    `json:"status"`
	ReconciledAt      time.Time `json:"reconciled_at"`
	CreatedAt         time.Time `json:"created_at"`
}
