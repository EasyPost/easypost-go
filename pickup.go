package easypost

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// PickupRate contains data about the cost of a pickup.
type PickupRate struct {
	ID        string     `json:"id,omitempty"`
	Object    string     `json:"object,omitempty"`
	Mode      string     `json:"mode,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	Service   string     `json:"service,omitempty"`
	Carrier   string     `json:"carrier,omitempty"`
	Rate      string     `json:"rate,omitempty"`
	Currency  string     `json:"currency,omitempty"`
	PickupID  string     `json:"pickup_id,omitempty"`
}

// A Pickup object represents a pickup from a carrier at a customer's residence
// or place of business.
type Pickup struct {
	ID               string            `json:"id,omitempty"`
	Object           string            `json:"object,omitempty"`
	Reference        string            `json:"reference,omitempty"`
	Mode             string            `json:"mode,omitempty"`
	CreatedAt        *time.Time        `json:"created_at,omitempty"`
	UpdatedAt        *time.Time        `json:"updated_at,omitempty"`
	Status           string            `json:"status,omitempty"`
	MinDatetime      *time.Time        `json:"min_datetime,omitempty"`
	MaxDatetime      *time.Time        `json:"max_datetime,omitempty"`
	IsAccountAddress bool              `json:"is_account_address,omitempty"`
	Instructions     string            `json:"instructions,omitempty"`
	Messages         []*CarrierMessage `json:"messages,omitempty"`
	Confirmation     string            `json:"confirmation,omitempty"`
	Shipment         *Shipment         `json:"shipment,omitempty"`
	Address          *Address          `json:"address,omitempty"`
	Batch            *Batch            `json:"batch,omitempty"`
	CarrierAccounts  []*CarrierAccount `json:"carrier_accounts,omitempty"`
	PickupRates      []*PickupRate     `json:"pickup_rates,omitempty"`
}

// ListPickupResult holds the results from the list Pickup API.
type ListPickupResult struct {
	Pickups []*Pickup `json:"pickups,omitempty"`
	// HasMore indicates if there are more responses to be fetched. If True,
	// additional responses can be fetched by updating the ListOptions
	// parameter's AfterID field with the ID of the last item in this object's
	// Pickups field.
	HasMore bool `json:"has_more,omitempty"`
}

type createPickupRequest struct {
	Pickup *Pickup `json:"pickup,omitempty"`
}

// CreatePickup creates a new Pickup object, and automatically fetches rates
// for the given time and location.
//	c := easypost.New(MyEasyPostAPIKey)
//	out, err := c.CreatePickup(
//		&easypost.Pickup{
//			Reference:        "my-first-pickup",
//			MinDatetime:      time.Date(2014, 10, 21, 0, 10, 0, 0, time.UTC),
//			MaxDatetime:      time.Date(2014, 10, 21, 15, 30, 0, 0, time.UTC),
//			Shipment:         &easypost.Shipment{ID: "shp_1"},
//			Address:          &easypost.Address{ID: "ad_1001"},
//			IsAccountAddress: false,
//			Instructions:     "Special pickup instructions",
//		},
//	)
func (c *Client) CreatePickup(in *Pickup) (out *Pickup, err error) {
	return c.CreatePickupWithContext(context.Background(), in)
}

// CreatePickupWithContext performs the same operation as CreatePickup, but
// allows specifying a context that can interrupt the request.
func (c *Client) CreatePickupWithContext(ctx context.Context, in *Pickup) (out *Pickup, err error) {
	err = c.post(ctx, "pickups", &createPickupRequest{Pickup: in}, &out)
	return
}

// GetPickup retrieves an existing Pickup object by ID.
func (c *Client) GetPickup(pickupID string) (out *Pickup, err error) {
	return c.GetPickupWithContext(context.Background(), pickupID)
}

// GetPickupWithContext performs the same operation as GetPickup, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetPickupWithContext(ctx context.Context, pickupID string) (out *Pickup, err error) {
	err = c.get(ctx, "pickups/"+pickupID, &out)
	return
}

// BuyPickup purchases and schedules a pickup.
//	c := easypost.New(MyEasyPostAPIKey)
//  rate := &PickupRate{Carrier: "UPS", Service: "Same-Day Pickup"}
//	out, err := c.BuyPickup("pck_1", rate)
func (c *Client) BuyPickup(pickupID string, rate *PickupRate) (out *Pickup, err error) {
	return c.BuyPickupWithContext(context.Background(), pickupID, rate)
}

// BuyPickupWithContext performs the same operation as BuyPickup, but allows
// specifying a context that can interrupt the request.
func (c *Client) BuyPickupWithContext(ctx context.Context, pickupID string, rate *PickupRate) (out *Pickup, err error) {
	vals := url.Values{
		"carrier": []string{rate.Carrier}, "service": []string{rate.Service},
	}
	err = c.post(ctx, "pickups/"+pickupID+"/buy", vals, &out)
	return
}

// CancelPickup cancels a scheduled pickup.
func (c *Client) CancelPickup(pickupID string) (out *Pickup, err error) {
	return c.CancelPickupWithContext(context.Background(), pickupID)
}

// CancelPickupWithContext performs the same operation as CancelPickup, but
// allows specifying a context that can interrupt the request.
func (c *Client) CancelPickupWithContext(ctx context.Context, pickupID string) (out *Pickup, err error) {
	err = c.post(ctx, "pickups/"+pickupID+"/cancel", nil, &out)
	return
}

// LowestPickupRate gets the lowest rate of a pickup
func (c *Client) LowestPickupRate(pickup *Pickup) (out PickupRate, err error) {
	return c.LowestPickupRateWithCarrier(pickup, nil)
}

// LowestPickupRateWithCarrier performs the same operation as LowestPickupRate,
// but allows specifying a list of carriers for the lowest rate
func (c *Client) LowestPickupRateWithCarrier(pickup *Pickup, carriers []string) (out PickupRate, err error) {
	return c.LowestPickupRateWithCarrierAndService(pickup, carriers, nil)
}

// LowestPickupRateWithCarrierAndService performs the same operation as LowestPickupRate,
// but allows specifying a list of carriers and service for the lowest rate
func (c *Client) LowestPickupRateWithCarrierAndService(pickup *Pickup, carriers []string, services []string) (out PickupRate, err error) {
	return c.lowestPickupRate(pickup.PickupRates, carriers, services)
}

// ListPickups provides a paginated result of Pickup objects.
func (c *Client) ListPickups(opts *ListOptions) (out *ListPickupResult, err error) {
	return c.ListPickupsWithContext(context.Background(), opts)
}

// ListPickupsWithContext performs the same operation as ListPickups, but
// allows specifying a context that can interrupt the request.
func (c *Client) ListPickupsWithContext(ctx context.Context, opts *ListOptions) (out *ListPickupResult, err error) {
	err = c.do(ctx, http.MethodGet, "pickups", c.convertOptsToURLValues(opts), &out)
	return
}
