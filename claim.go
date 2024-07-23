package easypost

import (
	"context"
	"net/http"
)

type ClaimHistoryEntry struct {
	Status          string `json:"status,omitempty"`
	StatusDetail    string `json:"status_detail,omitempty"`
	StatusTimestamp string `json:"status_timestamp,omitempty"`
}

// A Claim object represents a claim against insurance on a package.
type Claim struct {
	ID                   string              `json:"id,omitempty"`
	Object               string              `json:"object,omitempty"`
	Mode                 string              `json:"mode,omitempty"`
	CreatedAt            *DateTime           `json:"created_at,omitempty"`
	UpdatedAt            *DateTime           `json:"updated_at,omitempty"`
	ApprovedAmount       string              `json:"approved_amount,omitempty"`
	Attachments          []string            `json:"attachments,omitempty"`
	CheckDeliveryAddress string              `json:"check_delivery_address,omitempty"`
	ContactEmail         string              `json:"contact_email,omitempty"`
	Description          string              `json:"description,omitempty"`
	History              []ClaimHistoryEntry `json:"history,omitempty"`
	InsuranceAmount      string              `json:"insurance_amount,omitempty"`
	InsuranceId          string              `json:"insurance_id,omitempty"`
	PaymentMethod        string              `json:"payment_method,omitempty"`
	RecipientName        string              `json:"recipient_name,omitempty"`
	RequestedAmount      string              `json:"requested_amount,omitempty"`
	SalvageValue         string              `json:"salvage_value,omitempty"`
	ShipmentId           string              `json:"shipment_id,omitempty"`
	Status               string              `json:"status,omitempty"`
	StatusDetail         string              `json:"status_detail,omitempty"`
	StatusTimestamp      string              `json:"status_timestamp,omitempty"`
	TrackingCode         string              `json:"tracking_code,omitempty"`
	Type                 string              `json:"type,omitempty"`
}

// CreateClaimParameters is used to specify parameters for creating a claim.
type CreateClaimParameters struct {
	TrackingCode                       string   `json:"tracking_code,omitempty"`
	Type                               string   `json:"type,omitempty"`
	Amount                             float64  `json:"amount,omitempty"`
	EmailEvidenceAttachments           []string `json:"email_evidence_attachments,omitempty"`
	InvoiceAttachments                 []string `json:"invoice_attachments,omitempty"`
	SupportingDocumentationAttachments []string `json:"supporting_documentation_attachments,omitempty"`
	Description                        string   `json:"description,omitempty"`
	RecipientName                      string   `json:"recipient_name,omitempty"`
	ContactEmail                       string   `json:"contact_email,omitempty"`
	PaymentMethod                      string   `json:"payment_method,omitempty"`
}

// ListClaimsParameters is used to specify query parameters for listing claims.
type ListClaimsParameters struct {
	BeforeID      string    `url:"before_id,omitempty"`
	AfterID       string    `url:"after_id,omitempty"`
	StartDateTime *DateTime `url:"start_datetime,omitempty"`
	EndDateTime   *DateTime `url:"end_datetime,omitempty"`
	PageSize      int       `url:"page_size,omitempty"`
	Type          string    `url:"type,omitempty"`
	Status        string    `url:"status,omitempty"`
}

// ListClaimsResult holds the results from the list claims API.
type ListClaimsResult struct {
	Claims     []*Claim `json:"claims,omitempty"`
	Parameters ListClaimsParameters
	PaginatedCollection
}

// CreateClaim creates a Claim object for insurance purchased through EasyPost.
func (c *Client) CreateClaim(in *CreateClaimParameters) (out *Claim, err error) {
	return c.CreateClaimWithContext(context.Background(), in)
}

// CreateClaimWithContext performs the same operation as CreateClaim, but allows specifying a context that can interrupt the request.
func (c *Client) CreateClaimWithContext(ctx context.Context, in *CreateClaimParameters) (out *Claim, err error) {
	err = c.post(ctx, "claims", in, &out)
	return
}

// ListClaims provides a paginated result of Claim objects.
func (c *Client) ListClaims(opts *ListClaimsParameters) (out *ListClaimsResult, err error) {
	return c.ListClaimsWithContext(context.Background(), opts)
}

// ListClaimsWithContext performs the same operation as ListClaims, but allows specifying a context that can interrupt the request.
func (c *Client) ListClaimsWithContext(ctx context.Context, opts *ListClaimsParameters) (out *ListClaimsResult, err error) {
	err = c.do(ctx, http.MethodGet, "claims", c.convertOptsToURLValues(opts), &out)
	// Store the original query parameters for reuse when getting the next page
	out.Parameters = *opts
	return
}

// GetNextClaimPage returns the next page of shipments
func (c *Client) GetNextClaimPage(collection *ListClaimsResult) (out *ListClaimsResult, err error) {
	return c.GetNextClaimPageWithContext(context.Background(), collection)
}

// GetNextClaimPageWithPageSize returns the next page of claims with a specific page size
func (c *Client) GetNextClaimPageWithPageSize(collection *ListClaimsResult, pageSize int) (out *ListClaimsResult, err error) {
	return c.GetNextClaimPageWithPageSizeWithContext(context.Background(), collection, pageSize)
}

// GetNextClaimPageWithContext performs the same operation as GetNextClaimPage, but allows specifying a context that can interrupt the request.
func (c *Client) GetNextClaimPageWithContext(ctx context.Context, collection *ListClaimsResult) (out *ListClaimsResult, err error) {
	return c.GetNextClaimPageWithPageSizeWithContext(ctx, collection, 0)
}

// GetNextClaimPageWithPageSizeWithContext performs the same operation as GetNextClaimPageWithPageSize, but allows specifying a context that can interrupt the request.
func (c *Client) GetNextClaimPageWithPageSizeWithContext(ctx context.Context, collection *ListClaimsResult, pageSize int) (out *ListClaimsResult, err error) {
	if len(collection.Claims) == 0 {
		err = EndOfPaginationError
		return
	}
	lastID := collection.Claims[len(collection.Claims)-1].ID
	params, err := nextPageParameters(collection.HasMore, lastID, pageSize)
	if err != nil {
		return
	}
	claimParams := &collection.Parameters
	claimParams.BeforeID = params.BeforeID
	if pageSize > 0 {
		claimParams.PageSize = pageSize
	}
	return c.ListClaimsWithContext(ctx, claimParams)
}

// GetClaim returns the Claim object with the given ID or reference.
func (c *Client) GetClaim(claimID string) (out *Claim, err error) {
	return c.GetClaimWithContext(context.Background(), claimID)
}

// GetClaimWithContext performs the same operation as GetClaim, but allows specifying a context that can interrupt the request.
func (c *Client) GetClaimWithContext(ctx context.Context, claimID string) (out *Claim, err error) {
	err = c.get(ctx, "claims/"+claimID, &out)
	return
}

// CancelClaim refunds the Claim object with the given ID.
func (c *Client) CancelClaim(claimID string) (out *Claim, err error) {
	return c.CancelClaimWithContext(context.Background(), claimID)
}

// CancelClaimWithContext performs the same operation as CancelClaim, but allows specifying a context that can interrupt the request.
func (c *Client) CancelClaimWithContext(ctx context.Context, claimID string) (out *Claim, err error) {
	err = c.post(ctx, "claims/"+claimID+"/cancel", nil, &out)
	return
}
