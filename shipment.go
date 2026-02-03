package easypost

import (
	"context"
	"net/http"
)

// A Form represents a form associated with a Shipment.
type Form struct {
	ID                      string    `json:"id,omitempty" url:"id,omitempty"`
	Object                  string    `json:"object,omitempty" url:"object,omitempty"`
	Mode                    string    `json:"mode,omitempty" url:"mode,omitempty"`
	CreatedAt               *DateTime `json:"created_at,omitempty" url:"created_at,omitempty"`
	UpdatedAt               *DateTime `json:"updated_at,omitempty" url:"updated_at,omitempty"`
	FormType                string    `json:"form_type,omitempty" url:"form_type,omitempty"`
	FormURL                 string    `json:"form_url,omitempty" url:"form_url,omitempty"`
	SubmittedElectronically bool      `json:"submitted_electronically,omitempty" url:"submitted_electronically,omitempty"`
}

// PostageLabel provides details of a shipping label for a purchased shipment.
type PostageLabel struct {
	ID              string    `json:"id,omitempty" url:"id,omitempty"`
	Object          string    `json:"object,omitempty" url:"object,omitempty"`
	CreatedAt       *DateTime `json:"created_at,omitempty" url:"created_at,omitempty"`
	UpdatedAt       *DateTime `json:"updated_at,omitempty" url:"updated_at,omitempty"`
	IntegratedForm  string    `json:"integrated_form,omitempty" url:"integrated_form,omitempty"`
	LabelDate       *DateTime `json:"label_date,omitempty" url:"label_date,omitempty"`
	LabelEPL2URL    string    `json:"label_epl2_url,omitempty" url:"label_epl2_url,omitempty"`
	LabelFileType   string    `json:"label_file_type,omitempty" url:"label_file_type,omitempty"`
	LabelPDFURL     string    `json:"label_pdf_url,omitempty" url:"label_pdf_url,omitempty"`
	LabelResolution float64   `json:"label_resolution,omitempty" url:"label_resolution,omitempty"`
	LabelSize       string    `json:"label_size,omitempty" url:"label_size,omitempty"`
	LabelType       string    `json:"label_type,omitempty" url:"label_type,omitempty"`
	LabelURL        string    `json:"label_url,omitempty" url:"label_url,omitempty"`
	LabelZPLURL     string    `json:"label_zpl_url,omitempty" url:"label_zpl_url,omitempty"`
}

type LineItem struct {
	TotalLineValue  string `json:"total_line_value,omitempty" url:"total_line_value,omitempty"`
	ItemDescription string `json:"item_description,omitempty" url:"item_description,omitempty"`
}

