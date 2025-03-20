package easypost

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// A BetaPaymentRefund that has the refund details for the refund request.
type BetaPaymentRefund struct {
	RefundedAmount      int        `json:"refunded_amount,omitempty" url:"refunded_amount,omitempty"`
	RefundedPaymentLogs []string   `json:"refunded_payment_log,omitempty" url:"refunded_payment_log,omitempty"`
	PaymentLogId        string     `json:"payment_log_id,omitempty" url:"payment_log_id,omitempty"`
	Errors              []APIError `json:"errors,omitempty" url:"errors,omitempty"`
}

// A ReferralCustomer contains data about an EasyPost referral customer.
type ReferralCustomer struct {
	User
}

// ListReferralCustomersResult holds the results from the list referral customers API call.
type ListReferralCustomersResult struct {
	ReferralCustomers []*ReferralCustomer `json:"referral_customers,omitempty" url:"referral_customers,omitempty"`
	PaginatedCollection
}

// CreditCardOptions specifies options for creating or updating a credit card.
type CreditCardOptions struct {
	Number   string `json:"number,omitempty" url:"number,omitempty"`
	ExpMonth string `json:"expiration_month,omitempty" url:"expiration_month,omitempty"`
	ExpYear  string `json:"expiration_year,omitempty" url:"expiration_year,omitempty"`
	Cvc      string `json:"cvc,omitempty" url:"cvc,omitempty"`
}

// MandateData specifies the mandate data needed for adding a bank account from stripe.
type MandateData struct {
	IpAddress  string `json:"ip_address,omitempty" url:"ip_address,omitempty"`
	UserAgent  string `json:"user_agent,omitempty" url:"user_agent,omitempty"`
	AcceptedAt int64  `json:"accepted_at,omitempty" url:"accepted_at,omitempty"`
}

type stripeApiKeyResponse struct {
	PublicKey string `json:"public_key,omitempty" url:"public_key,omitempty"`
}

type stripeTokenResponse struct {
	Id string `json:"id,omitempty" url:"id,omitempty"`
}

type referralCustomerRequest struct {
	UserOptions *UserOptions `json:"user,omitempty" url:"user,omitempty"`
}

type creditCardCreateRequest struct {
	CreditCard *easypostCreditCardCreateOptions `json:"credit_card,omitempty" url:"credit_card,omitempty"`
}

type easypostCreditCardCreateOptions struct {
	StripeToken string `json:"stripe_object_id,omitempty" url:"stripe_object_id,omitempty"`
	Priority    string `json:"priority,omitempty" url:"priority,omitempty"`
}

type clientSecretResponse struct {
	ClientSecret string `json:"client_secret,omitempty" url:"client_secret,omitempty"`
}

// CreateReferralCustomer creates a new referral customer.
func (c *Client) CreateReferralCustomer(in *UserOptions) (out *ReferralCustomer, err error) {
	return c.CreateReferralCustomerWithContext(context.Background(), in)
}

// CreateReferralCustomerWithContext performs the same operation as CreateReferralCustomer, but allows
// specifying a context that can interrupt the request.
func (c *Client) CreateReferralCustomerWithContext(ctx context.Context, in *UserOptions) (out *ReferralCustomer, err error) {
	err = c.do(ctx, http.MethodPost, "referral_customers", &referralCustomerRequest{UserOptions: in}, &out)
	return
}

// ListReferralCustomers provides a paginated result of ReferralCustomer objects.
func (c *Client) ListReferralCustomers(opts *ListOptions) (out *ListReferralCustomersResult, err error) {
	return c.ListReferralCustomersWithContext(context.Background(), opts)
}

// ListReferralCustomersWithContext performs the same operation as ListReferralCustomers, but
// allows specifying a context that can interrupt the request.
func (c *Client) ListReferralCustomersWithContext(ctx context.Context, opts *ListOptions) (out *ListReferralCustomersResult, err error) {
	err = c.do(ctx, http.MethodGet, "referral_customers", opts, &out)
	return
}

