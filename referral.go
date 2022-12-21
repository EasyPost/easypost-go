package easypost

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// A ReferralCustomer contains data about an EasyPost referral customer.
type ReferralCustomer struct {
	User
}

// ListReferralCustomersResult holds the results from the list referral customers API call.
type ListReferralCustomersResult struct {
	ReferralCustomers []*ReferralCustomer `json:"referral_customers,omitempty"`
	// HasMore indicates if there are more responses to be fetched. If True,
	// additional responses can be fetched by updating the ListOptions
	// parameter's AfterID field with the ID of the last item in this object's
	// ReferralCustomers field.
	HasMore bool `json:"has_more,omitempty"`
}

// CreditCardOptions specifies options for creating or updating a credit card.
type CreditCardOptions struct {
	Number   string `json:"number,omitempty"`
	ExpMonth string `json:"expiration_month,omitempty"`
	ExpYear  string `json:"expiration_year,omitempty"`
	Cvc      string `json:"cvc,omitempty"`
}

type stripeApiKeyResponse struct {
	PublicKey string `json:"public_key"`
}

type stripeTokenResponse struct {
	Id string `json:"id"`
}

type referralCustomerRequest struct {
	UserOptions *UserOptions `json:"user,omitempty"`
}

type creditCardCreateRequest struct {
	CreditCard *easypostCreditCardCreateOptions `json:"credit_card,omitempty"`
}

type easypostCreditCardCreateOptions struct {
	StripeToken string `json:"stripe_object_id"`
	Priority    string `json:"priority"`
}

// CreateReferralCustomer creates a new referral customer.
func (c *Client) CreateReferralCustomer(in *UserOptions) (out *ReferralCustomer, err error) {
	return c.CreateReferralCustomerWithContext(context.Background(), in)
}

// CreateReferralCustomerWithContext performs the same operation as CreateReferralCustomer, but allows
// specifying a context that can interrupt the request.
func (c *Client) CreateReferralCustomerWithContext(ctx context.Context, in *UserOptions) (out *ReferralCustomer, err error) {
	err = c.post(ctx, "referral_customers", &referralCustomerRequest{UserOptions: in}, &out)
	return
}

// ListReferralCustomers provides a paginated result of ReferralCustomer objects.
func (c *Client) ListReferralCustomers(opts *ListOptions) (out *ListReferralCustomersResult, err error) {
	return c.ListReferralCustomersWithContext(context.Background(), opts)
}

// ListReferralCustomersWithContext performs the same operation as ListReferralCustomers, but
// allows specifying a context that can interrupt the request.
func (c *Client) ListReferralCustomersWithContext(ctx context.Context, opts *ListOptions) (out *ListReferralCustomersResult, err error) {
	err = c.do(ctx, http.MethodGet, "referral_customers", c.convertOptsToURLValues(opts), &out)
	return
}

// UpdateReferralCustomerEmail updates a ReferralCustomer's email address
func (c *Client) UpdateReferralCustomerEmail(userId string, email string) (out *ReferralCustomer, err error) {
	return c.UpdateReferralCustomerEmailWithContext(context.Background(), userId, email)
}

// UpdateReferralCustomerEmailWithContext performs the same operation as UpdateReferralCustomerEmail, but allows
// specifying a context that can interrupt the request.
func (c *Client) UpdateReferralCustomerEmailWithContext(ctx context.Context, userId string, email string) (out *ReferralCustomer, err error) {
	req := referralCustomerRequest{
		UserOptions: &UserOptions{
			Email: &email,
		},
	}
	err = c.put(ctx, "referral_customers/"+userId, req, &out)
	return
}

// AddReferralCustomerCreditCard adds a credit card to a referral customer's account.
func (c *Client) AddReferralCustomerCreditCard(referralCustomerApiKey string, creditCardOptions *CreditCardOptions, priority PaymentMethodPriority) (out *PaymentMethodObject, err error) {
	return c.AddReferralCustomerCreditCardWithContext(context.Background(), referralCustomerApiKey, creditCardOptions, priority)
}

// AddReferralCustomerCreditCardWithContext performs the same operation as AddReferralCustomerCreditCard, but allows
// specifying a context that can interrupt the request.
func (c *Client) AddReferralCustomerCreditCardWithContext(ctx context.Context, referralCustomerApiKey string, creditCardOptions *CreditCardOptions, priority PaymentMethodPriority) (out *PaymentMethodObject, err error) {
	stripeApiKeyResponse, err := c.retrieveEasypostStripeApiKey(ctx)
	if err != nil || stripeApiKeyResponse == nil || stripeApiKeyResponse.PublicKey == "" {
		return nil, errors.New("could not create Stripe token, please try again later")
	}

	stripeTokenResponse, err := c.createStripeToken(ctx, stripeApiKeyResponse.PublicKey, creditCardOptions)
	if err != nil || stripeTokenResponse == nil || stripeTokenResponse.Id == "" {
		return nil, errors.New("could not send card details to Stripe, please try again later")
	}

	return c.createEasypostCreditCard(ctx, referralCustomerApiKey, stripeTokenResponse.Id, priority)
}

func (c *Client) retrieveEasypostStripeApiKey(ctx context.Context) (out *stripeApiKeyResponse, err error) {
	err = c.get(ctx, "partners/stripe_public_key", &out)
	return
}

func (c *Client) createStripeToken(ctx context.Context, stripeApiKey string, creditCardOptions *CreditCardOptions) (out *stripeTokenResponse, err error) {
	data := url.Values{}
	data.Set("card[number]", creditCardOptions.Number)
	data.Set("card[exp_month]", creditCardOptions.ExpMonth)
	data.Set("card[exp_year]", creditCardOptions.ExpYear)
	data.Set("card[cvc]", creditCardOptions.Cvc)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.stripe.com/v1/tokens", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Bearer "+stripeApiKey)

	resp, err := c.client().Do(req) // use the current client's inner http.Client (configured to record) for the one-off request
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body) // deprecated, but we have to keep it for legacy compatibility
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, err
	}

	return
}

func (c *Client) createEasypostCreditCard(ctx context.Context, referralCustomerApiKey string, stripeToken string, priority PaymentMethodPriority) (out *PaymentMethodObject, err error) {
	client := &Client{
		APIKey: referralCustomerApiKey,
		Client: c.client(), // pass the current client's inner http.Client (configured to record) to the new client
	}

	priorityString := ""
	switch priority {
	case PrimaryPaymentMethodPriority:
		priorityString = "primary"
	case SecondaryPaymentMethodPriority:
		priorityString = "secondary"
	default:
		return nil, errors.New("invalid priority")
	}

	creditCardOptions := &easypostCreditCardCreateOptions{
		StripeToken: stripeToken,
		Priority:    priorityString,
	}

	creditCardRequest := &creditCardCreateRequest{
		CreditCard: creditCardOptions,
	}

	err = client.post(ctx, "credit_cards", creditCardRequest, &out)
	return
}
