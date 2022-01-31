package easypost

import (
	"context"
	"net/http"
	"time"
)

// TrackingLocation provides additional information about the location of a
// tracking event.
type TrackingLocation struct {
	Object  string `json:"object,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Country string `json:"country,omitempty"`
	Zip     string `json:"zip,omitempty"`
}

// TrackingDetail provides information about a tracking event.
type TrackingDetail struct {
	Object           string            `json:"object,omitempty"`
	Message          string            `json:"message,omitempty"`
	Description      string            `json:"description,omitempty"`
	Status           string            `json:"status,omitempty"`
	DateTime         string            `json:"datetime,omitempty"`
	Source           string            `json:"source,omitempty"`
	CarrierCode      string            `json:"carrier_code,omitempty"`
	TrackingLocation *TrackingLocation `json:"tracking_location,omitempty"`
}

// TrackingCarrierDetail provides additional tracking information from the
// carrier, when available.
type TrackingCarrierDetail struct {
	Object                      string            `json:"object,omitempty"`
	Service                     string            `json:"service,omitempty"`
	ContainerType               string            `json:"container_type,omitempty"`
	EstDeliveryDateLocal        string            `json:"est_delivery_date_local,omitempty"`
	EstDeliveryTimeLocal        string            `json:"est_delivery_time_local,omitempty"`
	OriginLocation              string            `json:"origin_locaion,omitempty"`
	OriginTrackingLocation      *TrackingLocation `json:"origin_tracking_location,omitempty"`
	DestinationLocation         string            `json:"destination_location,omitempty"`
	DestinationTrackingLocation *TrackingLocation `json:"destination_tracking_location,omitempty"`
	GuaranteedDeliveryDate      *time.Time        `json:"guaranteed_delivery_date,omitempty"`
	AlternateIdentifier         string            `json:"alternate_identifier,omitempty"`
	InitialDeliveryAttempt      *time.Time        `json:"initial_delivery_attempt,omitempty"`
}

// A Tracker object contains all of the tracking information for a package.
type Tracker struct {
	ID              string                 `json:"id,omitempty"`
	Object          string                 `json:"object,omitempty"`
	Mode            string                 `json:"mode,omitempty"`
	CreatedAt       *time.Time             `json:"created_at,omitempty"`
	UpdatedAt       *time.Time             `json:"updated_at,omitempty"`
	TrackingCode    string                 `json:"tracking_code,omitempty"`
	Status          string                 `json:"status,omitempty"`
	SignedBy        string                 `json:"signed_by,omitempty"`
	Weight          float64                `json:"weight,omitempty"`
	EstDeliveryDate *time.Time             `json:"est_delivery_date,omitempty"`
	ShipmentID      string                 `json:"shipment_id,omitempty"`
	Carrier         string                 `json:"carrier,omitempty"`
	TrackingDetails []*TrackingDetail      `json:"tracking_details,omitempty"`
	CarrierDetail   *TrackingCarrierDetail `json:"carrier_detail,omitempty"`
	PublicURL       string                 `json:"public_url,omitempty"`
	Fees            []*Fee                 `json:"fees,omitempty"`
	Finalized       bool                   `json:"finalized,omitempty"`
	IsReturn        bool                   `json:"is_return,omitempty"`
}

// CreateTrackerOptions specifies options for creating a new tracker.
type CreateTrackerOptions struct {
	TrackingCode    string
	Carrier         string
	Amount          string
	CarrierAccount  string
	IsReturn        bool
	FullTestTracker bool
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
	if c.FullTestTracker {
		optionsParams["full_test_tracker"] = true
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
	err = c.post(nil, "trackers", opts.toMap(), &out)
	return
}

// CreateTrackerWithContext performs the same operation as CreateTracker, but
// allows specifying a context that can interrupt the request.
func (c *Client) CreateTrackerWithContext(ctx context.Context, opts *CreateTrackerOptions) (out *Tracker, err error) {
	err = c.post(ctx, "trackers", opts.toMap(), &out)
	return
}

// CreateTrackerList asynchronously creates multiple trackers.
// Input a map of maps that contains multiple tracking codes
func (c *Client) CreateTrackerList(param map[string]interface{}) (bool, error) {
	// The data structure must look like the following when calling the API:
	// {
	//     "trackers": {
	//         "0": { "tracking_code": "EZ1000000001", "carrier": "USPS" },
	//         "1": { "tracking_code": "EZ1000000002", "carrier": "USPS" }
	//     }
	// }
	// The keys inside of the 'trackers' map (0, 1 in the example) get discarded
	// by the API endpoint, so are not important.
	req := map[string]interface{}{"trackers": param}
	//this endpoint does not return a response so we return true here
	return true, c.post(nil, "trackers/create_list", req, nil)
}

// CreateTrackerListWithContext performs the same operation as
// CreateTrackerList, but allows specifying a context that can interrupt the
// request.
func (c *Client) CreateTrackerListWithContext(ctx context.Context, param map[string]interface{}) (bool, error) {
	req := map[string]interface{}{"trackers": param}
	//this endpoint does not return a response so we return true here
	return true, c.post(ctx, "trackers/create_list", req, nil)
}

// ListTrackersOptions is used to specify query parameters for listing Tracker
// objects.
type ListTrackersOptions struct {
	BeforeID      string     `url:"before_id,omitempty"`
	AfterID       string     `url:"after_id,omitempty"`
	StartDateTime *time.Time `url:"start_datetime,omitempty"`
	EndDateTime   *time.Time `url:"end_datetime,omitempty"`
	PageSize      int        `url:"page_size,omitempty"`
	TrackingCodes []string   `url:"tracking_codes,omitempty"`
	Carrier       string     `url:"carrier,omitempty"`
}

// ListTrackersResult holds the results from the list trackers API.
type ListTrackersResult struct {
	Trackers []*Tracker `json:"trackers,omitempty"`
	// HasMore indicates if there are more responses to be fetched. If True,
	// additional responses can be fetched by updating the ListTrackersOptions
	// parameter's AfterID field with the ID of the last item in this object's
	// Trackers field.
	HasMore bool `json:"has_more,omitempty"`
}

// ListTrackers provides a paginated result of Tracker objects.
func (c *Client) ListTrackers(opts *ListTrackersOptions) (out *ListTrackersResult, err error) {
	return c.ListTrackersWithContext(nil, opts)
}

// ListTrackersWithContext performs the same operation as ListTrackers, but
// allows specifying a context that can interrupt the request.
func (c *Client) ListTrackersWithContext(ctx context.Context, opts *ListTrackersOptions) (out *ListTrackersResult, err error) {
	err = c.do(ctx, http.MethodGet, "trackers", c.convertOptsToURLValues(opts), &out)
	return
}

// ListTrackersUpdatedOptions specifies options for the list trackers updated
// API.
type ListTrackersUpdatedOptions struct {
	Page                 int        `json:"page,omitempty"`
	PageSize             int        `json:"page_size,omitempty"`
	StatusStart          *time.Time `json:"status_start,omitempty"`
	StatusEnd            *time.Time `json:"status_end,omitempty"`
	TrackingDetailsStart *time.Time `json:"tracking_details_start,omitempty"`
	TrackingDetailsEnd   *time.Time `json:"tracking_details_end,omitempty"`
}

// GetTracker retrieves a Tracker object by ID.
func (c *Client) GetTracker(trackerID string) (out *Tracker, err error) {
	err = c.get(nil, "trackers/"+trackerID, &out)
	return
}

// GetTrackerWithContext performs the same operation as GetTracker, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetTrackerWithContext(ctx context.Context, trackerID string) (out *Tracker, err error) {
	err = c.get(ctx, "trackers/"+trackerID, &out)
	return
}
