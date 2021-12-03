package easypost

import (
	"context"
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
	err = c.post(context.Background(), "pickups", &createPickupRequest{Pickup: in}, &out)
	return
}

// CreatePickupWithContext performs the same operation as CreatePickup, but
// allows specifying a context that can interrupt the request.
func (c *Client) CreatePickupWithContext(ctx context.Context, in *Pickup) (out *Pickup, err error) {
	err = c.post(ctx, "pickups", &createPickupRequest{Pickup: in}, &out)
	return
}

// GetPickup retrieves an existing Pickup object by ID.
func (c *Client) GetPickup(pickupID string) (out *Pickup, err error) {
	err = c.get(context.Background(), "pickups/"+pickupID, &out)
	return
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
	vals := url.Values{
		"carrier": []string{rate.Carrier}, "service": []string{rate.Service},
	}
	err = c.post(context.Background(), "pickups/"+pickupID+"/buy", vals, &out)
	return
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
	err = c.post(context.Background(), "pickups/"+pickupID+"/cancel", nil, &out)
	return
}

// CancelPickupWithContext performs the same operation as CancelPickup, but
// allows specifying a context that can interrupt the request.
func (c *Client) CancelPickupWithContext(ctx context.Context, pickupID string) (out *Pickup, err error) {
	err = c.post(ctx, "pickups/"+pickupID+"/cancel", nil, &out)
	return
}
