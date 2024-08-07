package easypost

import (
	"context"
	"net/http"
)

// A ScanForm object represents a document that can be scanned to mark all
// included tracking codes as "Accepted for Shipment" by the carrier.
type ScanForm struct {
	ID            string    `json:"id,omitempty" url:"id,omitempty"`
	Object        string    `json:"object,omitempty" url:"object,omitempty"`
	CreatedAt     *DateTime `json:"created_at,omitempty" url:"created_at,omitempty"`
	UpdatedAt     *DateTime `json:"updated_at,omitempty" url:"updated_at,omitempty"`
	Status        string    `json:"status,omitempty" url:"status,omitempty"`
	Message       string    `json:"message,omitempty" url:"message,omitempty"`
	Address       *Address  `json:"address,omitempty" url:"address,omitempty"`
	TrackingCodes []string  `json:"tracking_codes,omitempty" url:"tracking_codes,omitempty"`
	FormURL       string    `json:"form_url,omitempty" url:"form_url,omitempty"`
	FormFileType  string    `json:"form_file_type,omitempty" url:"form_file_type,omitempty"`
	BatchID       string    `json:"batch_id,omitempty" url:"batch_id,omitempty"`
}

type scanFormRequest struct {
	Shipments []*Shipment `json:"shipments,omitempty" url:"shipments,omitempty"`
}

// ListScanFormsResult holds the results from the list scan forms API.
type ListScanFormsResult struct {
	ScanForms []*ScanForm `json:"scan_forms,omitempty" url:"scan_forms,omitempty"`
	PaginatedCollection
}

func newScanFormRequest(shipmentIDs ...string) *scanFormRequest {
	length := len(shipmentIDs)

	req := &scanFormRequest{Shipments: make([]*Shipment, length)}
	for i := 0; i < length; i++ {
		req.Shipments[i] = &Shipment{ID: shipmentIDs[i]}
	}
	return req
}

// CreateScanForm creates a scan form for the given Shipments.
//
//	c := easypost.New(MyEasyPostAPIKey)
//	out, err := c.CreateScanForm("shp_1", "shp_2")
func (c *Client) CreateScanForm(shipmentIDs ...string) (out *ScanForm, err error) {
	return c.CreateScanFormWithContext(context.Background(), shipmentIDs...)
}

// CreateScanFormWithContext performs the same operation as CreateScanForm, but
// allows specifying a context that can interrupt the request.
func (c *Client) CreateScanFormWithContext(ctx context.Context, shipmentIDs ...string) (out *ScanForm, err error) {
	err = c.do(ctx, http.MethodPost, "scan_forms", newScanFormRequest(shipmentIDs...), &out)
	return
}

// ListScanForms provides a paginated result of ScanForm objects.
func (c *Client) ListScanForms(opts *ListOptions) (out *ListScanFormsResult, err error) {
	return c.ListScanFormsWithContext(context.Background(), opts)
}

// ListScanFormsWithContext performs the same operation as ListScanForms, but
// allows specifying a context that can interrupt the request.
func (c *Client) ListScanFormsWithContext(ctx context.Context, opts *ListOptions) (out *ListScanFormsResult, err error) {
	err = c.do(ctx, http.MethodGet, "scan_forms", opts, &out)
	return
}

// GetNextScanFormPage returns the next page of scan forms
func (c *Client) GetNextScanFormPage(collection *ListScanFormsResult) (out *ListScanFormsResult, err error) {
	return c.GetNextScanFormPageWithContext(context.Background(), collection)
}

// GetNextScanFormPageWithPageSize returns the next page of scan forms with a specific page size
func (c *Client) GetNextScanFormPageWithPageSize(collection *ListScanFormsResult, pageSize int) (out *ListScanFormsResult, err error) {
	return c.GetNextScanFormPageWithPageSizeWithContext(context.Background(), collection, pageSize)
}

// GetNextScanFormPageWithContext performs the same operation as GetNextScanFormPage, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetNextScanFormPageWithContext(ctx context.Context, collection *ListScanFormsResult) (out *ListScanFormsResult, err error) {
	return c.GetNextScanFormPageWithPageSizeWithContext(ctx, collection, 0)
}

// GetNextScanFormPageWithPageSizeWithContext performs the same operation as GetNextScanFormPageWithPageSize, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetNextScanFormPageWithPageSizeWithContext(ctx context.Context, collection *ListScanFormsResult, pageSize int) (out *ListScanFormsResult, err error) {
	if len(collection.ScanForms) == 0 {
		err = EndOfPaginationError
		return
	}
	lastID := collection.ScanForms[len(collection.ScanForms)-1].ID
	params, err := nextPageParameters(collection.HasMore, lastID, pageSize)
	if err != nil {
		return
	}
	return c.ListScanFormsWithContext(ctx, params)
}

// GetScanForm retrieves a ScanForm object by ID.
func (c *Client) GetScanForm(scanFormID string) (out *ScanForm, err error) {
	return c.GetScanFormWithContext(context.Background(), scanFormID)
}

// GetScanFormWithContext performs the same operation as GetScanForm, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetScanFormWithContext(ctx context.Context, scanFormID string) (out *ScanForm, err error) {
	err = c.do(ctx, http.MethodGet, "scan_forms/"+scanFormID, nil, &out)
	return
}
