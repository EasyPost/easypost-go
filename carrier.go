package easypost

import (
	"context"
	"net/http"
)

// CarrierField provides data for a single field in a carrier account.
type CarrierField struct {
	Visibility string `json:"visibility,omitempty" url:"visibility,omitempty"`
	Label      string `json:"label,omitempty" url:"label,omitempty"`
	Value      string `json:"value,omitempty" url:"value,omitempty"`
}

// CarrierFields contains the data for carrier account fields for production and/or test credentials.
type CarrierFields struct {
	Credentials     map[string]*CarrierField `json:"credentials,omitempty" url:"credentials,omitempty"`
	TestCredentials map[string]*CarrierField `json:"test_credentials,omitempty" url:"test_credentials,omitempty"`
	AutoLink        bool                     `json:"auto_link,omitempty" url:"auto_link,omitempty"`
	CustomWorkflow  bool                     `json:"custom_workflow,omitempty" url:"custom_workflow,omitempty"`
}

// CarrierAccount encapsulates credentials and other information related to a carrier account.
// This struct is also used as a parameter set for creating and updating non-UPS carrier accounts.
type CarrierAccount struct {
	ID               string                 `json:"id,omitempty" url:"id,omitempty"`
	Object           string                 `json:"object,omitempty" url:"object,omitempty"`
	Reference        string                 `json:"reference,omitempty" url:"reference,omitempty"`
	CreatedAt        *DateTime              `json:"created_at,omitempty" url:"created_at,omitempty"`
	UpdatedAt        *DateTime              `json:"updated_at,omitempty" url:"updated_at,omitempty"`
	Type             string                 `json:"type,omitempty" url:"type,omitempty"`
	Fields           *CarrierFields         `json:"fields,omitempty" url:"fields,omitempty"`
	Clone            bool                   `json:"clone,omitempty" url:"clone,omitempty"`
	Description      string                 `json:"description,omitempty" url:"description,omitempty"`
	Readable         string                 `json:"readable,omitempty" url:"readable,omitempty"`
	Credentials      map[string]string      `json:"credentials,omitempty" url:"credentials,omitempty"`
	TestCredentials  map[string]string      `json:"test_credentials,omitempty" url:"test_credentials,omitempty"`
	RegistrationData map[string]interface{} `json:"registration_data,omitempty" url:"registration_data,omitempty"`
	BillingType      string                 `json:"billing_type,omitempty" url:"billing_type,omitempty"`
}

// CarrierType contains information on a supported carrier. It can be used to determine the valid fields for a carrier account.
type CarrierType struct {
	Object   string         `json:"object,omitempty" url:"object,omitempty"`
	Type     string         `json:"type,omitempty" url:"type,omitempty"`
	Readable string         `json:"readable,omitempty" url:"readable,omitempty"`
	Logo     string         `json:"logo,omitempty" url:"logo,omitempty"`
	Fields   *CarrierFields `json:"fields,omitempty" url:"fields,omitempty"`
}

type carrierAccountRequest struct {
	Data *CarrierAccount `json:"carrier_account,omitempty" url:"carrier_account,omitempty"`
}

// UpsCarrierAccountCreationParameters contains the parameters needed to create a new UPS carrier account.
type UpsCarrierAccountCreationParameters struct {
	Type          string `json:"type,omitempty" url:"type,omitempty"`
	Description   string `json:"description,omitempty" url:"description,omitempty"`
	Reference     string `json:"reference,omitempty" url:"reference,omitempty"`
	AccountNumber string `json:"account_number,omitempty" url:"account_number,omitempty"`
}

type upsCarrierAccountCreationRequest struct {
	Data *UpsCarrierAccountCreationParameters `json:"carrier_account_oauth_registrations,omitempty" url:"carrier_account_oauth_registrations,omitempty"`
}

type oauthCarrierAccountCreationRequest struct {
	Data *CarrierAccount `json:"carrier_account_oauth_registrations,omitempty" url:"carrier_account_oauth_registrations,omitempty"`
}