// GetNextReferralCustomerPage returns the next page of referral customers
func (c *Client) GetNextReferralCustomerPage(collection *ListReferralCustomersResult) (out *ListReferralCustomersResult, err error) {
	return c.GetNextReferralCustomerPageWithContext(context.Background(), collection)
}

// GetNextReferralCustomerPageWithPageSize returns the next page of referral customers with a specific page size
func (c *Client) GetNextReferralCustomerPageWithPageSize(collection *ListReferralCustomersResult, pageSize int) (out *ListReferralCustomersResult, err error) {
	return c.GetNextReferralCustomerPageWithPageSizeWithContext(context.Background(), collection, pageSize)
}

// GetNextReferralCustomerPageWithContext performs the same operation as GetNextReferralCustomerPage, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetNextReferralCustomerPageWithContext(ctx context.Context, collection *ListReferralCustomersResult) (out *ListReferralCustomersResult, err error) {
	return c.GetNextReferralCustomerPageWithPageSizeWithContext(ctx, collection, 0)
}

// GetNextReferralCustomerPageWithPageSizeWithContext performs the same operation as GetNextReferralCustomerPageWithPageSize, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetNextReferralCustomerPageWithPageSizeWithContext(ctx context.Context, collection *ListReferralCustomersResult, pageSize int) (out *ListReferralCustomersResult, err error) {
	if len(collection.ReferralCustomers) == 0 {
		err = newEndOfPaginationError()
		return
	}
	lastID := collection.ReferralCustomers[len(collection.ReferralCustomers)-1].ID
	params, err := nextPageParameters(collection.HasMore, lastID, pageSize)
	if err != nil {
		return
	}
	return c.ListReferralCustomersWithContext(ctx, params)
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
	err = c.do(ctx, http.MethodPut, "referral_customers/"+userId, req, &out)
	return
}

// AddReferralCustomerCreditCard adds a credit card to EasyPost for a ReferralCustomer without needing a Stripe account.
func (c *Client) AddReferralCustomerCreditCard(referralCustomerApiKey string, creditCardOptions *CreditCardOptions, priority PaymentMethodPriority) (out *PaymentMethodObject, err error) {
	return c.AddReferralCustomerCreditCardWithContext(context.Background(), referralCustomerApiKey, creditCardOptions, priority)
}

// AddReferralCustomerCreditCardWithContext performs the same operation as AddReferralCustomerCreditCard, but allows
// specifying a context that can interrupt the request.
func (c *Client) AddReferralCustomerCreditCardWithContext(ctx context.Context, referralCustomerApiKey string, creditCardOptions *CreditCardOptions, priority PaymentMethodPriority) (out *PaymentMethodObject, err error) {
	stripeApiKeyResponse, err := c.retrieveEasypostStripeApiKey(ctx)
	if err != nil || stripeApiKeyResponse == nil || stripeApiKeyResponse.PublicKey == "" {
		return nil, &InternalServerError{
			APIError: APIError{
				Code:       "Could not create Stripe token, please try again later",
				StatusCode: 500,
			},
		}
	}

	stripeTokenResponse, err := c.createStripeToken(ctx, stripeApiKeyResponse.PublicKey, creditCardOptions)
	if err != nil || stripeTokenResponse == nil || stripeTokenResponse.Id == "" {
		return nil, newExternalApiError("Could not create Stripe token, please try again later")
	}

	return c.createEasypostCreditCard(ctx, referralCustomerApiKey, stripeTokenResponse.Id, priority)
}

// AddReferralCustomerCreditCardFromStripe adds a credit card to EasyPost for a ReferralCustomer with a payment method ID from Stripe. This function requires the ReferralCustomer User's API key.
func (c *Client) AddReferralCustomerCreditCardFromStripe(referralCustomerApiKey string, paymentMethodId string, priority PaymentMethodPriority) (out *PaymentMethodObject, err error) {
	return c.AddReferralCustomerCreditCardFromStripeWithContext(context.Background(), referralCustomerApiKey, paymentMethodId, priority)
}

