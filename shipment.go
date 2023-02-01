package easypost

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// A Form represents a form associated with a Shipment.
type Form struct {
	ID                      string     `json:"id,omitempty"`
	Object                  string     `json:"object,omitempty"`
	Mode                    string     `json:"mode,omitempty"`
	CreatedAt               *time.Time `json:"created_at,omitempty"`
	UpdatedAt               *time.Time `json:"updated_at,omitempty"`
	FormType                string     `json:"form_type,omitempty"`
	FormURL                 string     `json:"form_url,omitempty"`
	SubmittedElectronically bool       `json:"submitted_electronically,omitempty"`
}

// PostageLabel provides details of a shipping label for a purchased shipment.
type PostageLabel struct {
	ID              string     `json:"id,omitempty"`
	Object          string     `json:"object,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
	IntegratedForm  string     `json:"integrated_form,omitempty"`
	LabelDate       *time.Time `json:"label_date,omitempty"`
	LabelEPL2URL    string     `json:"label_epl2_url,omitempty"`
	LabelFileType   string     `json:"label_file_type,omitempty"`
	LabelPDFURL     string     `json:"label_pdf_url,omitempty"`
	LabelResolution float64    `json:"label_resolution,omitempty"`
	LabelSize       string     `json:"label_size,omitempty"`
	LabelType       string     `json:"label_type,omitempty"`
	LabelURL        string     `json:"label_url,omitempty"`
	LabelZPLURL     string     `json:"label_zpl_url,omitempty"`
}

// A Shipment represents its namesake, and is made up of a "to" and "from"
// addresses, the Parcel being shipped, and any customs forms required for
// international deliveries.
type Shipment struct {
	ID                string            `json:"id,omitempty"`
	Object            string            `json:"object,omitempty"`
	Reference         string            `json:"reference,omitempty"`
	Mode              string            `json:"mode,omitempty"`
	CreatedAt         *time.Time        `json:"created_at,omitempty"`
	UpdatedAt         *time.Time        `json:"updated_at,omitempty"`
	ToAddress         *Address          `json:"to_address,omitempty"`
	FromAddress       *Address          `json:"from_address,omitempty"`
	ReturnAddress     *Address          `json:"return_address,omitempty"`
	BuyerAddress      *Address          `json:"buyer_address,omitempty"`
	Parcel            *Parcel           `json:"parcel,omitempty"`
	Carrier           string            `json:"carrier,omitempty"`
	Service           string            `json:"service,omitempty"`
	CarrierAccountIDs []string          `json:"carrier_accounts,omitempty"`
	CustomsInfo       *CustomsInfo      `json:"customs_info,omitempty"`
	ScanForm          *ScanForm         `json:"scan_form,omitempty"`
	Forms             []*Form           `json:"forms,omitempty"`
	Insurance         string            `json:"insurance,omitempty"`
	Rates             []*Rate           `json:"rates,omitempty"`
	SelectedRate      *Rate             `json:"selected_rate,omitempty"`
	PostageLabel      *PostageLabel     `json:"postage_label,omitempty"`
	Messages          []*CarrierMessage `json:"messages,omitempty"`
	Options           *ShipmentOptions  `json:"options,omitempty"`
	IsReturn          bool              `json:"is_return,omitempty"`
	TrackingCode      string            `json:"tracking_code,omitempty"`
	USPSZone          int               `json:"usps_zone,omitempty"`
	Status            string            `json:"status,omitempty"`
	Tracker           *Tracker          `json:"tracker,omitempty"`
	Fees              []*Fee            `json:"fees,omitempty"`
	RefundStatus      string            `json:"refund_status,omitempty"`
	BatchID           string            `json:"batch_id,omitempty"`
	BatchStatus       string            `json:"batch_status,omitempty"`
	BatchMessage      string            `json:"batch_message,omitempty"`
	TaxIdentifiers    []*TaxIdentifier  `json:"tax_identifiers,omitempty"`
}

// ListShipmentsOptions is used to specify query parameters for listing Shipment
// objects.
type ListShipmentsOptions struct {
	BeforeID        string     `url:"before_id,omitempty"`
	AfterID         string     `url:"after_id,omitempty"`
	StartDateTime   *time.Time `url:"start_datetime,omitempty"`
	EndDateTime     *time.Time `url:"end_datetime,omitempty"`
	PageSize        int        `url:"page_size,omitempty"`
	Purchased       *bool      `url:"purchased,omitempty"`
	IncludeChildren *bool      `url:"include_children,omitempty"`
}

// ListShipmentsResult holds the results from the list shipments API.
type ListShipmentsResult struct {
	Shipments []*Shipment `json:"shipments,omitempty"`
	// HasMore indicates if there are more responses to be fetched. If True,
	// additional responses can be fetched by updating the ListShipmentsOptions
	// parameter's AfterID field with the ID of the last item in this object's
	// Shipments field.
	HasMore bool `json:"has_more,omitempty"`
}

type buyShipmentRequest struct {
	Rate         *Rate  `json:"rate,omitempty"`
	Insurance    string `json:"insurance,omitempty"`
	CarbonOffset bool   `json:"carbon_offset,omitempty"`
	EndShipperID string `json:"end_shipper_id,omitempty"`
}

type createShipmentRequest struct {
	Shipment     *Shipment `json:"shipment,omitempty"`
	CarbonOffset bool      `json:"carbon_offset,omitempty"`
}

type getShipmentRatesRequest struct {
	CarbonOffset bool `json:"carbon_offset,omitempty"`
}

type getShipmentRatesResponse struct {
	Rates *[]*Rate `json:"rates,omitempty"`
}

type getShipmentStatelessRatesResponse struct {
	Rates *[]*StatelessRate `json:"rates,omitempty"`
}

type generateFormRequest struct {
	Form map[string]interface{} `json:"form,omitempty"`
}

// CreateShipment creates a new Shipment object. The ToAddress, FromAddress and
// Parcel attributes are required. These objects may be fully-specified to
// create new ones at the same time as creating the Shipment, or they can refer
// to existing objects via their ID attribute. Passing in one or more carrier
// accounts to CreateShipment limits the returned rates to the specified
// carriers.
//
//	c := easypost.New(MyEasyPostAPIKey)
//	out, err := c.CreateShipment(
//		&easypost.Shipment{
//			ToAddress: &easypost.Address{
//				Name:    "Dr. Steve Brule",
//				Street1: "179 N Harbor Dr",
//				City:    "Redondo Beach",
//				State:   "CA",
//				Zip:     "90277",
//				Country: "US",
//				Phone:   "8573875756",
//				Email:   "dr_steve_brule@gmail.com",
//			},
//			FromAddress: &easypost.Address{ID: "adr_101"},
//			Parcel: &easypost.Parcel{
//				Length: 20.2,
//				Width:  10.9,
//				Height: 5,
//				Weight: 65.9,
//			},
//			CustomsInfo: &easypost.CustomsInfo{ID: "cstinfo_1"},
//		},
//	)
func (c *Client) CreateShipment(in *Shipment) (out *Shipment, err error) {
	return c.CreateShipmentWithContext(context.Background(), in)
}

// CreateShipmentWithCarbonOffset performs the same operation as CreateShipment, but includes a carbon offset.
func (c *Client) CreateShipmentWithCarbonOffset(in *Shipment) (out *Shipment, err error) {
	return c.CreateShipmentWithCarbonOffsetWithContext(context.Background(), in)
}

// CreateShipmentWithContext performs the same operation as CreateShipment, but
// allows specifying a context that can interrupt the request.
func (c *Client) CreateShipmentWithContext(ctx context.Context, in *Shipment) (out *Shipment, err error) {
	req := &createShipmentRequest{Shipment: in}
	err = c.post(ctx, "shipments", &req, &out)
	return
}

// CreateShipmentWithCarbonOffsetWithContext performs the same operation as CreateShipmentWithCarbonOffset, but
// allows specifying a context that can interrupt the request.
func (c *Client) CreateShipmentWithCarbonOffsetWithContext(ctx context.Context, in *Shipment) (out *Shipment, err error) {
	req := &createShipmentRequest{Shipment: in, CarbonOffset: true}
	err = c.post(ctx, "shipments", &req, &out)
	return
}

// ListShipments provides a paginated result of Shipment objects.
func (c *Client) ListShipments(opts *ListShipmentsOptions) (out *ListShipmentsResult, err error) {
	return c.ListShipmentsWithContext(context.Background(), opts)
}

// ListShipmentsWithContext performs the same operation as ListShipments, but
// allows specifying a context that can interrupt the request.
func (c *Client) ListShipmentsWithContext(ctx context.Context, opts *ListShipmentsOptions) (out *ListShipmentsResult, err error) {
	err = c.do(ctx, http.MethodGet, "shipments", c.convertOptsToURLValues(opts), &out)
	return
}

// GetShipment retrieves a Shipment object by ID.
func (c *Client) GetShipment(shipmentID string) (out *Shipment, err error) {
	return c.GetShipmentWithContext(context.Background(), shipmentID)
}

// GetShipmentWithContext performs the same operation as GetShipment, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetShipmentWithContext(ctx context.Context, shipmentID string) (out *Shipment, err error) {
	err = c.get(ctx, "shipments/"+shipmentID, &out)
	return
}

func (c *Client) buyShipment(ctx context.Context, shipmentID string, in *buyShipmentRequest) (out *Shipment, err error) {
	err = c.post(ctx, "shipments/"+shipmentID+"/buy", &in, &out)
	return
}

// BuyShipment purchases a shipment. If successful, the returned Shipment will
// have the PostageLabel attribute populated.
//
//	c := easypost.New(MyEasyPostAPIKey)
//	out, err := c.Buy("shp_100", &easypost.Rate{ID: "rate_1001"}, "249.99")
func (c *Client) BuyShipment(shipmentID string, rate *Rate, insurance string) (out *Shipment, err error) {
	return c.BuyShipmentWithContext(context.Background(), shipmentID, rate, insurance)
}

// BuyShipmentWithContext performs the same operation as BuyShipment, but allows
// specifying a context that can interrupt the request.
func (c *Client) BuyShipmentWithContext(ctx context.Context, shipmentID string, rate *Rate, insurance string) (out *Shipment, err error) {
	req := &buyShipmentRequest{Rate: rate, Insurance: insurance}
	return c.buyShipment(ctx, shipmentID, req)
}

// BuyShipmentWithCarbonOffset performs the same operation as BuyShipment, but includes a carbon offset.
func (c *Client) BuyShipmentWithCarbonOffset(shipmentID string, rate *Rate, insurance string) (out *Shipment, err error) {
	return c.BuyShipmentWithCarbonOffsetWithContext(context.Background(), shipmentID, rate, insurance)
}

// BuyShipmentWithCarbonOffsetWithContext performs the same operation as BuyShipmentWithCarbonOffset, but allows
// specifying a context that can interrupt the request.
func (c *Client) BuyShipmentWithCarbonOffsetWithContext(ctx context.Context, shipmentID string, rate *Rate, insurance string) (out *Shipment, err error) {
	req := &buyShipmentRequest{Rate: rate, Insurance: insurance, CarbonOffset: true}
	return c.buyShipment(ctx, shipmentID, req)
}

// BuyShipmentWithEndShipper performs the same operation as BuyShipment, but includes an EndShipper ID.
func (c *Client) BuyShipmentWithEndShipper(shipmentID string, rate *Rate, insurance string, endShipperID string) (out *Shipment, err error) {
	return c.BuyShipmentWithEndShipperWithContext(context.Background(), shipmentID, rate, insurance, endShipperID)
}

// BuyShipmentWithEndShipperWithContext performs the same operation as BuyShipmentWithEndShipper, but allows
// specifying a context that can interrupt the request.
func (c *Client) BuyShipmentWithEndShipperWithContext(ctx context.Context, shipmentID string, rate *Rate, insurance string, endShipperID string) (out *Shipment, err error) {
	req := &buyShipmentRequest{Rate: rate, Insurance: insurance, EndShipperID: endShipperID}
	return c.buyShipment(ctx, shipmentID, req)
}

// BuyShipmentWithCarbonOffsetAndEndShipper performs the same operation as BuyShipment, but includes a carbon offset and an EndShipper ID.
func (c *Client) BuyShipmentWithCarbonOffsetAndEndShipper(shipmentID string, rate *Rate, insurance string, endShipperID string) (out *Shipment, err error) {
	return c.BuyShipmentWithCarbonOffsetAndEndShipperWithContext(context.Background(), shipmentID, rate, insurance, endShipperID)
}

// BuyShipmentWithCarbonOffsetAndEndShipperWithContext performs the same operation as BuyShipmentWithCarbonOffsetAndEndShipper, but allows
// specifying a context that can interrupt the request.
func (c *Client) BuyShipmentWithCarbonOffsetAndEndShipperWithContext(ctx context.Context, shipmentID string, rate *Rate, insurance string, endShipperID string) (out *Shipment, err error) {
	req := &buyShipmentRequest{Rate: rate, Insurance: insurance, CarbonOffset: true, EndShipperID: endShipperID}
	return c.buyShipment(ctx, shipmentID, req)
}

// GetShipmentLabel enables retrieving the label for a shipment in a different
// format. The PostageLabel field in the returned Shipment object will reflect
// the new format.
func (c *Client) GetShipmentLabel(shipmentID, format string) (out *Shipment, err error) {
	return c.GetShipmentLabelWithContext(context.Background(), shipmentID, format)
}

// GetShipmentLabelWithContext performs the same operation as GetShipmentLabel,
// but allows specifying a context that can interrupt the request.
func (c *Client) GetShipmentLabelWithContext(ctx context.Context, shipmentID, format string) (out *Shipment, err error) {
	vals := url.Values{"file_format": []string{format}}
	err = c.do(ctx, http.MethodGet, "shipments/"+shipmentID+"/label", vals, &out)
	return
}

// GetShipmentSmartrates fetches the available smartrates for a shipment.
func (c *Client) GetShipmentSmartrates(shipmentID string) (out []*SmartRate, err error) {
	return c.GetShipmentSmartratesWithContext(context.Background(), shipmentID)
}

// GetShipmentSmartratesWithContext performs the same operation as GetShipmentRates,
// but allows specifying a context that can interrupt the request.
func (c *Client) GetShipmentSmartratesWithContext(ctx context.Context, shipmentID string) (out []*SmartRate, err error) {
	res := struct {
		Smartrates *[]*SmartRate `json:"result,omitempty"`
	}{Smartrates: &out}
	err = c.get(ctx, "shipments/"+shipmentID+"/smartrate", &res)
	return
}

// InsureShipment purchases insurance for the shipment. Insurance should be
// purchased after purchasing the shipment, but before it has been processed by
// the carrier. On success, the purchased insurance will be reflected in the
// returned Shipment object's Insurance field.
func (c *Client) InsureShipment(shipmentID, amount string) (out *Shipment, err error) {
	return c.InsureShipmentWithContext(context.Background(), shipmentID, amount)
}

// InsureShipmentWithContext performs the same operation as InsureShipment, but
// allows specifying a context that can interrupt the request.
func (c *Client) InsureShipmentWithContext(ctx context.Context, shipmentID, amount string) (out *Shipment, err error) {
	vals := url.Values{"amount": []string{amount}}
	err = c.post(ctx, "shipments/"+shipmentID+"/insure", vals, &out)
	return
}

// RefundShipment requests a refund from the carrier.
func (c *Client) RefundShipment(shipmentID string) (out *Shipment, err error) {
	return c.RefundShipmentWithContext(context.Background(), shipmentID)
}

// RefundShipmentWithContext performs the same operation as RefundShipment, but
// allows specifying a context that can interrupt the request.
func (c *Client) RefundShipmentWithContext(ctx context.Context, shipmentID string) (out *Shipment, err error) {
	err = c.post(ctx, "shipments/"+shipmentID+"/refund", nil, &out)
	return
}

// RerateShipment fetches the available rates for a shipment with the current rates.
func (c *Client) RerateShipment(shipmentID string) (out []*Rate, err error) {
	return c.RerateShipmentWithContext(context.Background(), shipmentID)
}

// RerateShipmentWithContext performs the same operation as RerateShipment,
// but allows specifying a context that can interrupt the request.
func (c *Client) RerateShipmentWithContext(ctx context.Context, shipmentID string) (out []*Rate, err error) {
	req := &getShipmentRatesRequest{CarbonOffset: false}
	res := &getShipmentRatesResponse{Rates: &out}
	err = c.post(ctx, "shipments/"+shipmentID+"/rerate", &req, &res)
	return
}

// RerateShipmentWithCarbonOffset performs the same operation as RerateShipment, but includes a carbon offset.
func (c *Client) RerateShipmentWithCarbonOffset(shipmentID string) (out []*Rate, err error) {
	return c.RerateShipmentWithCarbonOffsetWithContext(context.Background(), shipmentID)
}

// RerateShipmentWithCarbonOffsetWithContext performs the same operation as RerateShipmentWithCarbonOffset, but allows
// specifying a context that can interrupt the request.
func (c *Client) RerateShipmentWithCarbonOffsetWithContext(ctx context.Context, shipmentID string) (out []*Rate, err error) {
	req := &getShipmentRatesRequest{CarbonOffset: true}
	res := &getShipmentRatesResponse{Rates: &out}
	err = c.post(ctx, "shipments/"+shipmentID+"/rerate", &req, &res)
	return
}

// BetaGetStatelessRates fetches a list of stateless rates for a proposed shipment, without creating a shipment object.
func (c *Client) BetaGetStatelessRates(in *Shipment) (out []*StatelessRate, err error) {
	return c.BetaGetStatelessRatesWithContext(context.Background(), in)
}

// BetaGetStatelessRatesWithContext performs the same operation as BetaGetStatelessRates,
// but allows specifying a context that can interrupt the request.
func (c *Client) BetaGetStatelessRatesWithContext(ctx context.Context, in *Shipment) (out []*StatelessRate, err error) {
	req := &createShipmentRequest{Shipment: in}
	res := &getShipmentStatelessRatesResponse{Rates: &out}
	err = c.post(ctx, "/beta/rates", &req, &res)
	return
}

// Deprecated: Use LowestShipmentRate instead.
// LowestRate gets the lowest rate of a shipment
func (c *Client) LowestRate(shipment *Shipment) (out Rate, err error) {
	// TODO: Why did we deprecate this and replace it with an arguably worse name?
	return c.LowestShipmentRate(shipment)
}

// LowestShipmentRate gets the lowest rate of a shipment
func (c *Client) LowestShipmentRate(shipment *Shipment) (out Rate, err error) {
	return c.LowestShipmentRateWithCarrier(shipment, nil)
}

// LowestStatelessRate gets the lowest stateless rate from a list of stateless rates
func (c *Client) LowestStatelessRate(rates []*StatelessRate) (out StatelessRate, err error) {
	return c.LowestStatelessRateWithCarrier(rates, nil)
}

// Deprecated: Use LowestShipmentRateWithCarrier instead.
// LowestRateWithCarrier performs the same operation as LowestRate,
// but allows specifying a list of carriers for the lowest rate
func (c *Client) LowestRateWithCarrier(shipment *Shipment, carriers []string) (out Rate, err error) {
	// TODO: Why did we deprecate this and replace it with an arguably worse name?
	return c.LowestShipmentRateWithCarrier(shipment, carriers)
}

// LowestShipmentRateWithCarrier performs the same operation as LowestShipmentRate,
// but allows specifying a list of carriers for the lowest rate
func (c *Client) LowestShipmentRateWithCarrier(shipment *Shipment, carriers []string) (out Rate, err error) {
	return c.LowestShipmentRateWithCarrierAndService(shipment, carriers, nil)
}

// LowestStatelessRateWithCarrier performs the same operation as LowestStatelessRate,
// but allows specifying a list of carriers for the lowest rate
func (c *Client) LowestStatelessRateWithCarrier(rates []*StatelessRate, carriers []string) (out StatelessRate, err error) {
	return c.LowestStatelessRateWithCarrierAndService(rates, carriers, nil)
}

// Deprecated: Use LowestShipmentRateWithCarrierAndService instead.
// LowestRateWithCarrierAndService performs the same operation as LowestRate,
// but allows specifying a list of carriers and service for the lowest rate
func (c *Client) LowestRateWithCarrierAndService(shipment *Shipment, carriers []string, services []string) (out Rate, err error) {
	// TODO: Why did we deprecate this and replace it with an arguably worse name?
	return c.LowestShipmentRateWithCarrierAndService(shipment, carriers, services)
}

// LowestShipmentRateWithCarrierAndService performs the same operation as LowestShipmentRate,
// but allows specifying a list of carriers and service for the lowest rate
func (c *Client) LowestShipmentRateWithCarrierAndService(shipment *Shipment, carriers []string, services []string) (out Rate, err error) {
	return c.lowestObjectRate(shipment.Rates, carriers, services)
}

// LowestStatelessRateWithCarrierAndService performs the same operation as LowestStatelessRate,
// but allows specifying a list of carriers and service for the lowest rate
func (c *Client) LowestStatelessRateWithCarrierAndService(rates []*StatelessRate, carriers []string, services []string) (out StatelessRate, err error) {
	return c.lowestStatelessRate(rates, carriers, services)
}

// LowestSmartrate gets the lowest smartrate of a shipment with the specified delivery days and accuracy
func (c *Client) LowestSmartrate(shipment *Shipment, deliveryDays int, deliveryAccuracy string) (out SmartRate, err error) {
	smartrates, _ := c.GetShipmentSmartrates(shipment.ID)
	return c.lowestSmartRate(smartrates, deliveryDays, deliveryAccuracy)
}

// GenerateShipmentForm generates a form of a given type for a shipment
func (c *Client) GenerateShipmentForm(shipmentID string, formType string) (out *Shipment, err error) {
	return c.GenerateShipmentFormWithContext(context.Background(), shipmentID, formType)
}

// GenerateShipmentFormWithContext performs the same operation as GenerateShipmentForm,
// but allows specifying a context that can interrupt the request.
func (c *Client) GenerateShipmentFormWithContext(ctx context.Context, shipmentID string, formType string) (out *Shipment, err error) {
	return c.GenerateShipmentFormWithOptionsWithContext(ctx, shipmentID, formType, make(map[string]interface{}))
}

// GenerateShipmentFormWithOptions generates a form of a given type for a shipment, using provided options
func (c *Client) GenerateShipmentFormWithOptions(shipmentID string, formType string, formOptions map[string]interface{}) (out *Shipment, err error) {
	return c.GenerateShipmentFormWithOptionsWithContext(context.Background(), shipmentID, formType, formOptions)
}

// GenerateShipmentFormWithOptionsWithContext performs the same operation as GenerateShipmentFormWithOptions,
// but allows specifying a context that can interrupt the request.
func (c *Client) GenerateShipmentFormWithOptionsWithContext(ctx context.Context, shipmentID string, formType string, formOptions map[string]interface{}) (out *Shipment, err error) {
	formOptions["type"] = formType
	req := &generateFormRequest{Form: formOptions}
	err = c.post(ctx, "shipments/"+shipmentID+"/forms", &req, &out)
	return
}
