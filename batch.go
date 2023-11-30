package easypost

import (
	"context"
	"net/http"
	"net/url"
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
	CreatedAt    *DateTime    `json:"created_at,omitempty"`
	UpdatedAt    *DateTime    `json:"updated_at,omitempty"`
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

// ListBatchesResult holds the results from the list insurances API.
type ListBatchesResult struct {
	Batch []*Batch `json:"batches,omitempty"`
	PaginatedCollection
}

type addRemoveShipmentsRequest struct {
	Shipments []*Shipment `json:"shipments,omitempty"`
}

// CreateBatch creates a new batch of shipments. It optionally accepts one or
// more shipments to add to the new batch. If successful, a new batch object is
// returned.
//
//	c := easypost.New(MyEasyPostAPIKey)
//	out, err := c.CreateBatch(
//		&easypost.Shipment{ID: "shp_100"},
//		&easypost.Shipment{ID: "shp_101"},
//		&easypost.Shipment{ID: "shp_102"},
//	)
func (c *Client) CreateBatch(in ...*Shipment) (out *Batch, err error) {
	return c.CreateBatchWithContext(context.Background(), in...)
}

// CreateBatchWithContext performs the same operation as CreateBatch, but allows
// specifying a context that can interrupt the request.
func (c *Client) CreateBatchWithContext(ctx context.Context, in ...*Shipment) (out *Batch, err error) {
	req := batchRequest{Batch: &Batch{Shipments: in}}
	err = c.post(ctx, "batches", req, &out)
	return
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

// TODO: Add support for GetNextPage when the API supports it.

// AddShipmentsToBatch adds shipments to an existing batch, and returns the
// updated batch object.
func (c *Client) AddShipmentsToBatch(batchID string, shipments ...*Shipment) (out *Batch, err error) {
	return c.AddShipmentsToBatchWithContext(context.Background(), batchID, shipments...)
}

// AddShipmentsToBatchWithContext performs the same operation as
// AddShipmentsToBatch, but allows specifying a context that can interrupt the
// request.
func (c *Client) AddShipmentsToBatchWithContext(ctx context.Context, batchID string, shipments ...*Shipment) (out *Batch, err error) {
	req := addRemoveShipmentsRequest{Shipments: shipments}
	err = c.post(ctx, "batches/"+batchID+"/add_shipments", req, &out)
	return
}

// RemoveShipmentsFromBatch removes shipments from an existing batch, and
// returns the updated batch object.
func (c *Client) RemoveShipmentsFromBatch(batchID string, shipments ...*Shipment) (out *Batch, err error) {
	return c.RemoveShipmentsFromBatchWithContext(context.Background(), batchID, shipments...)
}

// RemoveShipmentsFromBatchWithContext performs the same operation as
// RemoveShipmentsFromBatch, but allows specifying a context that can interrupt
// the request.
func (c *Client) RemoveShipmentsFromBatchWithContext(ctx context.Context, batchID string, shipments ...*Shipment) (out *Batch, err error) {
	req := addRemoveShipmentsRequest{Shipments: shipments}
	err = c.post(ctx, "batches/"+batchID+"/remove_shipments", req, &out)
	return
}

// BuyBatch initializes purchases for the shipments in the batch. The updated
// batch object is returned.
func (c *Client) BuyBatch(batchID string) (out *Batch, err error) {
	return c.BuyBatchWithContext(context.Background(), batchID)
}

// BuyBatchWithContext performs the same operation as BuyBatch, but allows
// specifying a context that can interrupt the request.
func (c *Client) BuyBatchWithContext(ctx context.Context, batchID string) (out *Batch, err error) {
	err = c.post(ctx, "batches/"+batchID+"/buy", nil, &out)
	return
}

// GetBatch retrieves a Batch object by ID.
func (c *Client) GetBatch(batchID string) (out *Batch, err error) {
	return c.GetBatchWithContext(context.Background(), batchID)
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
	return c.GetBatchLabelsWithContext(context.Background(), batchID, format)
}

// GetBatchLabelsWithContext performs the same operation as GetBatchLabels, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetBatchLabelsWithContext(ctx context.Context, batchID, format string) (out *Batch, err error) {
	params := url.Values{"file_format": []string{format}}
	err = c.post(ctx, "batches/"+batchID+"/label", params, &out)
	return
}

// CreateBatchScanForms generates a scan form for the batch.
func (c *Client) CreateBatchScanForms(batchID, format string) (out *Batch, err error) {
	return c.CreateBatchScanFormsWithContext(context.Background(), batchID, format)
}

// CreateBatchScanFormsWithContext performs the same operation as
// CreateBatchScanForms, but allows specifying a context that can interrupt the
// request.
func (c *Client) CreateBatchScanFormsWithContext(ctx context.Context, batchID, format string) (out *Batch, err error) {
	vals := url.Values{"file_format": []string{format}}
	err = c.do(ctx, http.MethodPost, "batches/"+batchID+"/scan_form", vals, &out)
	return
}
