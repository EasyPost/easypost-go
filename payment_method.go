package easypost

import (
	"errors"
	"strings"
)

type PaymentMethodPriority int64

const (
	PrimaryPaymentMethodPriority PaymentMethodPriority = iota
	SecondaryPaymentMethodPriority
)

type PaymentMethodType int64

const (
	CreditCardPaymentType PaymentMethodType = iota
	BankAccountPaymentType
)

type PaymentMethod struct {
	ID                     string               `json:"id,omitempty"`
	Object                 string               `json:"object,omitempty"`
	PrimaryPaymentMethod   *PaymentMethodObject `json:"primary_payment_method,omitempty"`
	SecondaryPaymentMethod *PaymentMethodObject `json:"secondary_payment_method,omitempty"`
}

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

// getPaymentMethodObjectType returns the PaymentMethodType enum of a PaymentMethodObject.
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

// getPaymentMethodEndpoint returns the associated endpoint for a PaymentMethodObject based on its type.
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

// getPaymentMethodByPriority returns the PaymentMethodObject associated with the given PaymentMethodPriority.
func (c *Client) getPaymentMethodByPriority(priority PaymentMethodPriority) (out *PaymentMethodObject, err error) {
	paymentMethods, err := c.RetrievePaymentMethods()
	if err != nil {
		return out, err
	}

	switch priority {
	case PrimaryPaymentMethodPriority:
		out = paymentMethods.PrimaryPaymentMethod
	case SecondaryPaymentMethodPriority:
		out = paymentMethods.SecondaryPaymentMethod
	}

	if out == nil {
		return out, errors.New("the chosen payment method is not set up yet")
	}

	return out, nil
}
