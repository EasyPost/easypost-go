package easypost

// ShipmentOptions represents the various options that can be applied to a
// shipment at creation.
type ShipmentOptions struct {
	AdditionalHandling          bool      `json:"additional_handling,omitempty" url:"additional_handling,omitempty"`
	AddressValidationLevel      string    `json:"address_validation_level,omitempty" url:"address_validation_level,omitempty"`
	Alcohol                     bool      `json:"alcohol,omitempty" url:"alcohol,omitempty"`
	BillingRef                  string    `json:"billing_ref,omitempty" url:"billing_ref,omitempty"`
	BillReceiverAccount         string    `json:"bill_receiver_account,omitempty" url:"bill_receiver_account,omitempty"`
	BillReceiverPostalCode      string    `json:"bill_receiver_postal_code,omitempty" url:"bill_receiver_postal_code,omitempty"`
	BillThirdPartyAccount       string    `json:"bill_third_party_account,omitempty" url:"bill_third_party_account,omitempty"`
	BillThirdPartyCountry       string    `json:"bill_third_party_country,omitempty" url:"bill_third_party_country,omitempty"`
	BillThirdPartyPostalCode    string    `json:"bill_third_party_postal_code,omitempty" url:"bill_third_party_postal_code,omitempty"`
	ByDrone                     bool      `json:"by_drone,omitempty" url:"by_drone,omitempty"`
	CarbonNeutral               bool      `json:"carbon_neutral,omitempty" url:"carbon_neutral,omitempty"`
	CarrierInsuranceAmount      string    `json:"carrier_insurance_amount,omitempty" url:"carrier_insurance_amount,omitempty"`
	CertifiedMail               bool      `json:"certified_mail,omitempty" url:"certified_mail,omitempty"`
	CODAmount                   string    `json:"cod_amount,omitempty" url:"cod_amount,omitempty"`
	CODMethod                   string    `json:"cod_method,omitempty" url:"cod_method,omitempty"`
	CODAddressID                string    `json:"cod_address_id,omitempty" url:"cod_address_id,omitempty"`
	CommercialInvoiceSignature  string    `json:"commercial_invoice_signature,omitempty" url:"commercial_invoice_signature,omitempty"`
	CommercialInvoiceLetterhead string    `json:"commercial_invoice_letterhead,omitempty" url:"commercial_invoice_letterhead,omitempty"`
	ContentDescription          string    `json:"content_description,omitempty" url:"content_description,omitempty"`
	Currency                    string    `json:"currency,omitempty" url:"currency,omitempty"`
	DeliveryConfirmation        string    `json:"delivery_confirmation,omitempty" url:"delivery_confirmation,omitempty"`
	DeliveryMaxDatetime         *DateTime `json:"delivery_max_datetime,omitempty" url:"delivery_max_datetime,omitempty"`
	DutyPayment                 *Payment  `json:"duty_payment,omitempty" url:"duty_payment,omitempty"`
	DutyPaymentAccount          string    `json:"duty_payment_account,omitempty" url:"duty_payment_account,omitempty"`
	DropoffMaxDatetime          *DateTime `json:"dropoff_max_datetime,omitempty" url:"dropoff_max_datetime,omitempty"`
	DropoffType                 string    `json:"dropoff_type,omitempty" url:"dropoff_type,omitempty"`
	DryIce                      bool      `json:"dry_ice,omitempty" url:"dry_ice,omitempty"`
	DryIceMedical               bool      `json:"dry_ice_medical,omitempty,string" url:"dry_ice_medical,omitempty"`
	DryIceWeight                float64   `json:"dry_ice_weight,omitempty,string" url:"dry_ice_weight,omitempty"`
	Endorsement                 string    `json:"endorsement,omitempty" url:"endorsement,omitempty"`
	EndShipperID                string    `json:"end_shipper_id,omitempty" url:"end_shipper_id,omitempty"`
	FreightCharge               float64   `json:"freight_charge,omitempty" url:"freight_charge,omitempty"`
	HandlingInstructions        string    `json:"handling_instructions,omitempty" url:"handling_instructions,omitempty"`
	Hazmat                      string    `json:"hazmat,omitempty" url:"hazmat,omitempty"`
	HoldForPickup               bool      `json:"hold_for_pickup,omitempty" url:"hold_for_pickup,omitempty"`
	Incoterm                    string    `json:"incoterm,omitempty" url:"incoterm,omitempty"`
	InvoiceNumber               string    `json:"invoice_number,omitempty" url:"invoice_number,omitempty"`
	LabelDate                   *DateTime `json:"label_date,omitempty" url:"label_date,omitempty"`
	LabelFormat                 string    `json:"label_format,omitempty" url:"label_format,omitempty"`
	LabelSize                   string    `json:"label_size,omitempty" url:"label_size,omitempty"`
	Machinable                  bool      `json:"machinable,omitempty" url:"machinable,omitempty"`
	Payment                     *Payment  `json:"payment,omitempty" url:"payment,omitempty"`
	PickupMinDatetime           *DateTime `json:"pickup_min_datetime,omitempty" url:"pickup_min_datetime,omitempty"`
	PickupMaxDatetime           *DateTime `json:"pickup_max_datetime,omitempty" url:"pickup_max_datetime,omitempty"`
	PrintCustom1                string    `json:"print_custom_1,omitempty" url:"print_custom_1,omitempty"`
	PrintCustom2                string    `json:"print_custom_2,omitempty" url:"print_custom_2,omitempty"`
	PrintCustom3                string    `json:"print_custom_3,omitempty" url:"print_custom_3,omitempty"`
	PrintCustom1BarCode         bool      `json:"print_custom_1_barcode,omitempty" url:"print_custom_1_barcode,omitempty"`
	PrintCustom2BarCode         bool      `json:"print_custom_2_barcode,omitempty" url:"print_custom_2_barcode,omitempty"`
	PrintCustom3BarCode         bool      `json:"print_custom_3_barcode,omitempty" url:"print_custom_3_barcode,omitempty"`
	PrintCustom1Code            string    `json:"print_custom_1_code,omitempty" url:"print_custom_1_code,omitempty"`
	PrintCustom2Code            string    `json:"print_custom_2_code,omitempty" url:"print_custom_2_code,omitempty"`
	PrintCustom3Code            string    `json:"print_custom_3_code,omitempty" url:"print_custom_3_code,omitempty"`
	RegisteredMail              bool      `json:"registered_mail,omitempty" url:"registered_mail,omitempty"`
	RegisteredMailAmount        float64   `json:"registered_mail_amount,omitempty" url:"registered_mail_amount,omitempty"`
	ReturnReceipt               bool      `json:"return_receipt,omitempty" url:"return_receipt,omitempty"`
	SaturdayDelivery            bool      `json:"saturday_delivery,omitempty" url:"saturday_delivery,omitempty"`
	SpecialRatesEligibility     string    `json:"special_rates_eligibility,omitempty" url:"special_rates_eligibility,omitempty"`
	SmartpostHub                string    `json:"smartpost_hub,omitempty" url:"smartpost_hub,omitempty"`
	SmartpostManifest           string    `json:"smartpost_manifest,omitempty" url:"smartpost_manifest,omitempty"`
	SuppressETD                 bool      `json:"suppress_etd,omitempty" url:"suppress_etd,omitempty"`
}

// Payment provides information on how a shipment is billed.
type Payment struct {
	Type       string `json:"type,omitempty" url:"type,omitempty"`
	Account    string `json:"account,omitempty" url:"account,omitempty"`
	Country    string `json:"country,omitempty" url:"country,omitempty"`
	PostalCode string `json:"postal_code,omitempty" url:"postal_code,omitempty"`
}