// A Shipment represents its namesake, and is made up of a "to" and "from"
// addresses, the Parcel being shipped, and any customs forms required for
// international deliveries.
type Shipment struct {
	ID                string            `json:"id,omitempty" url:"id,omitempty"`
	Object            string            `json:"object,omitempty" url:"object,omitempty"`
	Reference         string            `json:"reference,omitempty" url:"reference,omitempty"`
	Mode              string            `json:"mode,omitempty" url:"mode,omitempty"`
	CreatedAt         *DateTime         `json:"created_at,omitempty" url:"created_at,omitempty"`
	UpdatedAt         *DateTime         `json:"updated_at,omitempty" url:"updated_at,omitempty"`
	ToAddress         *Address          `json:"to_address,omitempty" url:"to_address,omitempty"`
	FromAddress       *Address          `json:"from_address,omitempty" url:"from_address,omitempty"`
	ReturnAddress     *Address          `json:"return_address,omitempty" url:"return_address,omitempty"`
	BuyerAddress      *Address          `json:"buyer_address,omitempty" url:"buyer_address,omitempty"`
	Parcel            *Parcel           `json:"parcel,omitempty" url:"parcel,omitempty"`
	Carrier           string            `json:"carrier,omitempty" url:"carrier,omitempty"`
	Service           string            `json:"service,omitempty" url:"service,omitempty"`
	CarrierAccountIDs []string          `json:"carrier_accounts,omitempty" url:"carrier_accounts,omitempty"`
	CustomsInfo       *CustomsInfo      `json:"customs_info,omitempty" url:"customs_info,omitempty"`
	ScanForm          *ScanForm         `json:"scan_form,omitempty" url:"scan_form,omitempty"`
	Forms             []*Form           `json:"forms,omitempty" url:"forms,omitempty"`
	Insurance         string            `json:"insurance,omitempty" url:"insurance,omitempty"`
	Rates             []*Rate           `json:"rates,omitempty" url:"rates,omitempty"`
	SelectedRate      *Rate             `json:"selected_rate,omitempty" url:"selected_rate,omitempty"`
	PostageLabel      *PostageLabel     `json:"postage_label,omitempty" url:"postage_label,omitempty"`
	Messages          []*CarrierMessage `json:"messages,omitempty" url:"messages,omitempty"`
	Options           *ShipmentOptions  `json:"options,omitempty" url:"options,omitempty"`
	IsReturn          bool              `json:"is_return,omitempty" url:"is_return,omitempty"`
	TrackingCode      string            `json:"tracking_code,omitempty" url:"tracking_code,omitempty"`
	USPSZone          int               `json:"usps_zone,omitempty" url:"usps_zone,omitempty"`
	Status            string            `json:"status,omitempty" url:"status,omitempty"`
	Tracker           *Tracker          `json:"tracker,omitempty" url:"tracker,omitempty"`
	Fees              []*Fee            `json:"fees,omitempty" url:"fees,omitempty"`
	RefundStatus      string            `json:"refund_status,omitempty" url:"refund_status,omitempty"`
	BatchID           string            `json:"batch_id,omitempty" url:"batch_id,omitempty"`
	BatchStatus       string            `json:"batch_status,omitempty" url:"batch_status,omitempty"`
	BatchMessage      string            `json:"batch_message,omitempty" url:"batch_message,omitempty"`
	TaxIdentifiers    []*TaxIdentifier  `json:"tax_identifiers,omitempty" url:"tax_identifiers,omitempty"`
	LineItems         []*LineItem       `json:"line_items,omitempty" url:"line_items,omitempty"`
}

// ListShipmentsOptions is used to specify query parameters for listing Shipment
// objects.
type ListShipmentsOptions struct {
	BeforeID        string    `json:"before_id,omitempty" url:"before_id,omitempty"`
	AfterID         string    `json:"after_id,omitempty" url:"after_id,omitempty"`
	StartDateTime   *DateTime `json:"start_datetime,omitempty" url:"start_datetime,omitempty"`
	EndDateTime     *DateTime `json:"end_datetime,omitempty" url:"end_datetime,omitempty"`
	PageSize        int       `json:"page_size,omitempty" url:"page_size,omitempty"`
	Purchased       *bool     `json:"purchased,omitempty" url:"purchased,omitempty"`
	IncludeChildren *bool     `json:"include_children,omitempty" url:"include_children,omitempty"`
}

// ListShipmentsResult holds the results from the list shipments API.
type ListShipmentsResult struct {
	Shipments       []*Shipment `json:"shipments,omitempty" url:"shipments,omitempty"`
	Purchased       *bool       `json:"purchased,omitempty" url:"purchased,omitempty"`
	IncludeChildren *bool       `json:"include_children,omitempty" url:"include_children,omitempty"`
	PaginatedCollection
}

type buyShipmentRequest struct {
	Rate         *Rate  `json:"rate,omitempty" url:"rate,omitempty"`
	Insurance    string `json:"insurance,omitempty" url:"insurance,omitempty"`
	EndShipperID string `json:"end_shipper_id,omitempty" url:"end_shipper_id,omitempty"`
}

type createShipmentRequest struct {
	Shipment *Shipment `json:"shipment,omitempty" url:"shipment,omitempty"`
}

type getShipmentRatesResponse struct {
	Rates *[]*Rate `json:"rates,omitempty" url:"rates,omitempty"`
}

type generateFormRequest struct {
	Form map[string]interface{} `json:"form,omitempty" url:"form,omitempty"`
}

