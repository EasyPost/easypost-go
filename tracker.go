package easypost

import (
	"context"
	"net/http"
)

// TrackingLocation provides additional information about the location of a
// tracking event.
type TrackingLocation struct {
	Object  string `json:"object,omitempty" url:"object,omitempty"`
	City    string `json:"city,omitempty" url:"city,omitempty"`
	State   string `json:"state,omitempty" url:"state,omitempty"`
	Country string `json:"country,omitempty" url:"country,omitempty"`
	Zip     string `json:"zip,omitempty" url:"zip,omitempty"`
}

// TrackingDetail provides information about a tracking event.
type TrackingDetail struct {
	Object           string            `json:"object,omitempty" url:"object,omitempty"`
	Message          string            `json:"message,omitempty" url:"message,omitempty"`
	Description      string            `json:"description,omitempty" url:"description,omitempty"`
	Status           string            `json:"status,omitempty" url:"status,omitempty"`
	StatusDetail     string            `json:"status_detail,omitempty" url:"status_detail,omitempty"`
	DateTime         string            `json:"datetime,omitempty" url:"datetime,omitempty"`
	Source           string            `json:"source,omitempty" url:"source,omitempty"`
	CarrierCode      string            `json:"carrier_code,omitempty" url:"carrier_code,omitempty"`
	TrackingLocation *TrackingLocation `json:"tracking_location,omitempty" url:"tracking_location,omitempty"`
	EstDeliveryDate  *DateTime         `json:"est_delivery_date,omitempty" url:"est_delivery_date,omitempty"`
}

// TrackingCarrierDetail provides additional tracking information from the
// carrier, when available.
type TrackingCarrierDetail struct {
	Object                      string            `json:"object,omitempty" url:"object,omitempty"`
	Service                     string            `json:"service,omitempty" url:"service,omitempty"`
	ContainerType               string            `json:"container_type,omitempty" url:"container_type,omitempty"`
	EstDeliveryDateLocal        string            `json:"est_delivery_date_local,omitempty" url:"est_delivery_date_local,omitempty"`
	EstDeliveryTimeLocal        string            `json:"est_delivery_time_local,omitempty" url:"est_delivery_time_local,omitempty"`
	OriginLocation              string            `json:"origin_location,omitempty" url:"origin_location,omitempty"`
	OriginTrackingLocation      *TrackingLocation `json:"origin_tracking_location,omitempty" url:"origin_tracking_location,omitempty"`
	DestinationLocation         string            `json:"destination_location,omitempty" url:"destination_location,omitempty"`
	DestinationTrackingLocation *TrackingLocation `json:"destination_tracking_location,omitempty" url:"destination_tracking_location,omitempty"`
	GuaranteedDeliveryDate      *DateTime         `json:"guaranteed_delivery_date,omitempty" url:"guaranteed_delivery_date,omitempty"`
	AlternateIdentifier         string            `json:"alternate_identifier,omitempty" url:"alternate_identifier,omitempty"`
	InitialDeliveryAttempt      *DateTime         `json:"initial_delivery_attempt,omitempty" url:"initial_delivery_attempt,omitempty"`
}

// A Tracker object contains all the tracking information for a package.
type Tracker struct {
	ID              string                 `json:"id,omitempty" url:"id,omitempty"`
	Object          string                 `json:"object,omitempty" url:"object,omitempty"`
	Mode            string                 `json:"mode,omitempty" url:"mode,omitempty"`
	CreatedAt       *DateTime              `json:"created_at,omitempty" url:"created_at,omitempty"`
	UpdatedAt       *DateTime              `json:"updated_at,omitempty" url:"updated_at,omitempty"`
	TrackingCode    string                 `json:"tracking_code,omitempty" url:"tracking_code,omitempty"`
	Status          string                 `json:"status,omitempty" url:"status,omitempty"`
	StatusDetail    string                 `json:"status_detail,omitempty" url:"status_detail,omitempty"`
	SignedBy        string                 `json:"signed_by,omitempty" url:"signed_by,omitempty"`
	Weight          float64                `json:"weight,omitempty" url:"weight,omitempty"`
	EstDeliveryDate *DateTime              `json:"est_delivery_date,omitempty" url:"est_delivery_date,omitempty"`
	ShipmentID      string                 `json:"shipment_id,omitempty" url:"shipment_id,omitempty"`
	Carrier         string                 `json:"carrier,omitempty" url:"carrier,omitempty"`
	TrackingDetails []*TrackingDetail      `json:"tracking_details,omitempty" url:"tracking_details,omitempty"`
	CarrierDetail   *TrackingCarrierDetail `json:"carrier_detail,omitempty" url:"carrier_detail,omitempty"`
	PublicURL       string                 `json:"public_url,omitempty" url:"public_url,omitempty"`
	Fees            []*Fee                 `json:"fees,omitempty" url:"fees,omitempty"`
	Finalized       bool                   `json:"finalized,omitempty" url:"finalized,omitempty"`
	IsReturn        bool                   `json:"is_return,omitempty" url:"is_return,omitempty"`
}