// UpsCarrierAccountUpdateParameters contains the parameters needed to update a UPS carrier account.
type UpsCarrierAccountUpdateParameters struct {
	AccountNumber string `json:"account_number,omitempty" url:"account_number,omitempty"`
}

type upsCarrierAccountUpdateRequest struct {
	Data *UpsCarrierAccountUpdateParameters `json:"carrier_account,omitempty" url:"carrier_account,omitempty"`
}

func (c *Client) selectCarrierAccountCreationEndpoint(typ string) string {
	for _, carrier := range getCarrierAccountTypesWithCustomWorkflows() {
		if typ == carrier {
			return "carrier_accounts/register"
		}
	}

	for _, carrier := range getUpsCarrierAccountTypes() {
		if typ == carrier {
			return "carrier_accounts/register_oauth"
		}
	}

	for _, carrier := range getCarrierAccountTypesWithCustomOauth() {
		if typ == carrier {
			return "carrier_accounts/register_oauth"
		}
	}

	return "carrier_accounts"
}

// GetCarrierTypes returns a list of supported carrier types for the current user.
func (c *Client) GetCarrierTypes() (out []*CarrierType, err error) {
	return c.GetCarrierTypesWithContext(context.Background())
}

// GetCarrierTypesWithContext performs the same operation as GetCarrierTypes, but allows specifying a context that can interrupt the request.
func (c *Client) GetCarrierTypesWithContext(ctx context.Context) (out []*CarrierType, err error) {
	err = c.do(ctx, http.MethodGet, "carrier_types", nil, &out)
	return
}

// CreateCarrierAccount creates a new carrier account. It can only be used with a production API key.
//
//	c := easypost.New(MyEasyPostAPIKey)
//	out, err := c.CreateCarrierAccount(
//		&easypost.CarrierAccount{
//			Type:        "DhlEcsAccount",
//			Description: "CA Location DHL eCommerce Solutions Account",
//			Credentials: map[string]string{
//				"client_id":           "123456",
//				"client_secret":       "123abc",
//				"distribution_center": "USLAX1",
//				"pickup_id":           "123456",
//			},
//			TestCredentials: map[string]string{
//				"client_id":           "123456",
//				"client_secret":       "123abc",
//				"distribution_center": "USLAX1",
//				"pickup_id":           "123456",
//			},
//		},
//	)
//
// Users cannot create UPS accounts with this function, must use CreateUpsCarrierAccount. An error will be returned if the user tries to create a UPS account with this function.
func (c *Client) CreateCarrierAccount(in *CarrierAccount) (out *CarrierAccount, err error) {
	return c.CreateCarrierAccountWithContext(context.Background(), in)
}

// CreateCarrierAccountWithContext performs the same operation as CreateCarrierAccount, but allows specifying a context that can interrupt the request.
// Users cannot create UPS accounts with this function, must use CreateUpsCarrierAccount. An error will be returned if the user tries to create a UPS account with this function.
func (c *Client) CreateCarrierAccountWithContext(ctx context.Context, in *CarrierAccount) (out *CarrierAccount, err error) {
	// Users cannot create UPS accounts with this function, must use CreateUpsCarrierAccount
	for _, carrier := range getUpsCarrierAccountTypes() {
		if in.Type == carrier {
			return nil, newInvalidFunctionError("users must use CreateUpsCarrierAccount to create UPS accounts")
		}
	}

	for _, carrier := range getCarrierAccountTypesWithCustomOauth() {
		if in.Type == carrier {
			return nil, newInvalidFunctionError("users must use CreateUpsCarrierAccount to create Oauth accounts")
		}
	}

	req := &carrierAccountRequest{Data: in}
	endpoint := c.selectCarrierAccountCreationEndpoint(in.Type)
	err = c.do(ctx, http.MethodPost, endpoint, req, &out)
	return
}

// CreateUpsCarrierAccount creates a new UPS carrier account. It can only be used with a production API key.
// Users cannot create non-UPS accounts with this function, must use CreateCarrierAccount. An error will be returned if the user tries to create a non-UPS account with this function.
func (c *Client) CreateUpsCarrierAccount(in *UpsCarrierAccountCreationParameters) (out *CarrierAccount, err error) {
	return c.CreateUpsCarrierAccountWithContext(context.Background(), in)
}