type EasyPostTimeInTransitData struct {
	DaysInTransit                 TimeInTransit `json:"days_in_transit,omitempty" url:"days_in_transit,omitempty"`
	EasyPostEstimatedDeliveryDate string        `json:"easypost_estimated_delivery_date,omitempty" url:"easypost_estimated_delivery_date,omitempty"`
	PlannedShipDate               string        `json:"planned_ship_date,omitempty" url:"planned_ship_date,omitempty"`
}

type EstimatedDeliveryDate struct {
	EasyPostTimeInTransitData EasyPostTimeInTransitData `json:"easypost_time_in_transit_data,omitempty" url:"easypost_time_in_transit_data,omitempty"`
	Rate                      SmartRate                 `json:"rate,omitempty" url:"rate,omitempty"`
}

// RecommendShipDateForShipmentResult is the result of the RecommendShipDateForShipment method.
type RecommendShipDateForShipmentResult struct {
	Rate                      *SmartRate                       `json:"rate,omitempty" url:"rate,omitempty"`
	EasyPostTimeInTransitData *TimeInTransitDetailsForShipDate `json:"easypost_time_in_transit_data,omitempty" url:"easypost_time_in_transit_data,omitempty"`
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

// CreateShipmentWithContext performs the same operation as CreateShipment, but
// allows specifying a context that can interrupt the request.
func (c *Client) CreateShipmentWithContext(ctx context.Context, in *Shipment) (out *Shipment, err error) {
	req := &createShipmentRequest{Shipment: in}
	err = c.do(ctx, http.MethodPost, "shipments", &req, &out)
	return
}

// ListShipments provides a paginated result of Shipment objects.
func (c *Client) ListShipments(opts *ListShipmentsOptions) (out *ListShipmentsResult, err error) {
	return c.ListShipmentsWithContext(context.Background(), opts)
}

// ListShipmentsWithContext performs the same operation as ListShipments, but
// allows specifying a context that can interrupt the request.
func (c *Client) ListShipmentsWithContext(ctx context.Context, opts *ListShipmentsOptions) (out *ListShipmentsResult, err error) {
	err = c.do(ctx, http.MethodGet, "shipments", opts, &out)
	// Store the original query parameters for reuse when getting the next page
	out.Purchased = opts.Purchased
	out.IncludeChildren = opts.IncludeChildren
	return
}

// GetNextShipmentPage returns the next page of shipments
func (c *Client) GetNextShipmentPage(collection *ListShipmentsResult) (out *ListShipmentsResult, err error) {
	return c.GetNextShipmentPageWithContext(context.Background(), collection)
}

// GetNextShipmentPageWithPageSize returns the next page of shipments with a specific page size
func (c *Client) GetNextShipmentPageWithPageSize(collection *ListShipmentsResult, pageSize int) (out *ListShipmentsResult, err error) {
	return c.GetNextShipmentPageWithPageSizeWithContext(context.Background(), collection, pageSize)
}

// GetNextShipmentPageWithContext performs the same operation as GetNextShipmentPage, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetNextShipmentPageWithContext(ctx context.Context, collection *ListShipmentsResult) (out *ListShipmentsResult, err error) {
	return c.GetNextShipmentPageWithPageSizeWithContext(ctx, collection, 0)
}

// GetNextShipmentPageWithPageSizeWithContext performs the same operation as GetNextShipmentPageWithPageSize, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetNextShipmentPageWithPageSizeWithContext(ctx context.Context, collection *ListShipmentsResult, pageSize int) (out *ListShipmentsResult, err error) {
	if len(collection.Shipments) == 0 {
		err = newEndOfPaginationError()
		return
	}
	lastID := collection.Shipments[len(collection.Shipments)-1].ID
	params, err := nextPageParameters(collection.HasMore, lastID, pageSize)
	if err != nil {
		return
	}
	shipmentParams := &ListShipmentsOptions{
		BeforeID:        params.BeforeID,
		Purchased:       collection.Purchased,
		IncludeChildren: collection.IncludeChildren,
	}
	if pageSize > 0 {
		shipmentParams.PageSize = pageSize
	}
	return c.ListShipmentsWithContext(ctx, shipmentParams)
}

