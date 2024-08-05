package easypost

type PaymentLog struct {
	ID         string    `json:"id,omitempty" url:"id,omitempty"`
	Object     string    `json:"object,omitempty" url:"object,omitempty"`
	CreatedAt  *DateTime `json:"created_at,omitempty" url:"created_at,omitempty"`
	UpdatedAt  *DateTime `json:"updated_at,omitempty" url:"updated_at,omitempty"`
	SourceType string    `json:"source_type,omitempty" url:"source_type,omitempty"`
	TargetType string    `json:"target_type,omitempty" url:"target_type,omitempty"`
	Date       string    `json:"date,omitempty" url:"date,omitempty"`
	ChargeType string    `json:"charge_type,omitempty" url:"charge_type,omitempty"`
	Status     string    `json:"status,omitempty" url:"status,omitempty"`
	Amount     string    `json:"amount,omitempty" url:"amount,omitempty"`
	Last4      string    `json:"last4,omitempty" url:"last4,omitempty"`
}
