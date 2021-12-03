package easypost

import (
	"context"
	"net/http"
	"time"
)

// An Insurance object represents insurance for packages purchased both via the
// EasyPost API as well as shipments purchased through third parties and later
// registered with EasyPost.
type Insurance struct {
	ID           string     `json:"id,omitempty"`
	Object       string     `json:"object,omitempty"`
	Reference    string     `json:"reference,omitempty"`
	Mode         string     `json:"mode,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	Amount       string     `json:"amount,omitempty"`
	Carrier      string     `json:"carrier,omitempty"`
	Provider     string     `json:"provider,omitempty"`
	ProviderID   string     `json:"provider_id,omitempty"`
	ShipmentID   string     `json:"shipment_id,omitempty"`
	TrackingCode string     `json:"tracking_code,omitempty"`
	Status       string     `json:"status,omitempty"`
	Tracker      *Tracker   `json:"tracker,omitempty"`
	ToAddress    *Address   `json:"to_address,omitempty"`
	FromAddress  *Address   `json:"from_address,omitempty"`
	Fee          *Fee       `json:"fee,omitempty"`
	Messages     []string   `json:"messages,omitempty"`
}

type createInsuranceRequest struct {
	Insurance *Insurance `json:"insurance,omitempty"`
}

// CreateInsurance creats an insurance object for a shipment purchased outside
// of EasyPost. ToAddress, FromAddress, TrackingCode and Amount fields must be
// provided. Providing a value in the Carrier field is optional, but can help
// avoid ambiguity and provide a shorter response time.
//	c := easypost.New(MyEasyPostAPIKey)
//	out, err := c.CreateInsurace(
//		&easypost.Insurance{
//			ToAddress:    &easypost.Address{ID: "adr_102"},
//			FromAddress:  &easypost.Address{ID: "adr_101"},
//			TrackingCode: "9400110898825022579493",
//			Carrier:      "USPS",
//			Reference:    "insuranceRef1",
//			Amount:       100,
//	)
func (c *Client) CreateInsurance(in *Insurance) (out *Insurance, err error) {
	req := &createInsuranceRequest{Insurance: in}
	err = c.post(context.Background(), "insurances", req, &out)
	return
}

// CreateInsuranceWithContext performs the same operation as CreateInsurance,
// but allows specifying a context that can interrupt the request.
func (c *Client) CreateInsuranceWithContext(ctx context.Context, in *Insurance) (out *Insurance, err error) {
	req := &createInsuranceRequest{Insurance: in}
	err = c.post(ctx, "insurances", req, &out)
	return
}

// ListInsurancesResult holds the results from the list insurances API.
type ListInsurancesResult struct {
	Insurances []*Insurance `json:"insurances,omitempty"`
	// HasMore indicates if there are more responses to be fetched. If True,
	// additional responses can be fetched by updating the ListInsurancesOptions
	// parameter's AfterID field with the ID of the last item in this object's
	// Insurances field.
	HasMore bool `json:"has_more,omitempty"`
}

// ListInsurances provides a paginated result of Insurance objects.
func (c *Client) ListInsurances(opts *ListOptions) (out *ListInsurancesResult, err error) {
	return c.ListInsurancesWithContext(context.Background(), opts)
}

// ListInsurancesWithContext performs the same operation as ListInsurances, but
// allows specifying a context that can interrupt the request.
func (c *Client) ListInsurancesWithContext(ctx context.Context, opts *ListOptions) (out *ListInsurancesResult, err error) {
	err = c.do(ctx, http.MethodGet, "insurances", c.convertOptsToURLValues(opts), &out)
	return
}

// GetInsurance returns the Insurance object with the given ID or reference.
func (c *Client) GetInsurance(insuranceID string) (out *Insurance, err error) {
	err = c.get(context.Background(), "insurances/"+insuranceID, &out)
	return
}

// GetInsuranceWithContext performs the same operation as GetInsurance, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetInsuranceWithContext(ctx context.Context, insuranceID string) (out *Insurance, err error) {
	err = c.get(ctx, "insurances/"+insuranceID, &out)
	return
}
