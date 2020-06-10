package easypost

import "time"

type PaymentLog struct {
	ID         string     `json:"id,omitempty"`
	Object     string     `json:"object,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	SourceType string     `json:"source_type,omitempty"`
	TargetType string     `json:"target_type,omitempty"`
	Date       string     `json:"date,omitempty"`
	ChargeType string     `json:"charge_type,omitempty"`
	Status     string     `json:"status,omitempty"`
	Amount     string     `json:"amount,omitempty"`
	Last4      string     `json:"last4,omitempty"`
}
