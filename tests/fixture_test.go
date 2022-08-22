package easypost_test

import (
	"bufio"
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

// Reads fixture data from the fixtures JSON file
func readFixtureData() Fixture {
	currentDir, _ := os.Getwd()
	parentDir := filepath.Dir(currentDir)
	filePath := fmt.Sprintf("%s%s", parentDir, "/examples/official/fixtures/client-library-fixtures.json")

	data, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error opening fixture file:", err)
		os.Exit(1)
	}

	defer data.Close()

	byteData, _ := ioutil.ReadAll(data)

	var fixtures Fixture
	json.Unmarshal([]byte(byteData), &fixtures)

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

func (fixture *Fixture) EventBody() []byte {
	currentDir, _ := os.Getwd()
	parentDir := filepath.Dir(currentDir)
	filePath := fmt.Sprintf("%s%s", parentDir, "/examples/official/fixtures/event-body.json")

	data, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error opening fixture file:", err)
		os.Exit(1)
	}

	defer data.Close()

	scanner := bufio.NewScanner(data)
	var eventBody []byte

	for scanner.Scan() {
		eventBody = []byte(scanner.Text())
	}

	return eventBody
}

func (fixture *Fixture) RmaFormOptions() map[string]interface{} {
	return readFixtureData().FormOptions["rma"]
}
