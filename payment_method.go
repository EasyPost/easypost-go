package easypost

type PaymentMethod struct {
	ID                     string               `json:"id,omitempty"`
	Object                 string               `json:"object,omitempty"`
	PrimaryPaymentMethod   *PaymentMethodObject `json:"primary_payment_method,omitempty"`
	SecondaryPaymentMethod *PaymentMethodObject `json:"secondary_payment_method,omitempty"`
}

type PaymentMethodPriority int64

const (
	PrimaryPaymentMethodPriority PaymentMethodPriority = iota
	SecondaryPaymentMethodPriority
)