// AddReferralCustomerCreditCardFromStripeWithContext performs the same operation as AddReferralCustomerCreditCardFromStripe, but allows
// specifying a context that can interrupt the request.
func (c *Client) AddReferralCustomerCreditCardFromStripeWithContext(ctx context.Context, referralCustomerApiKey string, paymentMethodId string, priority PaymentMethodPriority) (out *PaymentMethodObject, err error) {
	client := &Client{
		APIKey: referralCustomerApiKey,
		Client: c.client(), // pass the current client's inner http.Client (configured to record) to the new client
	}

	params := map[string]interface{}{
		"credit_card": map[string]interface{}{
			"payment_method_id": paymentMethodId,
			"priority":          c.GetPaymentEndpointByPrimaryOrSecondary(priority),
		},
	}

	err = client.do(ctx, http.MethodPost, "credit_cards", params, &out)
	return
}

// AddReferralCustomerBankAccountFromStripe adds a credit card to EasyPost for a ReferralCustomer with a payment method ID from Stripe. This function requires the ReferralCustomer User's API key.
func (c *Client) AddReferralCustomerBankAccountFromStripe(referralCustomerApiKey string, paymentMethodId string, mandateData *MandateData, priority PaymentMethodPriority) (out *PaymentMethodObject, err error) {
	return c.AddReferralCustomerBankAccountFromStripeWithContext(context.Background(), referralCustomerApiKey, paymentMethodId, mandateData, priority)
}

// AddReferralCustomerBankAccountFromStripeWithContext performs the same operation as AddReferralCustomerBankAccountFromStripe, but allows
// specifying a context that can interrupt the request.
func (c *Client) AddReferralCustomerBankAccountFromStripeWithContext(ctx context.Context, referralCustomerApiKey string, financialConnectionsId string, mandateData *MandateData, priority PaymentMethodPriority) (out *PaymentMethodObject, err error) {
	client := &Client{
		APIKey: referralCustomerApiKey,
		Client: c.client(), // pass the current client's inner http.Client (configured to record) to the new client
	}

	params := map[string]interface{}{
		"financial_connections_id": financialConnectionsId,
		"mandata_data":             mandateData,
		"priority":                 c.GetPaymentEndpointByPrimaryOrSecondary(priority),
	}

	err = client.do(ctx, http.MethodPost, "bank_accounts", params, &out)
	return
}

func (c *Client) retrieveEasypostStripeApiKey(ctx context.Context) (out *stripeApiKeyResponse, err error) {
	err = c.do(ctx, http.MethodGet, "partners/stripe_public_key", nil, &out)
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

	body, err := io.ReadAll(resp.Body) // deprecated, but we have to keep it for legacy compatibility
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

	creditCardOptions := &easypostCreditCardCreateOptions{
		StripeToken: stripeToken,
		Priority:    c.GetPaymentEndpointByPrimaryOrSecondary(priority),
	}

	creditCardRequest := &creditCardCreateRequest{
		CreditCard: creditCardOptions,
	}

	err = client.do(ctx, http.MethodPost, "credit_cards", creditCardRequest, &out)
	return
}

// BetaAddPaymentMethod adds Stripe payment method to referral customer.
func (c *Client) BetaAddPaymentMethod(stripeCustomerId string, paymentMethodReference string, priority PaymentMethodPriority) (out *PaymentMethodObject, err error) {
	return c.BetaAddPaymentMethodWithContext(context.Background(), stripeCustomerId, paymentMethodReference, priority)
}

// BetaAddPaymentMethodWithContext performs the same operation as BetaAddPaymentMethod, but allows
// specifying a context that can interrupt the request.
func (c *Client) BetaAddPaymentMethodWithContext(ctx context.Context, stripeCustomerId string, paymentMethodReference string, priority PaymentMethodPriority) (out *PaymentMethodObject, err error) {
	wrappedParams := map[string]interface{}{
		"payment_method": map[string]interface{}{
			"stripe_customer_id":       stripeCustomerId,
			"payment_method_reference": paymentMethodReference,
			"priority":                 c.GetPaymentEndpointByPrimaryOrSecondary(priority),
		},
	}

	err = c.do(ctx, http.MethodPost, "/beta/referral_customers/payment_method", wrappedParams, out)
	return
}

