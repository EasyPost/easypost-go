package easypost

import (
	"context"
	"net/http"
	"time"
)

type Refund struct {
	ID                 string     `json:"id,omitempty"`
	Object             string     `json:"object,omitempty"`
	CreatedAt          *time.Time `json:"created_at,omitempty"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
	TrackingCode       string     `json:"tracking_code,omitempty"`
	ConfirmationNumber string     `json:"confirmation_number,omitempty"`
	Status             string     `json:"status,omitempty"`
	Carrier            string     `json:"carrier,omitempty"`
	ShipmentID         string     `json:"shipment_id,omitempty"`
}

type ListRefundResult struct {
	Refunds []*Refund `json:"refunds,omitempty"`
	HasMore bool      `json:"has_more,omitempty"`
}

// CreateRefund submits a refund request and return a list of refunds.
func (c *Client) CreateRefund(in map[string]interface{}) (out []*Refund, err error) {
	req := map[string]interface{}{"refund": in}
	err = c.post(context.Background(), "refunds", req, &out)
	return
}

// CreateRefundWithContext performs the same operation as CreateRefund, but
// allows specifying a context that can interrupt the request.
func (c *Client) CreateRefundWithContext(ctx context.Context, in map[string]interface{}) (out []*Refund, err error) {
	req := map[string]interface{}{"refund": in}
	err = c.post(ctx, "refunds", req, &out)
	return
}

// ListRefunds provides a paginated result of Refund objects.
func (c *Client) ListRefunds(opts *ListOptions) (out *ListRefundResult, err error) {
	err = c.do(context.Background(), http.MethodGet, "refunds", c.convertOptsToURLValues(opts), &out)
	return
}

// ListRefundsWithContext performs the same operation as ListRefunds, but
// allows specifying a context that can interrupt the request.
func (c *Client) ListRefundsWithContext(ctx context.Context, opts *ListOptions) (out *ListRefundResult, err error) {
	err = c.do(ctx, http.MethodGet, "refunds", c.convertOptsToURLValues(opts), &out)
	return
}

// retrieves a previously-created Refund by its ID.
func (c *Client) GetRefund(refundID string) (out *Refund, err error) {
	err = c.get(context.Background(), "refunds/"+refundID, &out)
	return
}

// GetRefundWithContext performs the same operation as GetRefund, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetRefundWithContext(ctx context.Context, refundID string) (out *Refund, err error) {
	err = c.get(ctx, "refunds/"+refundID, &out)
	return
}
