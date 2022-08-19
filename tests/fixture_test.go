package easypost_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/EasyPost/easypost-go/v2"
)

type Fixture struct {
	Addresses       map[string]*easypost.Address        `json:"addresses,omitempty"`
	CarrierAccounts map[string]*easypost.CarrierAccount `json:"carrier_accounts,omitempty"`
	CarrierStrings  map[string]string                   `json:"carrier_strings,omitempty"`
	CustomsInfos    map[string]*easypost.CustomsInfo    `json:"customs_infos,omitempty"`
	CustomsItems    map[string]*easypost.CustomsItem    `json:"customs_items,omitempty"`
	FormOptions     map[string]map[string]interface{}   `json:"form_options,omitempty"`
	Insurances      map[string]*easypost.Insurance      `json:"insurances,omitempty"`
	Orders          map[string]*easypost.Order          `json:"orders,omitempty"`
	PageSizes       map[string]int                      `json:"page_sizes,omitempty"`
	Parcels         map[string]*easypost.Parcel         `json:"parcels,omitempty"`
	Pickups         map[string]*easypost.Pickup         `json:"pickups,omitempty"`
	ReportTypes     map[string]string                   `json:"report_types,omitempty"`
	ServiceNames    map[string]map[string]string        `json:"service_names,omitempty"`
	Shipments       map[string]*easypost.Shipment       `json:"shipments,omitempty"`
	TaxIdentifiers  map[string]*easypost.TaxIdentifier  `json:"tax_identifiers,omitempty"`
	WebhookURL      string                              `json:"webhook_url,omitempty"`
}

// Replace the min_datetime and max_datetime to RFC3339 format
func replaceJsonFileDate(jsonFile string) error {
	byteValue, _ := ioutil.ReadFile(jsonFile)

	var result map[string]interface{}
	_ = json.Unmarshal(byteValue, &result)

	result["pickups"].(map[string]interface{})["basic"].(map[string]interface{})["min_datetime"] = "2022-08-01T00:00:00.00Z"
	result["pickups"].(map[string]interface{})["basic"].(map[string]interface{})["max_datetime"] = "2022-08-02T00:00:00.00Z"

	byteValue, _ = json.Marshal(result)

	return ioutil.WriteFile(jsonFile, byteValue, 0644)
}

// Reads fixture data from the fixtures JSON file
func readFixtureData() Fixture {
	currentDir, _ := os.Getwd()
	parentDir := filepath.Dir(currentDir)
	filePath := fmt.Sprintf("%s%s", parentDir, "/examples/official/fixtures/client-library-fixtures.json")
	_ = replaceJsonFileDate(filePath)
	data, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error opening fixture file:", err)
		os.Exit(1)
	}

	defer data.Close()

	byteData, _ := ioutil.ReadAll(data)

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

func (fixture *Fixture) WebhookUrl() string {
	fixtureData := readFixtureData()
	fmt.Println(fixtureData)
	return readFixtureData().WebhookURL
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
	return readFixtureData().CustomsItems["basic"]
}

func (fixture *Fixture) BasicCustomsInfo() *easypost.CustomsInfo {
	return readFixtureData().CustomsInfos["basic"]
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
	pickupDate := time.Date(2022, time.August, 22, 0, 0, 0, 0, time.UTC)

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

func (fixture *Fixture) BasicOrder() *easypost.Order {
	return readFixtureData().Orders["basic"]
}

// TODO: FIX THIS
func (fixture *Fixture) EventBody() []byte {
	// Test data MUST be a minified string to work correctly
	data := `{"result":{"id":"batch_123...","object":"Batch","mode":"test","state":"created","num_shipments":0,"reference":null,"created_at":"2022-07-26T17:22:32Z","updated_at":"2022-07-26T17:22:32Z","scan_form":null,"shipments":[],"status":{"created":0,"queued_for_purchase":0,"creation_failed":0,"postage_purchased":0,"postage_purchase_failed":0},"pickup":null,"label_url":null},"description":"batch.created","mode":"test","previous_attributes":null,"completed_urls":null,"user_id":"user_123...","status":"pending","object":"Event","id":"evt_123..."}`

	return []byte(data)
}

func (fixture *Fixture) RmaFormOptions() map[string]interface{} {
	return readFixtureData().FormOptions["rma"]
}
