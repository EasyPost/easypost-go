package easypost

import (
	"context"
)

// DeliveryDateForZipPairEstimate is a single zip-pair-based delivery date estimate for a carrier-service level combination.
type DeliveryDateForZipPairEstimate struct {
	Carrier                   string                     `json:"carrier,omitempty"`
	Service                   string                     `json:"service,omitempty"`
	EasyPostTimeInTransitData *EasyPostTimeInTransitData `json:"easypost_time_in_transit_data,omitempty"`
}

// EstimateDeliveryDateForZipPairResult is the result of the EstimateDeliveryDateForZipPair method, containing the estimated delivery date of each carrier-service level combination and additional metadata.
type EstimateDeliveryDateForZipPairResult struct {
	CarriersWithoutEstimates []string                          `json:"carriers_without_tint_estimates,omitempty"`
	FromZip                  string                            `json:"from_zip,omitempty"`
	ToZip                    string                            `json:"to_zip,omitempty"`
	SaturdayDelivery         bool                              `json:"saturday_delivery,omitempty"`
	PlannedShipDate          *DateTime                         `json:"planned_ship_date,omitempty"`
	Results                  []*DeliveryDateForZipPairEstimate `json:"results,omitempty"`
}

// EstimateDeliveryDateForZipPairParams are used in the EstimateDeliveryDateForZipPair method.
type EstimateDeliveryDateForZipPairParams struct {
	FromZip          string   `json:"from_zip,omitempty"`
	ToZip            string   `json:"to_zip,omitempty"`
	Carriers         []string `json:"carriers,omitempty"`
	PlannedShipDate  string   `json:"planned_ship_date,omitempty"`
	SaturdayDelivery bool     `json:"saturday_delivery,omitempty"`
}

// TimeInTransitDetailsForShipDate contains the time-in-transit details and estimated delivery date for a specific ShipDateForZipPairRecommendation or RecommendShipDateForShipmentResult.
type TimeInTransitDetailsForShipDate struct {
	DesiredDeliveryDate         *DateTime      `json:"desired_delivery_date,omitempty"`
	EasyPostRecommendedShipDate *DateTime      `json:"ship_on_date,omitempty"`
	DeliveryDateConfidence      float64        `json:"delivery_date_confidence,omitempty"`
	EstimatedTransitDays        int            `json:"estimated_transit_days,omitempty"`
	DaysInTransit               *TimeInTransit `json:"days_in_transit,omitempty"`
}

// ShipDateForZipPairRecommendation is a single zip-pair-based ship date recommendation for a carrier-service level combination.
type ShipDateForZipPairRecommendation struct {
	Carrier                   string                           `json:"carrier,omitempty"`
	Service                   string                           `json:"service,omitempty"`
	EasyPostTimeInTransitData *TimeInTransitDetailsForShipDate `json:"easypost_time_in_transit_data,omitempty"`
}

// RecommendShipDateForZipPairResult is the result of the RecommendShipDateForZipPair method, containing the recommended ship date of each carrier-service level combination and additional metadata.
type RecommendShipDateForZipPairResult struct {
	CarriersWithoutEstimates []string                            `json:"carriers_without_tint_estimates,omitempty"`
	FromZip                  string                              `json:"from_zip,omitempty"`
	ToZip                    string                              `json:"to_zip,omitempty"`
	SaturdayDelivery         bool                                `json:"saturday_delivery,omitempty"`
	DesiredDeliveryDate      *DateTime                           `json:"desired_delivery_date,omitempty"`
	Results                  []*ShipDateForZipPairRecommendation `json:"results,omitempty"`
}

// RecommendShipDateForZipPairParams are used in the RecommendShipDateForZipPair method.
type RecommendShipDateForZipPairParams struct {
	FromZip             string   `json:"from_zip,omitempty"`
	ToZip               string   `json:"to_zip,omitempty"`
	Carriers            []string `json:"carriers,omitempty"`
	DesiredDeliveryDate string   `json:"desired_delivery_date,omitempty"`
	SaturdayDelivery    bool     `json:"saturday_delivery,omitempty"`
}

// EstimateDeliveryDateForZipPair retrieves the estimated delivery date of each carrier-service level combination via the Smart Deliver By API, based on a specific ship date and origin-destination postal code pair.
func (c *Client) EstimateDeliveryDateForZipPair(params *EstimateDeliveryDateForZipPairParams) (out *EstimateDeliveryDateForZipPairResult, err error) {
	return c.EstimateDeliveryDateForZipPairWithContext(context.Background(), params)
}

// EstimateDeliveryDateForZipPairWithContext performs the same operation as EstimateDeliveryDateForZipPair, but allows specifying a context that can interrupt the request.
func (c *Client) EstimateDeliveryDateForZipPairWithContext(ctx context.Context, params *EstimateDeliveryDateForZipPairParams) (out *EstimateDeliveryDateForZipPairResult, err error) {
	err = c.post(ctx, "smartrate/deliver_by", params, &out)
	return
}

// RecommendShipDateForZipPair retrieves the recommended ship date of each carrier-service level combination via the Smart Deliver On API, based on a specific desired delivery date and origin-destination postal code pair.
func (c *Client) RecommendShipDateForZipPair(params *RecommendShipDateForZipPairParams) (out *RecommendShipDateForZipPairResult, err error) {
	return c.RecommendShipDateForZipPairWithContext(context.Background(), params)
}

// RecommendShipDateForZipPairWithContext performs the same operation as RecommendShipDateForZipPair, but allows specifying a context that can interrupt the request.
func (c *Client) RecommendShipDateForZipPairWithContext(ctx context.Context, params *RecommendShipDateForZipPairParams) (out *RecommendShipDateForZipPairResult, err error) {
	err = c.post(ctx, "smartrate/deliver_on", params, &out)
	return
}
