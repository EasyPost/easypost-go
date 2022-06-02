package easypost

import (
	"context"
	"errors"
)

type PaymentMethod struct {
	ID                     string      `json:"id,omitempty"`
	Object                 string      `json:"object,omitempty"`
	PrimaryPaymentMethod   *CreditCard `json:"primary_payment_method,omitempty"`
	SecondaryPaymentMethod *CreditCard `json:"secondary_payment_method,omitempty"`
}

// ListPaymentMethods returns the payment methods associated with the current user.
func (c *Client) ListPaymentMethods() (out *PaymentMethod, err error) {
	return c.ListPaymentMethodsWithContext(context.Background())
}

// ListPaymentMethodsWithContext performs the same as ListPaymentMethods, but
// allows specifying a context that can interrupt the request.
func (c *Client) ListPaymentMethodsWithContext(ctx context.Context) (out *PaymentMethod, err error) {
	err = c.get(ctx, "payment_methods", &out)

	if out.ID == "" {
		return out, errors.New("billing has not been setup for this user. please add a payment method")
	}

	return
}