// CreateTrackerOptions specifies options for creating a new tracker.
type CreateTrackerOptions struct {
	TrackingCode   string
	Carrier        string
	Amount         string
	CarrierAccount string
	IsReturn       bool
}

// ListTrackersOptions is used to specify query parameters for listing Tracker
// objects.
type ListTrackersOptions struct {
	BeforeID      string    `json:"before_id,omitempty" url:"before_id,omitempty"`
	AfterID       string    `json:"after_id,omitempty" url:"after_id,omitempty"`
	StartDateTime *DateTime `json:"start_datetime,omitempty" url:"start_datetime,omitempty"`
	EndDateTime   *DateTime `json:"end_datetime,omitempty" url:"end_datetime,omitempty"`
	PageSize      int       `json:"page_size,omitempty" url:"page_size,omitempty"`
	TrackingCode  string    `json:"tracking_code,omitempty" url:"tracking_code,omitempty"`
	TrackingCodes []string  `json:"tracking_codes,omitempty" url:"tracking_codes,omitempty"`
	Carrier       string    `json:"carrier,omitempty" url:"carrier,omitempty"`
}

// ListTrackersResult holds the results from the list trackers API.
type ListTrackersResult struct {
	Trackers     []*Tracker `json:"trackers,omitempty" url:"trackers,omitempty"`
	TrackingCode string     `json:"tracking_code,omitempty" url:"tracking_code,omitempty"`
	Carrier      string     `json:"carrier,omitempty" url:"carrier,omitempty"`
	PaginatedCollection
}

// ListTrackersUpdatedOptions specifies options for the list trackers updated
// API.
type ListTrackersUpdatedOptions struct {
	Page                 int       `json:"page,omitempty" url:"page,omitempty"`
	PageSize             int       `json:"page_size,omitempty" url:"page_size,omitempty"`
	StatusStart          *DateTime `json:"status_start,omitempty" url:"status_start,omitempty"`
	StatusEnd            *DateTime `json:"status_end,omitempty" url:"status_end,omitempty"`
	TrackingDetailsStart *DateTime `json:"tracking_details_start,omitempty" url:"tracking_details_start,omitempty"`
	TrackingDetailsEnd   *DateTime `json:"tracking_details_end,omitempty" url:"tracking_details_end,omitempty"`
}

func (c *CreateTrackerOptions) toMap() map[string]interface{} {
	trackerParams := make(map[string]interface{})
	if c.TrackingCode != "" {
		trackerParams["tracking_code"] = c.TrackingCode
	}
	if c.Carrier != "" {
		trackerParams["carrier"] = c.Carrier
	}
	if c.Amount != "" {
		trackerParams["amount"] = c.Amount
	}
	optionsParams := make(map[string]interface{})
	if c.CarrierAccount != "" {
		optionsParams["carrier_account"] = c.CarrierAccount
	}
	if c.IsReturn {
		optionsParams["is_return"] = true
	}
	m := make(map[string]interface{})
	if len(trackerParams) != 0 {
		m["tracker"] = trackerParams
	}
	if len(optionsParams) != 0 {
		m["options"] = optionsParams
	}
	if len(m) == 0 {
		return nil
	}
	return m
}