// CreateUpsCarrierAccountWithContext performs the same operation as CreateUpsCarrierAccount, but allows specifying a context that can interrupt the request.
// Users cannot create non-UPS accounts with this function, must use CreateCarrierAccount. An error will be returned if the user tries to create a non-UPS account with this function.
func (c *Client) CreateUpsCarrierAccountWithContext(ctx context.Context, in *UpsCarrierAccountCreationParameters) (out *CarrierAccount, err error) {
	// Users cannot create non-UPS accounts with this function, must use CreateCarrierAccount
	isUpsAccount := false
	for _, carrier := range getUpsCarrierAccountTypes() {
		if in.Type == carrier {
			isUpsAccount = true
			break
		}
	}
	if !isUpsAccount {
		return nil, newInvalidFunctionError("users must use CreateCarrierAccount to create non-UPS accounts")
	}

	req := &upsCarrierAccountCreationRequest{Data: in}
	endpoint := c.selectCarrierAccountCreationEndpoint(in.Type)
	err = c.do(ctx, http.MethodPost, endpoint, req, &out)
	return
}

// CreateOauthCarrierAccount creates a new carrier account that uses Oauth. It can only be used with a production API key.
// Users cannot create non-oauth accounts with this function, must use CreateCarrierAccount. An error will be returned if the user tries to create a non-oauth account with this function.
func (c *Client) CreateOauthCarrierAccount(in *CarrierAccount) (out *CarrierAccount, err error) {
	return c.CreateOauthCarrierAccountWithContext(context.Background(), in)
}

// CreateOauthCarrierAccountWithContext performs the same operation as CreateOauthCarrierAccount, but allows specifying a context that can interrupt the request.
// Users cannot create non-oauth accounts with this function, must use CreateCarrierAccount. An error will be returned if the user tries to create a non-oauth account with this function.
func (c *Client) CreateOauthCarrierAccountWithContext(ctx context.Context, in *CarrierAccount) (out *CarrierAccount, err error) {
	// Users cannot create non-oauth accounts with this function, must use CreateCarrierAccount
	isOauthAccount := false
	for _, carrier := range getCarrierAccountTypesWithCustomOauth() {
		if in.Type == carrier {
			isOauthAccount = true
			break
		}
	}
	if !isOauthAccount {
		return nil, newInvalidFunctionError("users must use CreateCarrierAccount to create non-oauth accounts")
	}

	req := &oauthCarrierAccountCreationRequest{Data: in}
	endpoint := c.selectCarrierAccountCreationEndpoint(in.Type)
	err = c.do(ctx, http.MethodPost, endpoint, req, &out)
	return
}

// ListCarrierAccounts returns a list of all carrier accounts available to the authenticated account.
func (c *Client) ListCarrierAccounts() (out []*CarrierAccount, err error) {
	return c.ListCarrierAccountsWithContext(context.Background())
}

// ListCarrierAccountsWithContext performs the same operation as ListCarrierAccounts, but allows specifying a context that can interrupt the request.
func (c *Client) ListCarrierAccountsWithContext(ctx context.Context) (out []*CarrierAccount, err error) {
	err = c.do(ctx, http.MethodGet, "carrier_accounts", nil, &out)
	return
}

// GetCarrierAccount retrieves a carrier account by its ID or reference.
func (c *Client) GetCarrierAccount(carrierAccountID string) (out *CarrierAccount, err error) {
	return c.GetCarrierAccountWithContext(context.Background(), carrierAccountID)
}

// GetCarrierAccountWithContext performs the same operation as GetCarrierAccount, but allows specifying a context that can interrupt the request.
func (c *Client) GetCarrierAccountWithContext(ctx context.Context, carrierAccountID string) (out *CarrierAccount, err error) {
	err = c.do(ctx, http.MethodGet, "carrier_accounts/"+carrierAccountID, nil, &out)
	return
}