// GetShipment retrieves a Shipment object by ID.
func (c *Client) GetShipment(shipmentID string) (out *Shipment, err error) {
	return c.GetShipmentWithContext(context.Background(), shipmentID)
}

// GetShipmentWithContext performs the same operation as GetShipment, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetShipmentWithContext(ctx context.Context, shipmentID string) (out *Shipment, err error) {
	err = c.do(ctx, http.MethodGet, "shipments/"+shipmentID, nil, &out)
	return
}

func (c *Client) buyShipment(ctx context.Context, shipmentID string, in *buyShipmentRequest) (out *Shipment, err error) {
	err = c.do(ctx, http.MethodPost, "shipments/"+shipmentID+"/buy", &in, &out)
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

// GetShipmentLabel enables retrieving the label for a shipment in a different
// format. The PostageLabel field in the returned Shipment object will reflect
// the new format.
func (c *Client) GetShipmentLabel(shipmentID, format string) (out *Shipment, err error) {
	return c.GetShipmentLabelWithContext(context.Background(), shipmentID, format)
}

// GetShipmentLabelWithContext performs the same operation as GetShipmentLabel,
// but allows specifying a context that can interrupt the request.
func (c *Client) GetShipmentLabelWithContext(ctx context.Context, shipmentID, format string) (out *Shipment, err error) {
	req := struct {
		FileFormat string `json:"file_format,omitempty" url:"file_format,omitempty"`
	}{}
	req.FileFormat = format
	err = c.do(ctx, http.MethodGet, "shipments/"+shipmentID+"/label", req, &out)
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
		Smartrates *[]*SmartRate `json:"result,omitempty" url:"result,omitempty"`
	}{Smartrates: &out}
	err = c.do(ctx, http.MethodGet, "shipments/"+shipmentID+"/smartrate", nil, &res)
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
	params := struct {
		Amount string `json:"amount,omitempty" url:"amount,omitempty"`
	}{Amount: amount}
	err = c.do(ctx, http.MethodPost, "shipments/"+shipmentID+"/insure", params, &out)
	return
}

// RefundShipment requests a refund from the carrier.
func (c *Client) RefundShipment(shipmentID string) (out *Shipment, err error) {
	return c.RefundShipmentWithContext(context.Background(), shipmentID)
}

// RefundShipmentWithContext performs the same operation as RefundShipment, but
// allows specifying a context that can interrupt the request.
func (c *Client) RefundShipmentWithContext(ctx context.Context, shipmentID string) (out *Shipment, err error) {
	err = c.do(ctx, http.MethodPost, "shipments/"+shipmentID+"/refund", nil, &out)
	return
}

// RerateShipment fetches the available rates for a shipment with the current rates.
func (c *Client) RerateShipment(shipmentID string) (out []*Rate, err error) {
	return c.RerateShipmentWithContext(context.Background(), shipmentID)
}

// RerateShipmentWithContext performs the same operation as RerateShipment,
// but allows specifying a context that can interrupt the request.
func (c *Client) RerateShipmentWithContext(ctx context.Context, shipmentID string) (out []*Rate, err error) {
	res := &getShipmentRatesResponse{Rates: &out}
	err = c.do(ctx, http.MethodPost, "shipments/"+shipmentID+"/rerate", nil, &res)
	return
}

// LowestShipmentRate gets the lowest rate of a shipment
func (c *Client) LowestShipmentRate(shipment *Shipment) (out Rate, err error) {
	return c.LowestShipmentRateWithCarrier(shipment, nil)
}

// LowestShipmentRateWithCarrier performs the same operation as LowestShipmentRate,
// but allows specifying a list of carriers for the lowest rate
func (c *Client) LowestShipmentRateWithCarrier(shipment *Shipment, carriers []string) (out Rate, err error) {
	return c.LowestShipmentRateWithCarrierAndService(shipment, carriers, nil)
}

// LowestShipmentRateWithCarrierAndService performs the same operation as LowestShipmentRate,
// but allows specifying a list of carriers and service for the lowest rate
func (c *Client) LowestShipmentRateWithCarrierAndService(shipment *Shipment, carriers []string, services []string) (out Rate, err error) {
	return c.lowestObjectRate(shipment.Rates, carriers, services)
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
	err = c.do(ctx, http.MethodPost, "shipments/"+shipmentID+"/forms", &req, &out)
	return
}

