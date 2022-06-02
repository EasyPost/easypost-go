package easypost

import (
	"context"
	"errors"
	"strings"
)

type CreditCard struct {
	ID              string `json:"id,omitempty"`
	DisbledAt       string `json:"disabled_at,omitempty"`
	Object          string `json:"object,omitempty"`
	Name            string `json:"name,omitempty"`
	Last4           string `json:"last4,omitempty"`
	ExpirationMonth string `json:"expiration_month,omitempty"`
	ExpirationYear  string `json:"expiration_year,omitempty"`
	Brand           string `json:"brand,omitempty"`
}

// FundCreditCard allows you to fund your EasyPost wallet by charging your primary or secondary card on file.
func (c *Client) FundCreditCard(amount string, primaryOrSecondary string) (out *CreditCard, err error) {
	return c.FundCreditCardWithContext(context.Background(), amount, primaryOrSecondary)
}

// FundCreditCardWithContext performs the same as FundCreditCard, but
// allows specifying a context that can interrupt the request.
func (c *Client) FundCreditCardWithContext(ctx context.Context, amount string, primaryOrSecondary string) (out *CreditCard, err error) {
	paymentMethod, _ := c.ListPaymentMethods()
	var cardId string

	switch primaryOrSecondary {
	case "primary":
		cardId = paymentMethod.PrimaryPaymentMethod.ID
	case "secondary":
		cardId = paymentMethod.SecondaryPaymentMethod.ID
	default:
		return out, errors.New("payment method must be either primary or secondary")
	}

	if cardId == "" || !strings.HasPrefix(cardId, "card_") {
		return out, errors.New("the chosen payment method is not a credit card. Please try again")
	}

	wrappedParams := map[string]interface{}{
		"amount": amount,
	}

	err = c.post(ctx, "credit_cards/"+cardId+"/charges", wrappedParams, &out)
	return
}

// DeleteCreditCard allows you to delete a credit card in your payment method.
func (c *Client) DeleteCreditCard(creditCardId string) (err error) {
	return c.DeleteCreditCardWithContext(context.Background(), creditCardId)
}

// DeleteCreditCardWithContext performs the same as DeleteCreditCard, but
// allows specifying a context that can interrupt the request.
func (c *Client) DeleteCreditCardWithContext(ctx context.Context, creditCardId string) (err error) {
	return c.del(ctx, "credit_cards/"+creditCardId)
}