// UpdateCarrierAccount updates the carrier account. Only the Description, Reference, Credentials and TestCredentials attributes can be updated.
//
//	c := easypost.New(MyEasyPostAPIKey)
//	out, err := c.UpdateCarrierAccount(
//		&easypost.CarrierAccount{
//			ID:          "ca_...",
//			Description: "FL Location DHL eCommerce Solutions Account",
//			Credentials: map[string]string{
//				"pickup_id": "abc123",
//			},
//		},
//	)
//
// Users cannot update UPS accounts with this function, must use UpdateUpsCarrierAccount. An error will be returned if the user tries to update a UPS account with this function.
func (c *Client) UpdateCarrierAccount(in *CarrierAccount) (out *CarrierAccount, err error) {
	return c.UpdateCarrierAccountWithContext(context.Background(), in)
}

// UpdateCarrierAccountWithContext performs the same operation as UpdateCarrierAccount, but allows specifying a context that can interrupt the request.
// Users cannot update UPS accounts with this function, must use UpdateUpsCarrierAccount. An error will be returned if the user tries to update a UPS account with this function.
func (c *Client) UpdateCarrierAccountWithContext(ctx context.Context, in *CarrierAccount) (out *CarrierAccount, err error) {
	account, err := c.GetCarrierAccount(in.ID)
	if err != nil {
		return nil, err
	}

	// Users cannot update UPS accounts with this function, must use UpdateUpsCarrierAccount
	for _, carrier := range getUpsCarrierAccountTypes() {
		if account.Type == carrier {
			return nil, newInvalidFunctionError("users must use UpdateUpsCarrierAccount to update UPS accounts")
		}
	}

	req := &carrierAccountRequest{Data: in}
	err = c.do(ctx, http.MethodPatch, "carrier_accounts/"+in.ID, req, &out)
	return
}

// UpdateUpsCarrierAccount updates a UPS carrier account.
// Users cannot update non-UPS accounts with this function, must use UpdateCarrierAccount. An error will be returned if the user tries to update a non-UPS account with this function.
func (c *Client) UpdateUpsCarrierAccount(id string, in *UpsCarrierAccountUpdateParameters) (out *CarrierAccount, err error) {
	return c.UpdateUpsCarrierAccountWithContext(context.Background(), id, in)
}

// UpdateUpsCarrierAccountWithContext performs the same operation as UpdateUpsCarrierAccount, but allows specifying a context that can interrupt the request.
// Users cannot update non-UPS accounts with this function, must use UpdateCarrierAccount. An error will be returned if the user tries to update a non-UPS account with this function.
func (c *Client) UpdateUpsCarrierAccountWithContext(ctx context.Context, id string, in *UpsCarrierAccountUpdateParameters) (out *CarrierAccount, err error) {
	account, err := c.GetCarrierAccount(id)
	if err != nil {
		return nil, err
	}

	// Users cannot update non-UPS accounts with this function, must use UpdateCarrierAccount
	isUpsAccount := false
	for _, carrier := range getUpsCarrierAccountTypes() {
		if account.Type == carrier {
			isUpsAccount = true
			break
		}
	}
	if !isUpsAccount {
		return nil, newInvalidFunctionError("users must use UpdateCarrierAccount to update non-UPS accounts")
	}

	req := &upsCarrierAccountUpdateRequest{Data: in}
	err = c.do(ctx, http.MethodPatch, "carrier_accounts/"+id, req, &out)
	return
}

// DeleteCarrierAccount removes the carrier account with the given ID.
func (c *Client) DeleteCarrierAccount(carrierAccountID string) error {
	return c.DeleteCarrierAccountWithContext(context.Background(), carrierAccountID)
}

// DeleteCarrierAccountWithContext performs the same operation as DeleteCarrierAccount, but allows specifying a context that can interrupt the request.
func (c *Client) DeleteCarrierAccountWithContext(ctx context.Context, carrierAccountID string) error {
	return c.do(ctx, http.MethodDelete, "carrier_accounts/"+carrierAccountID, nil, nil)
}
