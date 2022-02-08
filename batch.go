package easypost

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// BatchStatus contains counts of statuses for the shipments in a batch.
type BatchStatus struct {
	PostagePurchased      int `json:"postage_purchased,omitempty"`
	PostagePurchaseFailed int `json:"postage_purchase_failed,omitempty"`
	QueuedForPurchase     int `json:"queued_for_purchase,omitempty"`
	CreationFailed        int `json:"creation_failed,omitempty"`
}

// Batch represents a batch of shipments.
type Batch struct {
	ID           string       `json:"id,omitempty"`
	Object       string       `json:"object,omitempty"`
	Reference    string       `json:"reference,omitempty"`
	Mode         string       `json:"mode,omitempty"`
	CreatedAt    *time.Time   `json:"created_at,omitempty"`
	UpdatedAt    *time.Time   `json:"updated_at,omitempty"`
	State        string       `json:"state,omitempty"`
	NumShipments int          `json:"num_shipments,omitempty"`
	Shipments    []*Shipment  `json:"shipments,omitempty"`
	Status       *BatchStatus `json:"status,omitempty"`
	LabelURL     string       `json:"label_url,omitempty"`
	ScanForm     *ScanForm    `json:"scan_form,omitempty"`
	Pickup       *Pickup      `json:"pickup,omitempty"`
}

type batchRequest struct {
	Batch *Batch `json:"batch,omitempty"`
}

// CreateBatch creates a new batch of shipments. It optionally accepts one or
// more shipments to add to the new batch. If successful, a new batch object is
// returned.
//	c := easypost.New(MyEasyPostAPIKey)
//	out, err := c.CreateBatch(
//		&easypost.Shipment{ID: "shp_100"},
//		&easypost.Shipment{ID: "shp_101"},
//		&easypost.Shipment{ID: "shp_102"},
//	)
func (c *Client) CreateBatch(in ...*Shipment) (out *Batch, err error) {
	req := batchRequest{Batch: &Batch{Shipments: in}}
	err = c.post(context.Background(), "batches", req, &out)
	return
}

// CreateBatchWithContext performs the same operation as CreateBatch, but allows
// specifying a context that can interrupt the request.
func (c *Client) CreateBatchWithContext(ctx context.Context, in ...*Shipment) (out *Batch, err error) {
	req := batchRequest{Batch: &Batch{Shipments: in}}
	err = c.post(ctx, "batches", req, &out)
	return
}

// CreateAndBuyBatch creates and buys a new batch of shipments in one request.
func (c *Client) CreateAndBuyBatch(in ...*Shipment) (out *Batch, err error) {
	req := batchRequest{Batch: &Batch{Shipments: in}}
	err = c.post(context.Background(), "batches/create_and_buy", req, &out)
	return
}

// CreateAndBuyBatchWithContext performs the same operation as
// CreateAndBuyBatch, but allows specifying a context that can interrupt the
// request.
func (c *Client) CreateAndBuyBatchWithContext(ctx context.Context, in ...*Shipment) (out *Batch, err error) {
	req := batchRequest{Batch: &Batch{Shipments: in}}
	err = c.post(ctx, "batches/create_and_buy", req, &out)
	return
}

// ListBatchesResult holds the results from the list insurances API.
type ListBatchesResult struct {
	Insurances []*Insurance `json:"batches,omitempty"`
	// HasMore indicates if there are more responses to be fetched. If True,
	// additional responses can be fetched by updating the ListInsurancesOptions
	// parameter's AfterID field with the ID of the last item in this object's
	// Insurances field.
	HasMore bool `json:"has_more,omitempty"`
}

// ListBatches provides a paginated result of Insurance objects.
func (c *Client) ListBatches(opts *ListOptions) (out *ListBatchesResult, err error) {
	return c.ListBatchesWithContext(context.Background(), opts)
}

// ListBatchesWithContext performs the same operation as ListBatches, but
// allows specifying a context that can interrupt the request.
func (c *Client) ListBatchesWithContext(ctx context.Context, opts *ListOptions) (out *ListBatchesResult, err error) {
	err = c.do(ctx, http.MethodGet, "batches", c.convertOptsToURLValues(opts), &out)
	return
}

// AddShipmentsToBatch adds shipments to an existing batch, and returns the
// updated batch object.
func (c *Client) AddShipmentsToBatch(batchID string, in ...*Shipment) (out *Batch, err error) {
	req := batchRequest{Batch: &Batch{Shipments: in}}
	err = c.post(context.Background(), "batches/"+batchID+"/add_shipments", req, &out)
	return
}

