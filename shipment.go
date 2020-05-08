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
	EstDeliveryDays        int        `json:"est_delivery_dats,omitempty"`
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

// Payment provides information on how a shipment is billed.
type Payment struct {
	Type       string `json:"type,omitempty"`
	Account    string `json:"account,omitempty"`
	Country    string `json:"country,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
}

// ShipmentOptions represents the various options that can be applied to a
// shipment at creation.
type ShipmentOptions struct {
	AdditionalHandling       bool       `json:"additional_handling,omitempty"`
	AddressValidationLevel   string     `json:"address_validation_level,omitempty"`
	Alcohol                  bool       `json:"alcohol,omitempty"`
	BillReceiverAccount      string     `json:"bill_receiver_account,omitempty"`
	BillReceiverPostalCode   string     `json:"bill_receiver_postal_code,omitempty"`
	BillThirdPartyAccount    string     `json:"bill_third_party_account,omitempty"`
	BillThirdPartyCountry    string     `json:"bill_third_party_country,omitempty"`
	BillThirdPartyPostalCode string     `json:"bill_third_party_postal_code,omitempty"`
	ByDrone                  bool       `json:"by_drone,omitempty"`
	CarbonNeutral            bool       `json:"carbon_neutral,omitempty"`
	CertifiedMail            bool       `json:"certified_mail,omitempty"`
	CODAmount                float64    `json:"cod_amount,omitempty,string"`
	CODMethod                string     `json:"cod_method,omitempty"`
	CODAddressID             string     `json:"cod_address_id,omitempty"`
	Currency                 string     `json:"currency,omitempty"`
	DeliveryConfirmation     string     `json:"delivery_confirmation,omitempty"`
	DropoffType              string     `json:"dropoff_type,omitempty"`
	DryIce                   bool       `json:"dry_ice,omitempty"`
	DryIceMedical            bool       `json:"dry_ice_medical,omitempty,string"`
	DryIceWeight             float64    `json:"dry_ice_weight,omitempty,string"`
	Endorsement              string     `json:"endorsement,omitempty"`
	FreightCharge            float64    `json:"freight_charge,omitempty"`
	HandlingInstructions     string     `json:"handling_insructions,omitempty"`
	Hazmat                   string     `json:"hazmat,omitempty"`
	HoldForPickup            bool       `json:"hold_for_pickup,omitempty"`
	Incoterm                 string     `json:"incoterm,omitempty"`
	InvoiceNumber            string     `json:"invoice_number,omitempty"`
	LabelDate                *time.Time `json:"label_date,omitempty"`
	LabelFormat              string     `json:"label_format,omitempty"`
	Machinable               bool       `json:"machinable,omitempty"`
	Payment                  *Payment   `json:"payment,omitempty"`
	PrintCustom1             string     `json:"print_custom_1,omitempty"`
	PrintCustom2             string     `json:"print_custom_2,omitempty"`
	PrintCustom3             string     `json:"print_custom_3,omitempty"`
	PrintCustom1BarCode      bool       `json:"print_custom_1_barcode,omitempty"`
	PrintCustom2BarCode      bool       `json:"print_custom_2_barcode,omitempty"`
	PrintCustom3BarCode      bool       `json:"print_custom_3_barcode,omitempty"`
	PrintCustom1Code         string     `json:"print_custom_1_code,omitempty"`
	PrintCustom2Code         string     `json:"print_custom_2_code,omitempty"`
	PrintCustom3Code         string     `json:"print_custom_3_code,omitempty"`
	RegisteredMail           bool       `json:"registered_mail,omitempty"`
	RegisteredMailAmount     float64    `json:"registered_mail_amount,omitempty"`
	ReturnReceipt            bool       `json:"return_receipt,omitempty"`
	SaturdayDelivery         bool       `json:"saturday_delivery,omitempty"`
	SpecialRatesEligibility  string     `json:"special_rates_eligibility,omitempty"`
	SmartpostHub             string     `json:"smartpost_hub,omitempty"`
	SmartpostManifest        string     `json:"smartpost_manifest,omitempty"`
	BillingRef               string     `json:"billing_ref,omitempty"`
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
}

// CreateShipment creates a new Shipment object. The ToAddress, FromAddress and
// Parcel attributes are required. These objects may be fully-specified to
// create new ones at the same time as creating the Shipment, or they can refer
// to existing objects via their ID attribute. Passing in one or more carrier
// accounts to CreateShipment limits the returned rates to the specified
// carriers.
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
	req := struct {
		Shipment *Shipment `json:"shipment"`
	}{Shipment: in}
	err = c.post(nil, "shipments", &req, &out)
	return
}

// CreateShipmentWithContext performs the same operation as CreateShipment, but
// allows specifying a context that can interrupt the request.
func (c *Client) CreateShipmentWithContext(ctx context.Context, in *Shipment) (out *Shipment, err error) {
	req := struct {
		Shipment *Shipment `json:"shipment"`
	}{Shipment: in}
	err = c.post(ctx, "shipments", &req, &out)
	return
}

// ListShipmentsOptions is used to specify query parameters for listing Shipment
// objects.
type ListShipmentsOptions struct {
	BeforeID        string     `json:"before_id,omitempty"`
	AfterID         string     `json:"after_id,omitempty"`
	StartDateTime   *time.Time `json:"start_datetime,omitempty"`
	EndDateTime     *time.Time `json:"end_datetime,omitempty"`
	PageSize        int        `json:"page_size,omitempty"`
	Purchased       bool       `json:"purchased,omitempty"`
	IncludeChildren bool       `json:"include_children,omitempty"`
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

// ListShipments provides a paginated result of Shipment objects.
func (c *Client) ListShipments(opts *ListShipmentsOptions) (out *ListShipmentsResult, err error) {
	err = c.do(nil, http.MethodGet, "shipments", &opts, &out)
	return
}

// ListShipmentsWithContext performs the same operation as ListShipments, but
// allows specifying a context that can interrupt the request.
func (c *Client) ListShipmentsWithContext(ctx context.Context, opts *ListShipmentsOptions) (out *ListShipmentsResult, err error) {
	err = c.do(ctx, http.MethodGet, "shipments", &opts, &out)
	return
}

// GetShipment retrieves a Shipment object by ID.
func (c *Client) GetShipment(shipmentID string) (out *Shipment, err error) {
	err = c.get(nil, "shipments/"+shipmentID, &out)
	return
}

// GetShipmentWithContext performs the same operation as GetShipment, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetShipmentWithContext(ctx context.Context, shipmentID string) (out *Shipment, err error) {
	err = c.get(ctx, "shipments/"+shipmentID, &out)
	return
}

type buyShipmentRequest struct {
	Rate      *Rate  `json:"rate,omitempty"`
	Insurance string `json:"insurance,omitempty"`
}

// BuyShipment purchases a shipment. If successful, the returned Shipment will
// have the PostageLabel attribute populated.
//	c := easypost.New(MyEasyPostAPIKey)
//	out, err := c.Buy("shp_100", &easypost.Rate{ID: "rate_1001"}, "249.99")
func (c *Client) BuyShipment(shipmentID string, rate *Rate, insurance string) (out *Shipment, err error) {
	req := &buyShipmentRequest{Rate: rate, Insurance: insurance}
	err = c.post(nil, "shipments/"+shipmentID+"/buy", req, &out)
	return
}

// BuyShipmentWithContext performs the same operation as BuyShipment, but allows
// specifying a context that can interrupt the request.
func (c *Client) BuyShipmentWithContext(ctx context.Context, shipmentID string, rate *Rate, insurance string) (out *Shipment, err error) {
	req := &buyShipmentRequest{Rate: rate, Insurance: insurance}
	err = c.post(ctx, "shipments/"+shipmentID+"/buy", &req, &out)
	return
}

// GetShipmentLabel enables retrieving the label for a shipment in a different
// format. The PostageLabel field in the returned Shipment object will reflect
// the new format.
func (c *Client) GetShipmentLabel(shipmentID, format string) (out *Shipment, err error) {
	vals := url.Values{"file_format": []string{format}}
	err = c.do(nil, http.MethodGet, "shipments/"+shipmentID+"/label", vals, &out)
	return
}

// GetShipmentLabelWithContext performs the same operation as GetShipmentLabel,
// but allows specifying a context that can interrupt the request.
func (c *Client) GetShipmentLabelWithContext(ctx context.Context, shipmentID, format string) (out *Shipment, err error) {
	vals := url.Values{"file_format": []string{format}}
	err = c.do(ctx, http.MethodGet, "shipments/"+shipmentID+"/label", vals, &out)
	return
}

type getShipmentRatesResponse struct {
	Rates *[]*Rate `json:"rates,omitempty"`
}

// GetShipmentRates fetches the available rates for a shipment.
func (c *Client) GetShipmentRates(shipmentID string) (out []*Rate, err error) {
	res := &getShipmentRatesResponse{Rates: &out}
	err = c.get(nil, "shipments/"+shipmentID+"/rates", &res)
	return
}

// GetShipmentRatesWithContext performs the same operation as GetShipmentRates,
// but allows specifying a context that can interrupt the request.
func (c *Client) GetShipmentRatesWithContext(ctx context.Context, shipmentID string) (out []*Rate, err error) {
	res := &getShipmentRatesResponse{Rates: &out}
	err = c.get(ctx, "shipments/"+shipmentID+"/rates", &res)
	return
}

// InsureShipment purchases insurance for the shipment. Insurance should be
// purchased after purchasing the shipment, but before it has been processed by
// the carrier. On success, the purchased insurance will be reflected in the
// returned Shipment object's Insurance field.
func (c *Client) InsureShipment(shipmentID, amount string) (out *Shipment, err error) {
	vals := url.Values{"amount": []string{amount}}
	err = c.post(nil, "shipments/"+shipmentID+"/insure", vals, &out)
	return
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
	err = c.post(nil, "shipments/"+shipmentID+"/refund", nil, &out)
	return
}

// RefundShipmentWithContext performs the same operation as RefundShipment, but
// allows specifying a context that can interrupt the request.
func (c *Client) RefundShipmentWithContext(ctx context.Context, shipmentID string) (out *Shipment, err error) {
	err = c.post(ctx, "shipments/"+shipmentID+"/refund", nil, &out)
	return
}
