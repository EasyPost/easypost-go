package easypost

import (
	"context"
	"net/http"
)

type ClaimHistoryEntry struct {
	Status          string `json:"status,omitempty" url:"status,omitempty"`
	StatusDetail    string `json:"status_detail,omitempty" url:"status_detail,omitempty"`
	StatusTimestamp string `json:"status_timestamp,omitempty" url:"status_timestamp,omitempty"`
}

// A Claim object represents a claim against insurance on a package.
type Claim struct {
	ID                   string              `json:"id,omitempty" url:"id,omitempty"`
	Object               string              `json:"object,omitempty" url:"object,omitempty"`
	Mode                 string              `json:"mode,omitempty" url:"mode,omitempty"`
	CreatedAt            *DateTime           `json:"created_at,omitempty" url:"created_at,omitempty"`
	UpdatedAt            *DateTime           `json:"updated_at,omitempty" url:"updated_at,omitempty"`
	ApprovedAmount       string              `json:"approved_amount,omitempty" url:"approved_amount,omitempty"`
	Attachments          []string            `json:"attachments,omitempty" url:"attachments,omitempty"`
	CheckDeliveryAddress string              `json:"check_delivery_address,omitempty" url:"check_delivery_address,omitempty"`
	ContactEmail         string              `json:"contact_email,omitempty" url:"contact_email,omitempty"`
	Description          string              `json:"description,omitempty" url:"description,omitempty"`
	History              []ClaimHistoryEntry `json:"history,omitempty" url:"history,omitempty"`
	InsuranceAmount      string              `json:"insurance_amount,omitempty" url:"insurance_amount,omitempty"`
	InsuranceId          string              `json:"insurance_id,omitempty" url:"insurance_id,omitempty"`
	PaymentMethod        string              `json:"payment_method,omitempty" url:"payment_method,omitempty"`
	RecipientName        string              `json:"recipient_name,omitempty" url:"recipient_name,omitempty"`
	Reference            string              `json:"reference,omitempty" url:"reference,omitempty"`
	RequestedAmount      string              `json:"requested_amount,omitempty" url:"requested_amount,omitempty"`
	SalvageValue         string              `json:"salvage_value,omitempty" url:"salvage_value,omitempty"`
	ShipmentId           string              `json:"shipment_id,omitempty" url:"shipment_id,omitempty"`
	Status               string              `json:"status,omitempty" url:"status,omitempty"`
	StatusDetail         string              `json:"status_detail,omitempty" url:"status_detail,omitempty"`
	StatusTimestamp      string              `json:"status_timestamp,omitempty" url:"status_timestamp,omitempty"`
	TrackingCode         string              `json:"tracking_code,omitempty" url:"tracking_code,omitempty"`
	Type                 string              `json:"type,omitempty" url:"type,omitempty"`
}

// CreateClaimParameters is used to specify parameters for creating a claim.
type CreateClaimParameters struct {
	TrackingCode                       string   `json:"tracking_code,omitempty" url:"tracking_code,omitempty"`
	Type                               string   `json:"type,omitempty" url:"type,omitempty"`
	Amount                             float64  `json:"amount,omitempty" url:"amount,omitempty"`
	EmailEvidenceAttachments           []string `json:"email_evidence_attachments,omitempty" url:"email_evidence_attachments,omitempty"`
	InvoiceAttachments                 []string `json:"invoice_attachments,omitempty" url:"invoice_attachments,omitempty"`
	SupportingDocumentationAttachments []string `json:"supporting_documentation_attachments,omitempty" url:"supporting_documentation_attachments,omitempty"`
	Description                        string   `json:"description,omitempty" url:"description,omitempty"`
	RecipientName                      string   `json:"recipient_name,omitempty" url:"recipient_name,omitempty"`
	ContactEmail                       string   `json:"contact_email,omitempty" url:"contact_email,omitempty"`
	PaymentMethod                      string   `json:"payment_method,omitempty" url:"payment_method,omitempty"`
	Reference                          string   `json:"reference,omitempty" url:"reference,omitempty"`
}

// ListClaimsParameters is used to specify query parameters for listing claims.
type ListClaimsParameters struct {
	BeforeID      string    `json:"before_id,omitempty" url:"before_id,omitempty"`
	AfterID       string    `json:"after_id,omitempty" url:"after_id,omitempty"`
	StartDateTime *DateTime `json:"start_datetime,omitempty" url:"start_datetime,omitempty"`
	EndDateTime   *DateTime `json:"end_datetime,omitempty" url:"end_datetime,omitempty"`
	PageSize      int       `json:"page_size,omitempty" url:"page_size,omitempty"`
	Type          string    `json:"type,omitempty" url:"type,omitempty"`
	Status        string    `json:"status,omitempty" url:"status,omitempty"`
}

// ListClaimsResult holds the results from the list claims API.
type ListClaimsResult struct {
	Claims     []*Claim `json:"claims,omitempty" url:"claims,omitempty"`
	Parameters ListClaimsParameters
	PaginatedCollection
}

// CreateClaim creates a Claim object for insurance purchased through EasyPost.
func (c *Client) CreateClaim(in *CreateClaimParameters) (out *Claim, err error) {
	return c.CreateClaimWithContext(context.Background(), in)
}

// CreateClaimWithContext performs the same operation as CreateClaim, but allows specifying a context that can interrupt the request.
func (c *Client) CreateClaimWithContext(ctx context.Context, in *CreateClaimParameters) (out *Claim, err error) {
	err = c.do(ctx, http.MethodPost, "claims", in, &out)
	return
}

// ListClaims provides a paginated result of Claim objects.
func (c *Client) ListClaims(opts *ListClaimsParameters) (out *ListClaimsResult, err error) {
	return c.ListClaimsWithContext(context.Background(), opts)
}

// ListClaimsWithContext performs the same operation as ListClaims, but allows specifying a context that can interrupt the request.
func (c *Client) ListClaimsWithContext(ctx context.Context, opts *ListClaimsParameters) (out *ListClaimsResult, err error) {
	err = c.do(ctx, http.MethodGet, "claims", opts, &out)
	// Store the original query parameters for reuse when getting the next page
	out.Parameters = *opts
	return
}

// GetNextClaimPage returns the next page of claims
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
		err = newEndOfPaginationError()
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
	err = c.do(ctx, http.MethodGet, "claims/"+claimID, nil, &out)
	return
}

// CancelClaim refunds the Claim object with the given ID.
func (c *Client) CancelClaim(claimID string) (out *Claim, err error) {
	return c.CancelClaimWithContext(context.Background(), claimID)
}

// CancelClaimWithContext performs the same operation as CancelClaim, but allows specifying a context that can interrupt the request.
func (c *Client) CancelClaimWithContext(ctx context.Context, claimID string) (out *Claim, err error) {
	err = c.do(ctx, http.MethodPost, "claims/"+claimID+"/cancel", nil, &out)
	return
}
