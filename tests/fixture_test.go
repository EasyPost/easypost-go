package easypost_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/EasyPost/easypost-go/v5"
)

type Fixture struct {
	Addresses       map[string]*easypost.Address               `json:"addresses,omitempty" url:"addresses,omitempty"`
	Billing         billingFixture                             `json:"billing,omitempty" url:"billing,omitempty"`
	CarrierAccounts map[string]*easypost.CarrierAccount        `json:"carrier_accounts,omitempty" url:"carrier_accounts,omitempty"`
	CarrierStrings  map[string]string                          `json:"carrier_strings,omitempty" url:"carrier_strings,omitempty"`
	Claims          map[string]*easypost.CreateClaimParameters `json:"claims,omitempty" url:"claims,omitempty"`
	CustomsInfos    map[string]*easypost.CustomsInfo           `json:"customs_infos,omitempty" url:"customs_infos,omitempty"`
	CustomsItems    map[string]*easypost.CustomsItem           `json:"customs_items,omitempty" url:"customs_items,omitempty"`
	CreditCards     map[string]*easypost.CreditCardOptions     `json:"credit_cards,omitempty" url:"credit_cards,omitempty"`
	FormOptions     map[string]map[string]interface{}          `json:"form_options,omitempty" url:"form_options,omitempty"`
	Insurances      map[string]*easypost.Insurance             `json:"insurances,omitempty" url:"insurances,omitempty"`
	Luma            map[string]interface{}                     `json:"luma,omitempty" url:"luma,omitempty"`
	Orders          map[string]*easypost.Order                 `json:"orders,omitempty" url:"orders,omitempty"`
	PageSizes       map[string]int                             `json:"page_sizes,omitempty" url:"page_sizes,omitempty"`
	Parcels         map[string]*easypost.Parcel                `json:"parcels,omitempty" url:"parcels,omitempty"`
	Pickups         map[string]*easypost.Pickup                `json:"pickups,omitempty" url:"pickups,omitempty"`
	ReportTypes     map[string]string                          `json:"report_types,omitempty" url:"report_types,omitempty"`
	ServiceNames    map[string]map[string]string               `json:"service_names,omitempty" url:"service_names,omitempty"`
	Shipments       map[string]*easypost.Shipment              `json:"shipments,omitempty" url:"shipments,omitempty"`
	TaxIdentifiers  map[string]*easypost.TaxIdentifier         `json:"tax_identifiers,omitempty" url:"tax_identifiers,omitempty"`
	Users           map[string]*easypost.UserOptions           `json:"users,omitempty" url:"users,omitempty"`
	Webhooks        map[string]interface{}                     `json:"webhooks,omitempty" url:"webhooks,omitempty"`
}

type billingFixture struct {
	PaymentMethodID        string                `json:"payment_method_id,omitempty" url:"payment_method_id,omitempty"`
	FinancialConnectionsID string                `json:"financial_connections_id,omitempty" url:"financial_connections_id,omitempty"`
	MandateData            *easypost.MandateData `json:"mandate_data,omitempty" url:"mandate_data,omitempty"`
	Priority               string                `json:"priority,omitempty" url:"priority,omitempty"`
}

// Reads fixture data from the fixtures JSON file
func readFixtureData() Fixture {
	currentDir, _ := os.Getwd()
	parentDir := filepath.Dir(currentDir)
	filePath := fmt.Sprintf("%s%s", parentDir, "/examples/official/fixtures/client-library-fixtures.json")

	/* #nosec */
	data, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error opening fixture file:", err)
		os.Exit(1)
	}

	defer func() { _ = data.Close() }()

	byteData, _ := io.ReadAll(data)

	var fixtures Fixture
	_ = json.Unmarshal([]byte(byteData), &fixtures)

	return fixtures
}

// We keep the page_size of retrieving `all` records small so cassettes stay small
func (fixture *Fixture) pageSize() int {
	return readFixtureData().PageSizes["five_results"]
}

// This is the USPS carrier account ID that comes with your EasyPost account by default and should be used for all tests
func (fixture *Fixture) USPSCarrierAccountID() string {
	uspsCarrierAccountID := os.Getenv("USPS_CARRIER_ACCOUNT_ID")

	// Fallback to the EasyPost Go Client Library Test User USPS carrier account ID
	if len(uspsCarrierAccountID) == 0 {
		return "ca_e606176ddb314afe896733636dba2f3b"
	}

	return uspsCarrierAccountID
}

func (fixture *Fixture) USPS() string {
	return readFixtureData().CarrierStrings["usps"]
}

func (fixture *Fixture) USPSService() string {
	return readFixtureData().ServiceNames["usps"]["first_service"]
}

func (fixture *Fixture) PickupService() string {
	return readFixtureData().ServiceNames["usps"]["pickup_service"]
}

func (fixture *Fixture) ReportType() string {
	return readFixtureData().ReportTypes["shipment"]
}

func (fixture *Fixture) ReportDate() string {
	return "2022-04-11"
}

func (fixture *Fixture) CaAddress1() *easypost.Address {
	return readFixtureData().Addresses["ca_address_1"]
}

func (fixture *Fixture) CaAddress2() *easypost.Address {
	return readFixtureData().Addresses["ca_address_2"]
}

func (fixture *Fixture) IncorrectAddress() *easypost.Address {
	return readFixtureData().Addresses["incorrect"]
}

