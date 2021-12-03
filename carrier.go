package easypost

import (
	"context"
	"time"
)

// CarrierField provides data for a single field in a carrier account.
type CarrierField struct {
	Visibility string `json:"visibility,omitempty"`
	Label      string `json:"label,omitempty"`
	Value      string `json:"value,omitempty"`
}

// CarrierFields contains the data for carrier account fields for production
// and/or test credentials.
type CarrierFields struct {
	Credentials     map[string]*CarrierField `json:"credentials,omitempty"`
	TestCredentials map[string]*CarrierField `json:"test_credentials,omitempty"`
	AutoLink        bool                     `json:"auto_link,omitempty"`
	CustomWorkflow  bool                     `json:"custom_workflow,omitempty"`
}

// CarrierAccount encapsulates credentials and other information related to a
// carrier account.
type CarrierAccount struct {
	ID              string            `json:"id,omitempty"`
	Object          string            `json:"object,omitempty"`
	Reference       string            `json:"reference,omitempty"`
	CreatedAt       *time.Time        `json:"created_at,omitempty"`
	UpdatedAt       *time.Time        `json:"updated_at,omitempty"`
	Type            string            `json:"type,omitempty"`
	Fields          *CarrierFields    `json:"fields,omitempty"`
	Clone           bool              `json:"clone,omitempty"`
	Description     string            `json:"description,omitempty"`
	Readable        string            `json:"readable,omitempty"`
	Credentials     map[string]string `json:"credentials"`
	TestCredentials map[string]string `json:"test_credentials"`
}

// CarrierType contains information on a supported carrier. It can be used to
// determine the valid fields for a carrier account.
type CarrierType struct {
	Object string         `json:"object,omitempty"`
	Type   string         `json:"type,omitempty"`
	Fields *CarrierFields `json:"fields,omitempty"`
}

// GetCarrierTypes returns a list of supported carrier types for the current
// user.
func (c *Client) GetCarrierTypes() (out []*CarrierType, err error) {
	err = c.get(context.Background(), "carrier_types", &out)
	return
}

// GetCarrierTypesWithContext performs the same operation as GetCarrierTypes,
// but allows specifying a context that can interrupt the request.
func (c *Client) GetCarrierTypesWithContext(ctx context.Context) (out []*CarrierType, err error) {
	err = c.get(ctx, "carrier_types", &out)
	return
}

type carrierAccountRequest struct {
	CarrierAccount *CarrierAccount `json:"carrier_account,omitempty"`
}

// CreateCarrierAccount creates a new carrier account. It can only be used with
// a production API key.
//	c := easypost.New(MyEasyPostAPIKey)
//	out, err := c.CreateCarrierAccount(
//		&easypost.CarrierAccount{
//			Type:        "UpsAccount",
//			Description: "NY Location UPS Account",
//			Reference:   "my reference",
//			Credentials: map[string]string{
//				"user_id":               "USERID",
//				"password":              "PASSWORD",
//				"access_license_number": "ALN",
//			},
//		},
//	)
func (c *Client) CreateCarrierAccount(in *CarrierAccount) (out *CarrierAccount, err error) {
	req := &carrierAccountRequest{CarrierAccount: in}
	err = c.post(context.Background(), "carrier_accounts", req, &out)
	return
}

// CreateCarrierAccountWithContext performs the same operation as
// CreateCarrierAccount, but allows specifying a context that can interrupt the
// request.
func (c *Client) CreateCarrierAccountWithContext(ctx context.Context, in *CarrierAccount) (out *CarrierAccount, err error) {
	req := &carrierAccountRequest{CarrierAccount: in}
	err = c.post(ctx, "carrier_accounts", req, &out)
	return
}

// ListCarrierAccounts returns a list of all carrier accounts available to the
// authenticated account.
func (c *Client) ListCarrierAccounts() (out []*CarrierAccount, err error) {
	err = c.get(context.Background(), "carrier_accounts", &out)
	return
}

// ListCarrierAccountsWithContext performs the same operation as
// ListCarrierAccounts, but allows specifying a context that can interrupt the
// request.
func (c *Client) ListCarrierAccountsWithContext(ctx context.Context) (out []*CarrierAccount, err error) {
	err = c.get(ctx, "carrier_accounts", &out)
	return
}

// GetCarrierAccount retrieves a carrier account by its ID or reference.
func (c *Client) GetCarrierAccount(carrierAccountID string) (out *CarrierAccount, err error) {
	err = c.get(context.Background(), "carrier_accounts/"+carrierAccountID, &out)
	return
}

// GetCarrierAccountWithContext performs the same operation as
// GetCarrierAccount, but allows specifying a context that can interrupt the
// request.
func (c *Client) GetCarrierAccountWithContext(ctx context.Context, carrierAccountID string) (out *CarrierAccount, err error) {
	err = c.get(ctx, "carrier_accounts/"+carrierAccountID, &out)
	return
}

// UpdateCarrierAccount updates the carrier account. Only the Description,
// Reference, Credentials and TestCredentials attributes can be updated.
//	c := easypost.New(MyEasyPostAPIKey)
//	out, err := c.UpdateCarrierAccount(
//		"ca_1001",
//		&easypost.CarrierAccount{
//			Description: "FL Location UPS Account",
//			Credentials: map[string]string{
//				"account_number": "B2B2B2",
//			},
//		},
//	)
func (c *Client) UpdateCarrierAccount(in *CarrierAccount) (out *CarrierAccount, err error) {
	req := &carrierAccountRequest{CarrierAccount: in}
	err = c.put(context.Background(), "carrier_accounts/"+in.ID, req, &out)
	return
}

// UpdateCarrierAccountWithContext performs the same operation as
// UpdateCarrierAccount, but allows specifying a context that can interrupt the
// request.
func (c *Client) UpdateCarrierAccountWithContext(ctx context.Context, in *CarrierAccount) (out *CarrierAccount, err error) {
	req := &carrierAccountRequest{CarrierAccount: in}
	err = c.put(ctx, "carrier_accounts/"+in.ID, req, &out)
	return
}

// DeleteCarrierAccount removes the carrier account with the given ID.
func (c *Client) DeleteCarrierAccount(carrierAccountID string) error {
	return c.del(context.Background(), "carrier_accounts/"+carrierAccountID)
}

// DeleteCarrierAccountWithContext performs the same operation as
// DeleteCarrierAccount, but allows specifying a context that can interrupt the
// request.
func (c *Client) DeleteCarrierAccountWithContext(ctx context.Context, carrierAccountID string) error {
	return c.del(ctx, "carrier_accounts/"+carrierAccountID)
}
