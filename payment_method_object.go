package easypost

import (
	"errors"
	"strings"
)

type PaymentMethodObject struct {
	BankName        string `json:"bank_name,omitempty"`        // bank account
	Brand           string `json:"brand,omitempty"`            // credit card
	Country         string `json:"country,omitempty"`          // bank account
	DisabledAt      string `json:"disabled_at,omitempty"`      // both
	ExpirationMonth string `json:"expiration_month,omitempty"` // credit card
	ExpirationYear  string `json:"expiration_year,omitempty"`  // credit card
	ID              string `json:"id,omitempty"`               // both
	Last4           string `json:"last4,omitempty"`            // both
	Name            string `json:"name,omitempty"`             // credit card
	Object          string `json:"object,omitempty"`           // both
	Verified        bool   `json:"verified,omitempty"`         // bank account
}

type PaymentMethodType int64

const (
	CreditCardPaymentType PaymentMethodType = iota
	BankAccountPaymentType
)

func (c *Client) getPaymentMethodObjectType(object *PaymentMethodObject) (out PaymentMethodType, err error) {
	if strings.HasPrefix(object.ID, "card_") {
		out = CreditCardPaymentType
	} else if strings.HasPrefix(object.ID, "bank_") {
		out = BankAccountPaymentType
	} else {
		return out, errors.New("no matching payment method type found")
	}
	return
}

func (c *Client) getPaymentMethodEndpoint(object *PaymentMethodObject) (out string, err error) {
	paymentMethodType, _ := c.getPaymentMethodObjectType(object)
	switch paymentMethodType {
	case CreditCardPaymentType:
		out = "credit_cards"
	case BankAccountPaymentType:
		out = "bank_accounts"
	default:
		return out, errors.New("no matching payment method type found")
	}
	return
}