// BetaRefundByAmount refunds a recent payment by amount in cents.
func (c *Client) BetaRefundByAmount(refundAmount int) (out *BetaPaymentRefund, err error) {
	return c.BetaRefundByAmountWithContext(context.Background(), refundAmount)
}

// BetaRefundByAmountWithContext performs the same operation as BetaRefundByAmount, but allows
// specifying a context that can interrupt the request.
func (c *Client) BetaRefundByAmountWithContext(ctx context.Context, refundAmount int) (out *BetaPaymentRefund, err error) {
	params := map[string]interface{}{
		"refund_amount": refundAmount,
	}

	err = c.do(ctx, http.MethodPost, "/beta/referral_customers/refunds", params, out)
	return
}

// BetaRefundByPaymentLog refunds a payment by payment log ID.
func (c *Client) BetaRefundByPaymentLog(paymentLogId string) (out *BetaPaymentRefund, err error) {
	return c.BetaRefundByPaymentLogWithContext(context.Background(), paymentLogId)
}

// BetaRefundByPaymentLogWithContext performs the same operation as BetaRefundByPaymentLog, but allows
// specifying a context that can interrupt the request.
func (c *Client) BetaRefundByPaymentLogWithContext(ctx context.Context, paymentLogId string) (out *BetaPaymentRefund, err error) {
	params := map[string]interface{}{
		"payment_log_id": paymentLogId,
	}

	err = c.do(ctx, http.MethodPost, "/beta/referral_customers/refunds", params, out)
	return
}

// BetaCreateCreditCardClientSecret creates a client secret to use with Stripe when adding a credit card.
func (c *Client) BetaCreateCreditCardClientSecret() (out *clientSecretResponse, err error) {
	return c.BetaCreateCreditCardClientSecretWithContext(context.Background())
}

// BetaCreateCreditCardClientSecretWithContext performs the same operation as BetaCreateCreditCardClientSecret, but allows
// specifying a context that can interrupt the request.
func (c *Client) BetaCreateCreditCardClientSecretWithContext(ctx context.Context) (out *clientSecretResponse, err error) {
	out = &clientSecretResponse{}
	err = c.do(ctx, http.MethodPost, "/beta/setup_intents", nil, out)
	return
}

// BetaCreateBankAccountClientSecret creates a client secret to use with Stripe when adding a bank account.
func (c *Client) BetaCreateBankAccountClientSecret() (out *clientSecretResponse, err error) {
	return c.BetaCreateBankAccountClientSecretWithContext(context.Background())
}

// BetaCreateBankAccountClientSecretWithContext performs the same operation as BetaCreateBankAccountClientSecret, but allows
// specifying a context that can interrupt the request.
func (c *Client) BetaCreateBankAccountClientSecretWithContext(ctx context.Context) (out *clientSecretResponse, err error) {
	out = &clientSecretResponse{}
	err = c.do(ctx, http.MethodPost, "/beta/financial_connections_sessions", nil, out)
	return
}

// BetaCreateBankAccountClientSecretWithReturlUrl creates a client secret to use with Stripe when adding a bank account.
func (c *Client) BetaCreateBankAccountClientSecretWithReturlUrl(returnUrl string) (out *clientSecretResponse, err error) {
	return c.BetaCreateBankAccountClientSecretWithReturnUrlWithContext(context.Background(), returnUrl)
}

// BetaCreateBankAccountClientSecretWithReturnUrlWithContext performs the same operation as BetaCreateBankAccountClientSecretWithReturlUrl, but allows
// specifying a context that can interrupt the request.
func (c *Client) BetaCreateBankAccountClientSecretWithReturnUrlWithContext(ctx context.Context, returnUrl string) (out *clientSecretResponse, err error) {
	out = &clientSecretResponse{}
	params := map[string]interface{}{
		"return_url": returnUrl,
	}
	err = c.do(ctx, http.MethodPost, "/beta/financial_connections_sessions", params, out)
	return
}
