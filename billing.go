package easypost

import (
	"context"
	"net/http"
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
	return c.do(ctx, http.MethodDelete, endpoint+"/"+paymentMethodObject.ID, nil, nil)
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
	err = c.do(ctx, http.MethodPost, endpoint+"/"+paymentMethod.ID+"/charges", wrappedParams, &out)
	return err
}

// RetrievePaymentMethods returns the payment methods associated with the current user.
func (c *Client) RetrievePaymentMethods() (out *PaymentMethod, err error) {
	return c.RetrievePaymentMethodsWithContext(context.Background())
}

// RetrievePaymentMethodsWithContext performs the same as RetrievePaymentMethods, but
// allows specifying a context that can interrupt the request.
func (c *Client) RetrievePaymentMethodsWithContext(ctx context.Context) (out *PaymentMethod, err error) {
	err = c.do(ctx, http.MethodGet, "payment_methods", nil, &out)

	if out.ID == "" {
		return out, newInvalidObjectError(NoPaymentMethods)
	}

	return
}
