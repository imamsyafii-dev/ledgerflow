package model

import "time"

type JournalLine struct {
	AccountCode string `json:"account_code"`
	AccountName string `json:"account_name"`
	Debit       int64  `json:"debit"`
	Credit      int64  `json:"credit"`
	Currency    string `json:"currency"`
}

type JournalEntry struct {
	ID             string        `json:"id"`
	OrganizationID string        `json:"organization_id"`
	ReferenceType  string        `json:"reference_type"`
	ReferenceID    string        `json:"reference_id"`
	Description    string        `json:"description"`
	EntryDate      time.Time     `json:"entry_date"`
	CreatedAt      time.Time     `json:"created_at"`
	Lines          []JournalLine `json:"lines"`
}
