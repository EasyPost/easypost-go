package easypost_test

import (
	"os"
	"time"

	"github.com/EasyPost/easypost-go/v2"
)

type Fixture struct {
}

// We keep the page_size of retrieving `all` records small so cassettes stay small
func (fixture *Fixture) pageSize() int {
	return 5
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
	return "USPS"
}

func (fixture *Fixture) USPSService() string {
	return "First"
}

func (fixture *Fixture) PickupService() string {
	return "NextDay"
}

func (fixture *Fixture) ReportType() string {
	return "shipment"
}

// If you need to re-record cassettes, increment this date by 1
func (fixture *Fixture) ReportDate() string {
	return "2022-04-11"
}

func (fixture *Fixture) WebhookUrl() string {
	return "http://example.com"
}

func (fixture *Fixture) BasicAddress() *easypost.Address {
	return &easypost.Address{
		Name:    "Jack Sparrow",
		Company: "EasyPost",
		Street1: "388 Townsend St",
		Street2: "Apt 20",
		City:    "San Francisco",
		State:   "CA",
		Zip:     "94107",
		Phone:   "5555555555",
	}
}

func (fixture *Fixture) IncorrectAddressToVerify() *easypost.Address {
	return &easypost.Address{
		Street1: "417 montgomery street",
		Street2: "FL 5",
		City:    "San Francisco",
		State:   "CA",
		Zip:     "94104",
		Country: "US",
		Company: "EasyPost",
		Phone:   "415-123-4567",
	}

}

func (fixture *Fixture) PickupAddress() *easypost.Pickup {
	return &easypost.Pickup{
		Address: &easypost.Address{
			Name:    "Dr. Steve Brule",
			Street1: "179 N Harbor Dr",
			City:    "Redondo Beach",
			State:   "CA",
			Zip:     "90277",
			Country: "US",
			Phone:   "3331114444"},
	}
}

func (fixture *Fixture) BasicParcel() *easypost.Parcel {
	return &easypost.Parcel{
		Length: 10,
		Width:  8,
		Height: 4,
		Weight: 15.4,
	}
}

func (fixture *Fixture) BasicCustomsItem() *easypost.CustomsItem {
	return &easypost.CustomsItem{
		Description:    "Sweet shirts",
		Quantity:       2,
		Weight:         11,
		Value:          23,
		HSTariffNumber: "654321",
		OriginCountry:  "US",
	}
}

func (fixture *Fixture) BasicCustomsInfo() *easypost.CustomsInfo {
	return &easypost.CustomsInfo{
		EELPFC:              "NOEEI 30.37(a)",
		CustomsCertify:      true,
		CustomsSigner:       "Steve Brule",
		ContentsType:        "merchandise",
		ContentsExplanation: "",
		RestrictionType:     "none",
		NonDeliveryOption:   "return",
		CustomsItems:        []*easypost.CustomsItem{fixture.BasicCustomsItem()},
	}
}

func (fixture *Fixture) TaxIdentifier() *easypost.TaxIdentifier {
	return &easypost.TaxIdentifier{
		Entity:         "SENDER",
		TaxIdType:      "IOSS",
		TaxId:          "12345",
		IssuingCountry: "GB",
	}
}

func (fixture *Fixture) BasicShipment() *easypost.Shipment {
	return &easypost.Shipment{
		ToAddress:   fixture.BasicAddress(),
		FromAddress: fixture.BasicAddress(),
		Parcel:      fixture.BasicParcel(),
	}
}

func (fixture *Fixture) FullShipment() *easypost.Shipment {
	return &easypost.Shipment{
		ToAddress:   fixture.BasicAddress(),
		FromAddress: fixture.BasicAddress(),
		Parcel:      fixture.BasicParcel(),
		CustomsInfo: fixture.BasicCustomsInfo(),
		Options: &easypost.ShipmentOptions{
			LabelFormat:   "PNG", // Must be PNG so we can convert to ZPL later
			InvoiceNumber: "123",
		},
		Reference: "123",
	}
}

func (fixture *Fixture) OneCallBuyShipment() *easypost.Shipment {
	return &easypost.Shipment{
		ToAddress:         fixture.BasicAddress(),
		FromAddress:       fixture.BasicAddress(),
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
	pickupDate := time.Date(2022, time.June, 15, 0, 0, 0, 0, time.UTC)

	return &easypost.Pickup{
		Address:      fixture.BasicAddress(),
		MinDatetime:  &pickupDate,
		MaxDatetime:  &pickupDate,
		Instructions: "Pickup at front door",
	}
}

func (fixture *Fixture) BasicCarrierAccount() *easypost.CarrierAccount {
	return &easypost.CarrierAccount{
		Type: "UpsAccount",
		Credentials: map[string]string{
			"account_number":        "A1A1A1",
			"user_id":               "USERID",
			"password":              "PASSWORD",
			"access_license_number": "ALN",
		},
	}
}

// This fixture will require you to add a `tracking_code` key with a tracking code from a shipment
func (fixture *Fixture) BasicInsurance() *easypost.Insurance {
	return &easypost.Insurance{
		ToAddress:   fixture.BasicAddress(),
		FromAddress: fixture.BasicAddress(),
		Carrier:     fixture.USPS(),
		Amount:      "100",
	}
}

func (fixture *Fixture) BasicOrder() *easypost.Order {
	return &easypost.Order{
		ToAddress:   fixture.BasicAddress(),
		FromAddress: fixture.BasicAddress(),
		Shipments:   []*easypost.Shipment{fixture.BasicShipment()},
	}
}
