package easypost

import (
	"context"
	"net/url"
)

// An Order object represents a collection of packages and can be used for
// multi-piece Shipments.
type Order struct {
	ID            string            `json:"id,omitempty"`
	Object        string            `json:"object,omitempty"`
	Reference     string            `json:"reference,omitempty"`
	Mode          string            `json:"mode,omitempty"`
	CreatedAt     *DateTime         `json:"created_at,omitempty"`
	UpdatedAt     *DateTime         `json:"updated_at,omitempty"`
	ToAddress     *Address          `json:"to_address,omitempty"`
	FromAddress   *Address          `json:"from_address,omitempty"`
	ReturnAddress *Address          `json:"return_address,omitempty"`
	BuyerAddress  *Address          `json:"buyer_address,omitempty"`
	Shipments     []*Shipment       `json:"shipments,omitempty"`
	Rates         []*Rate           `json:"rates,omitempty"`
	Messages      []*CarrierMessage `json:"messages,omitempty"`
	IsReturn      bool              `json:"is_return"`
	Service       string            `json:"service,omitempty"`
	CustomsInfo   *CustomsInfo      `json:"customs_info,omitempty"`
}

type createOrderRequest struct {
	Order struct {
		*Order
		CarrierAccounts []*CarrierAccount `json:"carrier_accounts,omitempty"`
	} `json:"order,omitempty"`
}

// CreateOrder creates a new order object. If the `accounts` parameter is given,
// the provided carrier accounts will be used to limit the returned rates to
// the given carrier(s).
//
//	c := easypost.New(MyEasyPostAPIKey)
//	out, err := c.CreateOrder(
//		&easypost.Order{
//			ToAddress:   &easypost.Address{ID: "adr_1001"},
//			FromAddress: &easypost.Address{Id: "adr_101"},
//			Shipments:   []*easypost.Shipment{
//				&easypost.Shipment{
//					Parcel: &easypost.Parcel{
//						PredefinedPackage: "FedExBox",
//						Weight:            10.2,
//					},
//				},
//				&easypost.Shipment{
//					Parcel: &easypost.Parcel{
//						PredefinedPackage: "FedExBox",
//						Weight:            17.5,
//					},
//				},
//			},
//		},
//		&easypost.CarrierAccount{ID: "ca_101"},
//		&easypost.CarrierAccount{ID: "ca_102"},
//	)
func (c *Client) CreateOrder(in *Order, accounts ...*CarrierAccount) (out *Order, err error) {
	return c.CreateOrderWithContext(context.Background(), in, accounts...)
}

// CreateOrderWithContext performs the same operation as CreateOrder, but allows
// specifying a context that can interrupt the request.
func (c *Client) CreateOrderWithContext(ctx context.Context, in *Order, accounts ...*CarrierAccount) (out *Order, err error) {
	var req createOrderRequest
	req.Order.Order, req.Order.CarrierAccounts = in, accounts
	err = c.post(ctx, "orders", &req, &out)
	return
}

// GetOrder retrieves an existing Order object by ID.
func (c *Client) GetOrder(orderID string) (out *Order, err error) {
	return c.GetOrderWithContext(context.Background(), orderID)
}

// GetOrderWithContext performs the same operation as GetOrder, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetOrderWithContext(ctx context.Context, orderID string) (out *Order, err error) {
	err = c.get(ctx, "orders/"+orderID, &out)
	return
}

// GetOrderRates refreshes rates for an Order.
func (c *Client) GetOrderRates(orderID string) (out *Order, err error) {
	return c.GetOrderRatesWithContext(context.Background(), orderID)
}

// GetOrderRatesWithContext performs the same operation as GetOrderRates, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetOrderRatesWithContext(ctx context.Context, orderID string) (out *Order, err error) {
	err = c.get(ctx, "orders/"+orderID+"/rates", &out)
	return
}

// BuyOrder purchases an order. This operation populates the TrackingCode and
// PostageLabel attributes of each Shipment.
//
//	c := easypost.New(MyEasyPostAPIKey)
//	out, err := c.Buy("order_1", "FedEx", "FEDEX_GROUND")
func (c *Client) BuyOrder(orderID, carrier, service string) (out *Order, err error) {
	return c.BuyOrderWithContext(context.Background(), orderID, carrier, service)
}

// BuyOrderWithContext performs the same operation as GBuyOrder, but allows
// specifying a context that can interrupt the request.
func (c *Client) BuyOrderWithContext(ctx context.Context, orderID, carrier, service string) (out *Order, err error) {
	vals := url.Values{
		"carrier": []string{carrier},
		"service": []string{service},
	}
	err = c.post(ctx, "orders/"+orderID+"/buy", vals, &out)
	return
}

// LowestOrderRate gets the lowest rate of an order
func (c *Client) LowestOrderRate(order *Order) (out Rate, err error) {
	return c.LowestOrderRateWithCarrier(order, nil)
}

// LowestOrderRateWithCarrier performs the same operation as LowestOrderRate,
// but allows specifying a list of carriers for the lowest rate
func (c *Client) LowestOrderRateWithCarrier(order *Order, carriers []string) (out Rate, err error) {
	return c.LowestOrderRateWithCarrierAndService(order, carriers, nil)
}

// LowestOrderRateWithCarrierAndService performs the same operation as LowestOrderRate,
// but allows specifying a list of carriers and service for the lowest rate
func (c *Client) LowestOrderRateWithCarrierAndService(order *Order, carriers []string, services []string) (out Rate, err error) {
	return c.lowestObjectRate(order.Rates, carriers, services)
}