// CreateTracker creates a new Tracker object with the provided options.
// Providing a carrier is optional, but helps to avoid ambiguity in detecting
// the carrier based on the tracking code format.
func (c *Client) CreateTracker(opts *CreateTrackerOptions) (out *Tracker, err error) {
	return c.CreateTrackerWithContext(context.Background(), opts)
}

// CreateTrackerWithContext performs the same operation as CreateTracker, but
// allows specifying a context that can interrupt the request.
func (c *Client) CreateTrackerWithContext(ctx context.Context, opts *CreateTrackerOptions) (out *Tracker, err error) {
	err = c.do(ctx, http.MethodPost, "trackers", opts.toMap(), &out)
	return
}

// ListTrackers provides a paginated result of Tracker objects.
func (c *Client) ListTrackers(opts *ListTrackersOptions) (out *ListTrackersResult, err error) {
	return c.ListTrackersWithContext(context.Background(), opts)
}

// ListTrackersWithContext performs the same operation as ListTrackers, but
// allows specifying a context that can interrupt the request.
func (c *Client) ListTrackersWithContext(ctx context.Context, opts *ListTrackersOptions) (out *ListTrackersResult, err error) {
	err = c.do(ctx, http.MethodGet, "trackers", opts, &out)
	// Store the original query parameters for reuse when getting the next page
	out.TrackingCode = opts.TrackingCode
	out.Carrier = opts.Carrier
	return
}

// GetNextTrackerPage returns the next page of trackers
func (c *Client) GetNextTrackerPage(collection *ListTrackersResult) (out *ListTrackersResult, err error) {
	return c.GetNextTrackerPageWithContext(context.Background(), collection)
}

// GetNextTrackerPageWithPageSize returns the next page of trackers with a specific page size
func (c *Client) GetNextTrackerPageWithPageSize(collection *ListTrackersResult, pageSize int) (out *ListTrackersResult, err error) {
	return c.GetNextTrackerPageWithPageSizeWithContext(context.Background(), collection, pageSize)
}

// GetNextTrackerPageWithContext performs the same operation as GetNextTrackerPage, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetNextTrackerPageWithContext(ctx context.Context, collection *ListTrackersResult) (out *ListTrackersResult, err error) {
	return c.GetNextTrackerPageWithPageSizeWithContext(ctx, collection, 0)
}

// GetNextTrackerPageWithPageSizeWithContext performs the same operation as GetNextTrackerPageWithPageSize, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetNextTrackerPageWithPageSizeWithContext(ctx context.Context, collection *ListTrackersResult, pageSize int) (out *ListTrackersResult, err error) {
	if len(collection.Trackers) == 0 {
		err = newEndOfPaginationError()
		return
	}
	lastID := collection.Trackers[len(collection.Trackers)-1].ID
	params, err := nextPageParameters(collection.HasMore, lastID, pageSize)
	if err != nil {
		return
	}
	trackerParams := &ListTrackersOptions{
		BeforeID:     params.BeforeID,
		Carrier:      collection.Carrier,
		TrackingCode: collection.TrackingCode,
	}
	if pageSize > 0 {
		trackerParams.PageSize = pageSize
	}
	return c.ListTrackersWithContext(ctx, trackerParams)
}

// GetTracker retrieves a Tracker object by ID.
func (c *Client) GetTracker(trackerID string) (out *Tracker, err error) {
	return c.GetTrackerWithContext(context.Background(), trackerID)
}

// GetTrackerWithContext performs the same operation as GetTracker, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetTrackerWithContext(ctx context.Context, trackerID string) (out *Tracker, err error) {
	err = c.do(ctx, http.MethodGet, "trackers/"+trackerID, nil, &out)
	return
}

// RetrieveTrackerBatch retrieves a batch of Tracker objects.
func (c *Client) RetrieveTrackerBatch(opts *ListTrackersOptions) (out *ListTrackersResult, err error) {
	return c.RetrieveTrackerBatchWithContext(context.Background(), opts)
}

// RetrieveTrackerBatchWithContext performs the same operation as ListTrackers, but
// allows specifying a context that can interrupt the request.
func (c *Client) RetrieveTrackerBatchWithContext(ctx context.Context, opts *ListTrackersOptions) (out *ListTrackersResult, err error) {
	err = c.do(ctx, http.MethodPost, "trackers/batch", opts, &out)
	return
}
