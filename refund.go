package easypost

import (
	"context"
	"net/http"
)

type Refund struct {
	ID                 string    `json:"id,omitempty" url:"id,omitempty"`
	Object             string    `json:"object,omitempty" url:"object,omitempty"`
	CreatedAt          *DateTime `json:"created_at,omitempty" url:"created_at,omitempty"`
	UpdatedAt          *DateTime `json:"updated_at,omitempty" url:"updated_at,omitempty"`
	TrackingCode       string    `json:"tracking_code,omitempty" url:"tracking_code,omitempty"`
	ConfirmationNumber string    `json:"confirmation_number,omitempty" url:"confirmation_number,omitempty"`
	Status             string    `json:"status,omitempty" url:"status,omitempty"`
	Carrier            string    `json:"carrier,omitempty" url:"carrier,omitempty"`
	ShipmentID         string    `json:"shipment_id,omitempty" url:"shipment_id,omitempty"`
}

type ListRefundResult struct {
	Refunds []*Refund `json:"refunds,omitempty" url:"refunds,omitempty"`
	PaginatedCollection
}

// CreateRefund submits a refund request and return a list of refunds.
func (c *Client) CreateRefund(in map[string]interface{}) (out []*Refund, err error) {
	return c.CreateRefundWithContext(context.Background(), in)
}

// CreateRefundWithContext performs the same operation as CreateRefund, but
// allows specifying a context that can interrupt the request.
func (c *Client) CreateRefundWithContext(ctx context.Context, in map[string]interface{}) (out []*Refund, err error) {
	req := map[string]interface{}{"refund": in}
	err = c.do(ctx, http.MethodPost, "refunds", req, &out)
	return
}

// ListRefunds provides a paginated result of Refund objects.
func (c *Client) ListRefunds(opts *ListOptions) (out *ListRefundResult, err error) {
	return c.ListRefundsWithContext(context.Background(), opts)
}

// ListRefundsWithContext performs the same operation as ListRefunds, but
// allows specifying a context that can interrupt the request.
func (c *Client) ListRefundsWithContext(ctx context.Context, opts *ListOptions) (out *ListRefundResult, err error) {
	err = c.do(ctx, http.MethodGet, "refunds", opts, &out)
	return
}

// GetNextRefundPage returns the next page of refunds
func (c *Client) GetNextRefundPage(collection *ListRefundResult) (out *ListRefundResult, err error) {
	return c.GetNextRefundPageWithContext(context.Background(), collection)
}

// GetNextRefundPageWithPageSize returns the next page of refunds with a specific page size
func (c *Client) GetNextRefundPageWithPageSize(collection *ListRefundResult, pageSize int) (out *ListRefundResult, err error) {
	return c.GetNextRefundPageWithPageSizeWithContext(context.Background(), collection, pageSize)
}

// GetNextRefundPageWithContext performs the same operation as GetNextRefundPage, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetNextRefundPageWithContext(ctx context.Context, collection *ListRefundResult) (out *ListRefundResult, err error) {
	return c.GetNextRefundPageWithPageSizeWithContext(ctx, collection, 0)
}

// GetNextRefundPageWithPageSizeWithContext performs the same operation as GetNextRefundPageWithPageSize, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetNextRefundPageWithPageSizeWithContext(ctx context.Context, collection *ListRefundResult, pageSize int) (out *ListRefundResult, err error) {
	if len(collection.Refunds) == 0 {
		err = EndOfPaginationError
		return
	}
	lastID := collection.Refunds[len(collection.Refunds)-1].ID
	params, err := nextPageParameters(collection.HasMore, lastID, pageSize)
	if err != nil {
		return
	}
	return c.ListRefundsWithContext(ctx, params)
}

// retrieves a previously-created Refund by its ID.
func (c *Client) GetRefund(refundID string) (out *Refund, err error) {
	return c.GetRefundWithContext(context.Background(), refundID)
}

// GetRefundWithContext performs the same operation as GetRefund, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetRefundWithContext(ctx context.Context, refundID string) (out *Refund, err error) {
	err = c.do(ctx, http.MethodGet, "refunds/"+refundID, nil, &out)
	return
}
