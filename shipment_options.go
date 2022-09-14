package easypost

import "time"

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
	CODAmount                string     `json:"cod_amount,omitempty"`
	CODMethod                string     `json:"cod_method,omitempty"`
	CODAddressID             string     `json:"cod_address_id,omitempty"`
	Currency                 string     `json:"currency,omitempty"`
	DeliveryConfirmation     string     `json:"delivery_confirmation,omitempty"`
	DutyPayment              *Payment   `json:"duty_payment,omitempty"`
	DutyPaymentAccount       string     `json:"duty_payment_account,omitempty"`
	DropoffType              string     `json:"dropoff_type,omitempty"`
	DryIce                   bool       `json:"dry_ice,omitempty"`
	DryIceMedical            bool       `json:"dry_ice_medical,omitempty,string"`
	DryIceWeight             float64    `json:"dry_ice_weight,omitempty,string"`
	Endorsement              string     `json:"endorsement,omitempty"`
	EndShipperID             string     `json:"end_shipper_id,omitempty"`
	FreightCharge            float64    `json:"freight_charge,omitempty"`
	HandlingInstructions     string     `json:"handling_instructions,omitempty"`
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

// Payment provides information on how a shipment is billed.
type Payment struct {
	Type       string `json:"type,omitempty"`
	Account    string `json:"account,omitempty"`
	Country    string `json:"country,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
}
