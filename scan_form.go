package easypost

import (
	"context"
	"net/http"
	"time"
)

// A ScanForm object represents a document that can be scanned to mark all
// included tracking codes as "Accepted for Shipment" by the carrier.
type ScanForm struct {
	ID            string     `json:"id,omitempty"`
	Object        string     `json:"object,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	Status        string     `json:"status,omitempty"`
	Message       string     `json:"message,omitempty"`
	Address       *Address   `json:"address,omitempty"`
	TrackingCodes []string   `json:"tracking_codes,omitempty"`
	FormURL       string     `json:"form_url,omitempty"`
	FormFileType  string     `json:"form_file_type,omitempty"`
	BatchID       string     `json:"batch_id,omitempty"`
}

type scanFormRequest struct {
	Shipments []*Shipment `json:"shipments,omitempty"`
}

func newscanFormRequest(shipmentIDs ...string) *scanFormRequest {
	l := len(shipmentIDs)
	if l == 0 {
		return nil
	}
	req := &scanFormRequest{Shipments: make([]*Shipment, l)}
	for i := 0; i < l; i++ {
		req.Shipments[i] = &Shipment{ID: shipmentIDs[i]}
	}
	return req
}

// CreateScanForm creates a scan form for the given Shipments.
//	c := easypost.New(MyEasyPostAPIKey)
//	out, err := c.CreateScanForm("shp_1", "shp_2")
func (c *Client) CreateScanForm(shipmentIDs ...string) (out *ScanForm, err error) {
	err = c.post(context.Background(), "scan_forms", newscanFormRequest(shipmentIDs...), &out)
	return
}

// CreateScanFormWithContext performs the same operation as CreateScanForm, but
// allows specifying a context that can interrupt the request.
func (c *Client) CreateScanFormWithContext(ctx context.Context, shipmentIDs ...string) (out *ScanForm, err error) {
	err = c.post(ctx, "scan_forms", newscanFormRequest(shipmentIDs...), &out)
	return
}

// ListScanFormsResult holds the results from the list scan forms API.
type ListScanFormsResult struct {
	ScanForms []*ScanForm `json:"scan_forms,omitempty"`
	// HasMore indicates if there are more responses to be fetched. If True,
	// additional responses can be fetched by updating the ListScanFormsOptions
	// parameter's AfterID field with the ID of the last item in this object's
	// ScanForms field.
	HasMore bool `json:"has_more,omitempty"`
}

// ListScanForms provides a paginated result of ScanForm objects.
func (c *Client) ListScanForms(opts *ListOptions) (out *ListScanFormsResult, err error) {
	return c.ListScanFormsWithContext(context.Background(), opts)
}

// ListScanFormsWithContext performs the same operation as ListScanForms, but
// allows specifying a context that can interrupt the request.
func (c *Client) ListScanFormsWithContext(ctx context.Context, opts *ListOptions) (out *ListScanFormsResult, err error) {
	err = c.do(ctx, http.MethodGet, "scan_forms", c.convertOptsToURLValues(opts), &out)
	return
}

// GetScanForm retrieves a ScanForm object by ID.
func (c *Client) GetScanForm(scanFormID string) (out *ScanForm, err error) {
	err = c.get(context.Background(), "scan_forms/"+scanFormID, &out)
	return
}

// GetScanFormWithContext performs the same operation as GetScanForm, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetScanFormWithContext(ctx context.Context, scanFormID string) (out *ScanForm, err error) {
	err = c.get(ctx, "scan_forms/"+scanFormID, &out)
	return
}
