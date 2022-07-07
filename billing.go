package easypost

import (
	"context"
	"errors"
)

// DeletePaymentMethod allows you to delete a payment method in your wallet.
func (c *Client) DeletePaymentMethod(priority PaymentMethodPriority) (err error) {
	return c.DeletePaymentMethodWithContext(context.Background(), priority)
}

// DeletePaymentMethodWithContext performs the same as DeletePaymentMethod, but
// allows specifying a context that can interrupt the request.
func (c *Client) DeletePaymentMethodWithContext(ctx context.Context, priority PaymentMethodPriority) (err error) {
	paymentMethodObject, _ := c.getPaymentMethodByPriority(priority)
	endpoint, _ := c.getPaymentMethodEndpoint(paymentMethodObject)
	return c.del(ctx, endpoint+"/"+paymentMethodObject.ID)
}

// FundWallet allows you to fund your EasyPost wallet by charging your primary or secondary payment method on file.
func (c *Client) FundWallet(amount string, priority PaymentMethodPriority) (err error) {
	return c.FundWalletWithContext(context.Background(), amount, priority)
}

// FundWalletWithContext performs the same as FundWallet, but
// allows specifying a context that can interrupt the request.
func (c *Client) FundWalletWithContext(ctx context.Context, amount string, priority PaymentMethodPriority) (err error) {
	paymentMethod, _ := c.getPaymentMethodByPriority(priority)
	endpoint, _ := c.getPaymentMethodEndpoint(paymentMethod)

	wrappedParams := map[string]interface{}{
		"amount": amount,
	}

	out := PaymentMethodObject{}
	err = c.post(ctx, endpoint+"/"+paymentMethod.ID+"/charges", wrappedParams, &out)
	return err
}

// RetrievePaymentMethods returns the payment methods associated with the current user.
func (c *Client) RetrievePaymentMethods() (out *PaymentMethod, err error) {
	return c.RetrievePaymentMethodsWithContext(context.Background())
}

// RetrievePaymentMethodsWithContext performs the same as RetrievePaymentMethods, but
// allows specifying a context that can interrupt the request.
func (c *Client) RetrievePaymentMethodsWithContext(ctx context.Context) (out *PaymentMethod, err error) {
	err = c.get(ctx, "payment_methods", &out)

	if out.ID == "" {
		return out, errors.New("billing has not been setup for this user. please add a payment method")
	}

	return
}

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
