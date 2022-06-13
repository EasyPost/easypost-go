package easypost

import (
	"context"
	"time"
)

// A Rate contains information on shipping cost and delivery time.
type Rate struct {
	ID                     string     `json:"id,omitempty"`
	Object                 string     `json:"object,omitempty"`
	Mode                   string     `json:"mode,omitempty"`
	CreatedAt              *time.Time `json:"created_at,omitempty"`
	UpdatedAt              *time.Time `json:"updated_at,omitempty"`
	Service                string     `json:"service,omitempty"`
	Carrier                string     `json:"carrier,omitempty"`
	CarrierAccountID       string     `json:"carrier_account_id,omitempty"`
	ShipmentID             string     `json:"shipment_id,omitempty"`
	Rate                   string     `json:"rate,omitempty"`
	Currency               string     `json:"currency,omitempty"`
	RetailRate             string     `json:"retail_rate,omitempty"`
	RetailCurrency         string     `json:"retail_currency,omitempty"`
	ListRate               string     `json:"list_rate,omitempty"`
	ListCurrency           string     `json:"list_currency,omitempty"`
	DeliveryDays           int        `json:"delivery_days,omitempty"`
	DeliveryDate           *time.Time `json:"delivery_date,omitempty"`
	DeliveryDateGuaranteed bool       `json:"delivery_date_guaranteed,omitempty"`
	EstDeliveryDays        int        `json:"est_delivery_days,omitempty"`
	BillingType            string     `json:"billing_type,omitempty"`
}

// A SmartRate contains information on shipping cost and delivery time in addition to time-in-transit details.
type SmartRate struct {
	ID                     string         `json:"id,omitempty"`
	Object                 string         `json:"object,omitempty"`
	Mode                   string         `json:"mode,omitempty"`
	CreatedAt              *time.Time     `json:"created_at,omitempty"`
	UpdatedAt              *time.Time     `json:"updated_at,omitempty"`
	Service                string         `json:"service,omitempty"`
	Carrier                string         `json:"carrier,omitempty"`
	CarrierAccountID       string         `json:"carrier_account_id,omitempty"`
	ShipmentID             string         `json:"shipment_id,omitempty"`
	Rate                   float64        `json:"rate,omitempty"`
	Currency               string         `json:"currency,omitempty"`
	RetailRate             float64        `json:"retail_rate,omitempty"`
	RetailCurrency         string         `json:"retail_currency,omitempty"`
	ListRate               float64        `json:"list_rate,omitempty"`
	ListCurrency           string         `json:"list_currency,omitempty"`
	DeliveryDays           int            `json:"delivery_days,omitempty"`
	DeliveryDate           *time.Time     `json:"delivery_date,omitempty"`
	DeliveryDateGuaranteed bool           `json:"delivery_date_guaranteed,omitempty"`
	EstDeliveryDays        int            `json:"est_delivery_days,omitempty"`
	TimeInTransit          *TimeInTransit `json:"time_in_transit,omitempty"`
}

// TimeInTransit provides details on the probability your package will arrive within a certain number of days
type TimeInTransit struct {
	Percentile50 int `json:"percentile_50,omitempty"`
	Percentile75 int `json:"percentile_75,omitempty"`
	Percentile85 int `json:"percentile_85,omitempty"`
	Percentile90 int `json:"percentile_90,omitempty"`
	Percentile95 int `json:"percentile_95,omitempty"`
	Percentile97 int `json:"percentile_97,omitempty"`
	Percentile99 int `json:"percentile_99,omitempty"`
}

// GetRate retrieves a previously-created rate by its ID.
func (c *Client) GetRate(rateID string) (out *Rate, err error) {
	err = c.get(context.Background(), "rates/"+rateID, &out)
	return
}

// GetRateWithContext performs the same operation as GetRate, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetRateWithContext(ctx context.Context, rateID string) (out *Rate, err error) {
	err = c.get(ctx, "rates/"+rateID, &out)
	return
}