// AddShipmentsToBatchWithContext performs the same operation as
// AddShipmentsToBatch, but allows specifying a context that can interrupt the
// request.
func (c *Client) AddShipmentsToBatchWithContext(ctx context.Context, batchID string, in ...*Shipment) (out *Batch, err error) {
	req := batchRequest{Batch: &Batch{Shipments: in}}
	err = c.post(ctx, "batches/"+batchID+"/add_shipments", req, &out)
	return
}

// RemoveShipmentsFromBatch removes shipments fro, an existing batch, and
// returns the updated batch object.
func (c *Client) RemoveShipmentsFromBatch(batchID string, shipmentIDs ...string) (out *Batch, err error) {
	req := batchRequest{
		Batch: &Batch{Shipments: make([]*Shipment, len(shipmentIDs))},
	}
	for i := range shipmentIDs {
		req.Batch.Shipments[i] = &Shipment{ID: shipmentIDs[i]}
	}
	err = c.post(context.Background(), "batches/"+batchID+"/remove_shipments", req, &out)
	return
}

// RemoveShipmentsFromBatchWithContext performs the same operation as
// RemoveShipmentsFromBatch, but allows specifying a context that can interrupt
// the request.
func (c *Client) RemoveShipmentsFromBatchWithContext(ctx context.Context, batchID string, shipmentIDs ...string) (out *Batch, err error) {
	req := batchRequest{
		Batch: &Batch{Shipments: make([]*Shipment, len(shipmentIDs))},
	}
	for i := range shipmentIDs {
		req.Batch.Shipments[i] = &Shipment{ID: shipmentIDs[i]}
	}
	err = c.post(ctx, "batches/"+batchID+"/remove_shipments", req, &out)
	return
}

// BuyBatch initializes purchases for the shipments in the batch. The updated
// batch object is returned.
func (c *Client) BuyBatch(batchID string) (out *Batch, err error) {
	err = c.post(context.Background(), "batches/"+batchID+"/buy", nil, &out)
	return
}

// BuyBatchWithContext performs the same operation as BuyBatch, but allows
// specifying a context that can interrupt the request.
func (c *Client) BuyBatchWithContext(ctx context.Context, batchID string) (out *Batch, err error) {
	err = c.post(ctx, "batches/"+batchID+"/buy", nil, &out)
	return
}

// GetBatch retrieves a Batch object by ID.
func (c *Client) GetBatch(batchID string) (out *Batch, err error) {
	err = c.get(context.Background(), "batches/"+batchID, &out)
	return
}

// GetBatchWithContext performs the same operation as GetBatch, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetBatchWithContext(ctx context.Context, batchID string) (out *Batch, err error) {
	err = c.get(ctx, "batches/"+batchID, &out)
	return
}

// GetBatchLabels generates a label for the batch. This can only be done once
// per batch, and all shipments must have a "postage_purchased" status.
func (c *Client) GetBatchLabels(batchID, format string) (out *Batch, err error) {
	vals := url.Values{"file_format": []string{format}}
	err = c.do(context.Background(), http.MethodGet, "batches/"+batchID+"/label", vals, &out)
	return
}

// GetBatchLabelsWithContext performs the same operation as GetBatchLabels, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetBatchLabelsWithContext(ctx context.Context, batchID, format string) (out *Batch, err error) {
	vals := url.Values{"file_format": []string{format}}
	err = c.do(ctx, http.MethodGet, "batches/"+batchID+"/label", vals, &out)
	return
}

// CreateBatchScanForms generates a scan form for the batch.
func (c *Client) CreateBatchScanForms(batchID, format string) (out *Batch, err error) {
	vals := url.Values{"file_format": []string{format}}
	err = c.do(context.Background(), http.MethodPost, "batches/"+batchID+"/scan_form", vals, &out)
	return
}

// CreateBatchScanFormsWithContext performs the same operation as
// CreateBatchScanForms, but allows specifying a context that can interrupt the
// request.
func (c *Client) CreateBatchScanFormsWithContext(ctx context.Context, batchID, format string) (out *Batch, err error) {
	vals := url.Values{"file_format": []string{format}}
	err = c.do(ctx, http.MethodPost, "batches/"+batchID+"/scan_form", vals, &out)
	return
}