// GetShipmentEstimatedDeliveryDate retrieves the estimated delivery date of each rate for a Shipment via the Delivery Date Estimator API, based on a specific ship date.
func (c *Client) GetShipmentEstimatedDeliveryDate(shipmentID string, plannedShipDate string) (out []*EstimatedDeliveryDate, err error) {
	return c.GetShipmentEstimatedDeliveryDateWithContext(context.Background(), shipmentID, plannedShipDate)
}

// GetShipmentEstimatedDeliveryDateWithContext performs the same operation as EstimateDeliveryDateForShipment, but allows specifying a context that can interrupt the request.
func (c *Client) GetShipmentEstimatedDeliveryDateWithContext(ctx context.Context, shipmentID string, plannedShipDate string) (out []*EstimatedDeliveryDate, err error) {
	params := struct {
		PlannedShipDate string `json:"planned_ship_date,omitempty" url:"planned_ship_date,omitempty"`
	}{PlannedShipDate: plannedShipDate}
	res := struct {
		EstimatedDeliveryDates *[]*EstimatedDeliveryDate `json:"rates,omitempty" url:"rates,omitempty"`
	}{EstimatedDeliveryDates: &out}
	err = c.do(ctx, http.MethodGet, "shipments/"+shipmentID+"/smartrate/delivery_date", params, &res)
	return
}

// RecommendShipDateForShipment retrieves the recommended ship date of each rate for a Shipment via the Precision Shipping API, based on a specific desired delivery date.
func (c *Client) RecommendShipDateForShipment(shipmentID string, desiredDeliveryDate string) (out []*RecommendShipDateForShipmentResult, err error) {
	return c.RecommendShipDateForShipmentWithContext(context.Background(), shipmentID, desiredDeliveryDate)
}

// RecommendShipDateForShipmentWithContext performs the same operation as RecommendShipDateForShipment, but allows specifying a context that can interrupt the request.
func (c *Client) RecommendShipDateForShipmentWithContext(ctx context.Context, shipmentID string, desiredDeliveryDate string) (out []*RecommendShipDateForShipmentResult, err error) {
	params := struct {
		DesiredDeliveryDate string `json:"desired_delivery_date,omitempty" url:"desired_delivery_date,omitempty"`
	}{
		DesiredDeliveryDate: desiredDeliveryDate,
	}
	res := struct {
		Results *[]*RecommendShipDateForShipmentResult `json:"rates,omitempty" url:"rates,omitempty"`
	}{Results: &out}
	err = c.do(ctx, http.MethodGet, "shipments/"+shipmentID+"/smartrate/precision_shipping", params, &res)
	return
}

// CreateAndBuyLumaShipment creates and buys a Shipment using Luma.
func (c *Client) CreateAndBuyLumaShipment(in *LumaRequest) (out *Shipment, err error) {
	return c.CreateAndBuyLumaShipmentWithContext(context.Background(), in)
}

// CreateAndBuyLumaShipmentWithContext performs the same operation as CreateAndBuyLumaShipment, but allows specifying a context that can interrupt the request.
func (c *Client) CreateAndBuyLumaShipmentWithContext(ctx context.Context, in *LumaRequest) (out *Shipment, err error) {
	req := &WrappedLumaRequest{Shipment: in}
	err = c.do(ctx, http.MethodPost, "shipments/luma", &req, &out)
	return
}

// BuyLumaShipment buys a Shipment using Luma.
func (c *Client) BuyLumaShipment(shipmentID string, in *LumaRequest) (out *Shipment, err error) {
	return c.BuyLumaShipmentWithContext(context.Background(), shipmentID, in)
}

// BuyLumaShipmentWithContext performs the same operation as BuyLumaShipment, but allows specifying a context that can interrupt the request.
func (c *Client) BuyLumaShipmentWithContext(ctx context.Context, shipmentID string, in *LumaRequest) (out *Shipment, err error) {
	err = c.do(ctx, http.MethodPost, "shipments/"+shipmentID+"/luma", in, &out)
	return
}
