package easypost

import (
	"context"
	"net/http"
)

// An Insurance object represents insurance for packages purchased both via the
// EasyPost API and shipments purchased through third parties and later
// registered with EasyPost.
type Insurance struct {
	ID           string    `json:"id,omitempty"`
	Object       string    `json:"object,omitempty"`
	Reference    string    `json:"reference,omitempty"`
	Mode         string    `json:"mode,omitempty"`
	CreatedAt    *DateTime `json:"created_at,omitempty"`
	UpdatedAt    *DateTime `json:"updated_at,omitempty"`
	Amount       string    `json:"amount,omitempty"`
	Carrier      string    `json:"carrier,omitempty"`
	Provider     string    `json:"provider,omitempty"`
	ProviderID   string    `json:"provider_id,omitempty"`
	ShipmentID   string    `json:"shipment_id,omitempty"`
	TrackingCode string    `json:"tracking_code,omitempty"`
	Status       string    `json:"status,omitempty"`
	Tracker      *Tracker  `json:"tracker,omitempty"`
	ToAddress    *Address  `json:"to_address,omitempty"`
	FromAddress  *Address  `json:"from_address,omitempty"`
	Fee          *Fee      `json:"fee,omitempty"`
	Messages     []string  `json:"messages,omitempty"`
}

type createInsuranceRequest struct {
	Insurance *Insurance `json:"insurance,omitempty"`
}

// ListInsurancesResult holds the results from the list insurances API.
type ListInsurancesResult struct {
	Insurances []*Insurance `json:"insurances,omitempty"`
	PaginatedCollection
}

// CreateInsurance creates an insurance object for a shipment purchased outside
// EasyPost. ToAddress, FromAddress, TrackingCode and Amount fields must be
// provided. Providing a value in the Carrier field is optional, but can help
// avoid ambiguity and provide a shorter response time.
//
//	c := easypost.New(MyEasyPostAPIKey)
//	out, err := c.CreateInsurance(
//		&easypost.Insurance{
//			ToAddress:    &easypost.Address{ID: "adr_102"},
//			FromAddress:  &easypost.Address{ID: "adr_101"},
//			TrackingCode: "9400110898825022579493",
//			Carrier:      "USPS",
//			Reference:    "insuranceRef1",
//			Amount:       100,
//	)
func (c *Client) CreateInsurance(in *Insurance) (out *Insurance, err error) {
	return c.CreateInsuranceWithContext(context.Background(), in)
}

// CreateInsuranceWithContext performs the same operation as CreateInsurance,
// but allows specifying a context that can interrupt the request.
func (c *Client) CreateInsuranceWithContext(ctx context.Context, in *Insurance) (out *Insurance, err error) {
	req := &createInsuranceRequest{Insurance: in}
	err = c.post(ctx, "insurances", req, &out)
	return
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

// GetNextInsurancePage returns the next page of insurance records
func (c *Client) GetNextInsurancePage(collection *ListInsurancesResult) (out *ListInsurancesResult, err error) {
	return c.GetNextInsurancePageWithContext(context.Background(), collection)
}

// GetNextInsurancePageWithPageSize returns the next page of insurance records with a specific page size
func (c *Client) GetNextInsurancePageWithPageSize(collection *ListInsurancesResult, pageSize int) (out *ListInsurancesResult, err error) {
	return c.GetNextInsurancePageWithPageSizeWithContext(context.Background(), collection, pageSize)
}

// GetNextInsurancePageWithContext performs the same operation as GetNextInsurancePage, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetNextInsurancePageWithContext(ctx context.Context, collection *ListInsurancesResult) (out *ListInsurancesResult, err error) {
	return c.GetNextInsurancePageWithPageSizeWithContext(ctx, collection, 0)
}

// GetNextInsurancePageWithPageSizeWithContext performs the same operation as GetNextInsurancePageWithPageSize, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetNextInsurancePageWithPageSizeWithContext(ctx context.Context, collection *ListInsurancesResult, pageSize int) (out *ListInsurancesResult, err error) {
	if len(collection.Insurances) == 0 {
		err = EndOfPaginationError
		return
	}
	lastID := collection.Insurances[len(collection.Insurances)-1].ID
	params, err := nextPageParameters(collection.HasMore, lastID, pageSize)
	if err != nil {
		return
	}
	return c.ListInsurancesWithContext(ctx, params)
}

// GetInsurance returns the Insurance object with the given ID or reference.
func (c *Client) GetInsurance(insuranceID string) (out *Insurance, err error) {
	return c.GetInsuranceWithContext(context.Background(), insuranceID)
}

// GetInsuranceWithContext performs the same operation as GetInsurance, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetInsuranceWithContext(ctx context.Context, insuranceID string) (out *Insurance, err error) {
	err = c.get(ctx, "insurances/"+insuranceID, &out)
	return
}

// RefundInsurance refunds the Insurance object with the given ID.
func (c *Client) RefundInsurance(insuranceID string) (out *Insurance, err error) {
	return c.RefundInsuranceWithContext(context.Background(), insuranceID)
}

// RefundInsuranceWithContext performs the same operation as RefundInsurance, but
// allows specifying a context that can interrupt the request.
func (c *Client) RefundInsuranceWithContext(ctx context.Context, insuranceID string) (out *Insurance, err error) {
	err = c.post(ctx, "insurances/"+insuranceID+"/refund", nil, &out)
	return
}
