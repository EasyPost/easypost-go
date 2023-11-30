package easypost

import (
	"context"
)

// A Rate contains information on shipping cost and delivery time.
type Rate struct {
	ID                     string    `json:"id,omitempty"`
	Object                 string    `json:"object,omitempty"`
	Mode                   string    `json:"mode,omitempty"`
	CreatedAt              *DateTime `json:"created_at,omitempty"`
	UpdatedAt              *DateTime `json:"updated_at,omitempty"`
	Service                string    `json:"service,omitempty"`
	Carrier                string    `json:"carrier,omitempty"`
	CarrierAccountID       string    `json:"carrier_account_id,omitempty"`
	ShipmentID             string    `json:"shipment_id,omitempty"`
	Rate                   string    `json:"rate,omitempty"`
	Currency               string    `json:"currency,omitempty"`
	RetailRate             string    `json:"retail_rate,omitempty"`
	RetailCurrency         string    `json:"retail_currency,omitempty"`
	ListRate               string    `json:"list_rate,omitempty"`
	ListCurrency           string    `json:"list_currency,omitempty"`
	DeliveryDays           int       `json:"delivery_days,omitempty"`
	DeliveryDate           *DateTime `json:"delivery_date,omitempty"`
	DeliveryDateGuaranteed bool      `json:"delivery_date_guaranteed,omitempty"`
	EstDeliveryDays        int       `json:"est_delivery_days,omitempty"`
	BillingType            string    `json:"billing_type,omitempty"`
}

// A SmartRate contains information on shipping cost and delivery time in addition to time-in-transit details.
type SmartRate struct {
	ID                     string         `json:"id,omitempty"`
	Object                 string         `json:"object,omitempty"`
	Mode                   string         `json:"mode,omitempty"`
	CreatedAt              *DateTime      `json:"created_at,omitempty"`
	UpdatedAt              *DateTime      `json:"updated_at,omitempty"`
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
	DeliveryDate           *DateTime      `json:"delivery_date,omitempty"`
	DeliveryDateGuaranteed bool           `json:"delivery_date_guaranteed,omitempty"`
	EstDeliveryDays        int            `json:"est_delivery_days,omitempty"`
	TimeInTransit          *TimeInTransit `json:"time_in_transit,omitempty"`
	BillingType            string         `json:"billing_type,omitempty"`
}

// A StatelessRate contains information on shipping cost and delivery time, but does not have an ID (is ephemeral).
type StatelessRate struct {
	BillingType            string    `json:"billing_type,omitempty"`
	Carrier                string    `json:"carrier,omitempty"`
	CarrierAccountID       string    `json:"carrier_account_id,omitempty"`
	Currency               string    `json:"currency,omitempty"`
	DeliveryDate           *DateTime `json:"delivery_date,omitempty"`
	DeliveryDateGuaranteed bool      `json:"delivery_date_guaranteed,omitempty"`
	DeliveryDays           int       `json:"delivery_days,omitempty"`
	EstDeliveryDays        int       `json:"est_delivery_days,omitempty"`
	ListCurrency           string    `json:"list_currency,omitempty"`
	ListRate               string    `json:"list_rate,omitempty"`
	Mode                   string    `json:"mode,omitempty"`
	Object                 string    `json:"object,omitempty"`
	Rate                   string    `json:"rate,omitempty"`
	RetailCurrency         string    `json:"retail_currency,omitempty"`
	RetailRate             string    `json:"retail_rate,omitempty"`
	Service                string    `json:"service,omitempty"`
	ShipmentID             string    `json:"shipment_id,omitempty"`
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

type getStatelessRatesResponse struct {
	Rates *[]*StatelessRate `json:"rates,omitempty"`
}

// GetRate retrieves a previously-created rate by its ID.
func (c *Client) GetRate(rateID string) (out *Rate, err error) {
	return c.GetRateWithContext(context.Background(), rateID)
}

// GetRateWithContext performs the same operation as GetRate, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetRateWithContext(ctx context.Context, rateID string) (out *Rate, err error) {
	err = c.get(ctx, "rates/"+rateID, &out)
	return
}

// BetaGetStatelessRates fetches a list of stateless rates.
func (c *Client) BetaGetStatelessRates(in *Shipment) (out []*StatelessRate, err error) {
	return c.BetaGetStatelessRatesWithContext(context.Background(), in)
}

// BetaGetStatelessRatesWithContext performs the same operation as BetaGetStatelessRates,
// but allows specifying a context that can interrupt the request.
func (c *Client) BetaGetStatelessRatesWithContext(ctx context.Context, in *Shipment) (out []*StatelessRate, err error) {
	req := &createShipmentRequest{Shipment: in}
	res := &getStatelessRatesResponse{Rates: &out}
	err = c.post(ctx, "/beta/rates", &req, &res)
	return
}

// LowestStatelessRate gets the lowest stateless rate from a list of stateless rates
func (c *Client) LowestStatelessRate(rates []*StatelessRate) (out StatelessRate, err error) {
	return c.LowestStatelessRateWithCarrier(rates, nil)
}

// LowestStatelessRateWithCarrier performs the same operation as LowestStatelessRate,
// but allows specifying a list of carriers for the lowest rate
func (c *Client) LowestStatelessRateWithCarrier(rates []*StatelessRate, carriers []string) (out StatelessRate, err error) {
	return c.LowestStatelessRateWithCarrierAndService(rates, carriers, nil)
}

// LowestStatelessRateWithCarrierAndService performs the same operation as LowestStatelessRate,
// but allows specifying a list of carriers and service for the lowest rate
func (c *Client) LowestStatelessRateWithCarrierAndService(rates []*StatelessRate, carriers []string, services []string) (out StatelessRate, err error) {
	return c.lowestStatelessRate(rates, carriers, services)
}