func (fixture *Fixture) BasicParcel() *easypost.Parcel {
	return readFixtureData().Parcels["basic"]
}

func (fixture *Fixture) BasicCustomsItem() *easypost.CustomsItem {
	customsItem := readFixtureData().CustomsItems["basic"]

	// Json unmarshalling doesn't handle float64 well, need to manually set the value
	customsItem.Value = 23.25

	return customsItem
}

func (fixture *Fixture) BasicCustomsInfo() *easypost.CustomsInfo {
	customsInfo := readFixtureData().CustomsInfos["basic"]

	// Json unmarshalling doesn't handle float64 well, need to manually set the value
	for _, customsItem := range customsInfo.CustomsItems {
		customsItem.Value = 23.25
	}
	return customsInfo
}

func (fixture *Fixture) TaxIdentifier() *easypost.TaxIdentifier {
	return readFixtureData().TaxIdentifiers["basic"]
}

func (fixture *Fixture) BasicShipment() *easypost.Shipment {
	return readFixtureData().Shipments["basic_domestic"]
}

func (fixture *Fixture) FullShipment() *easypost.Shipment {
	return readFixtureData().Shipments["full"]
}

func (fixture *Fixture) OneCallBuyShipment() *easypost.Shipment {
	return &easypost.Shipment{
		ToAddress:         fixture.CaAddress1(),
		FromAddress:       fixture.CaAddress2(),
		Parcel:            fixture.BasicParcel(),
		Service:           fixture.USPSService(),
		CarrierAccountIDs: []string{fixture.USPSCarrierAccountID()},
		Carrier:           fixture.USPS(),
	}
}

// This fixture will require you to add a `shipment` key with a Shipment object from a test.
// If you need to re-record cassettes, increment the date below and ensure it is one day in the future,
// USPS only does "next-day" pickups including Saturday but not Sunday or Holidays.
func (fixture *Fixture) BasicPickup() *easypost.Pickup {
	pickupDate := easypost.NewDateTime(2024, time.August, 14, 0, 0, 0, 0, time.UTC)

	pickupData := readFixtureData().Pickups["basic"]
	pickupData.MinDatetime = &pickupDate
	pickupData.MaxDatetime = &pickupDate

	return pickupData
}

func (fixture *Fixture) BasicCarrierAccount() *easypost.CarrierAccount {
	return readFixtureData().CarrierAccounts["basic"]
}

// This fixture will require you to add a `tracking_code` key with a tracking code from a shipment
func (fixture *Fixture) BasicInsurance() *easypost.Insurance {
	return readFixtureData().Insurances["basic"]
}

// This fixture will require you to add a `tracking_code` key with a tracking code from a shipment and an `amount` key with a float value
func (fixture *Fixture) BasicClaim() *easypost.CreateClaimParameters {
	return readFixtureData().Claims["basic"]
}

func (fixture *Fixture) BasicOrder() *easypost.Order {
	return readFixtureData().Orders["basic"]
}

func (fixture *Fixture) EventBody() []byte {
	currentDir, _ := os.Getwd()
	parentDir := filepath.Dir(currentDir)
	filePath := fmt.Sprintf("%s%s", parentDir, "/examples/official/fixtures/event-body.json")

	/* #nosec */
	data, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error opening fixture file:", err)
		os.Exit(1)
	}

	defer func() { _ = data.Close() }()

	scanner := bufio.NewScanner(data)
	var eventBody []byte

	for scanner.Scan() {
		eventBody = []byte(scanner.Text())
	}

	return eventBody
}

func (fixture *Fixture) WebhookHmacSignature() string {
	return readFixtureData().Webhooks["hmac_signature"].(string)
}

func (fixture *Fixture) WebhookSecret() string {
	return readFixtureData().Webhooks["secret"].(string)
}

func (fixture *Fixture) WebhookUrl() string {
	return readFixtureData().Webhooks["url"].(string)
}

func (fixture *Fixture) WebhookCustomHeaders() []easypost.WebhookCustomHeader {
	var header easypost.WebhookCustomHeader

	inputCustomHeaders := readFixtureData().Webhooks["custom_headers"]
	customHeaders := inputCustomHeaders.([]interface{})
	headerMap := customHeaders[0].(map[string]interface{})

	header.Name = headerMap["name"].(string)
	header.Value = headerMap["value"].(string)
	webhookCustomHeaders := []easypost.WebhookCustomHeader{header}

	return webhookCustomHeaders
}

func (fixture *Fixture) RmaFormOptions() map[string]interface{} {
	return readFixtureData().FormOptions["rma"]
}

func (fixture *Fixture) ReferralUser() *easypost.UserOptions {
	return readFixtureData().Users["referral"]
}

func (fixture *Fixture) TestCreditCard() *easypost.CreditCardOptions {
	return readFixtureData().CreditCards["test"]
}

func (fixture *Fixture) PlannedShipDate() string {
	return "2024-08-14"
}

func (fixture *Fixture) DesiredDeliveryDate() string {
	return "2024-08-14"
}

func (fixture *Fixture) BillingData() billingFixture {
	return readFixtureData().Billing
}

func (fixture *Fixture) LumaRulesetName() string {
	return readFixtureData().Luma["ruleset_name"].(string)
}

func (fixture *Fixture) LumaPlannedShipDate() string {
	return "2025-06-16"
}
