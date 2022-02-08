package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/EasyPost/easypost-go"
)

func main() {
	apiKey := os.Getenv("EASYPOST_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "missing API key")
		os.Exit(1)
		return
	}
	client := easypost.New(apiKey)

	toAddress := &easypost.Address{
		Name:    "Bugs Bunny",
		Street1: "4000 Warner Blvd",
		City:    "Burbank",
		State:   "CA",
		Zip:     "91522",
		Phone:   "8018982020",
	}

	fromAddress := &easypost.Address{
		Company: "EasyPost",
		Street1: "One Montgomery St",
		Street2: "Ste 400",
		City:    "San Francisco",
		State:   "CA",
		Zip:     "94104",
		Phone:   "415-555-1212",
	}

	parcel := &easypost.Parcel{
		Length: 10.2,
		Width:  7.8,
		Height: 4.3,
		Weight: 21.2,
	}

	customsItem := &easypost.CustomsItem{
		Description:    "EasyPost t-shirts",
		HSTariffNumber: "123456",
		OriginCountry:  "US",
		Quantity:       2,
		Value:          96.27,
		Weight:         21.1,
	}

	customsInfo := &easypost.CustomsInfo{
		CustomsCertify:    true,
		CustomsSigner:     "Wile E. Coyote",
		ContentsType:      "gift",
		EELPFC:            "NOEEI 30.37(a)",
		NonDeliveryOption: "return",
		RestrictionType:   "none",
		CustomsItems:      []*easypost.CustomsItem{customsItem},
	}

	// Create the shipment
	shipment, err := client.CreateShipment(
		&easypost.Shipment{
			ToAddress:   toAddress,
			FromAddress: fromAddress,
			Parcel:      parcel,
			CustomsInfo: customsInfo,
		},
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating shipment:", err)
		os.Exit(1)
		return
	}

	// Get the lowest rate of the shipment
	rate, err := client.LowestRate(shipment)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error getting lowest rate:", err)
		os.Exit(1)
		return
	}

	prettyJSON, err := json.MarshalIndent(rate, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating JSON:", err)
		os.Exit(1)
		return
	}
	fmt.Println(string(prettyJSON))
}
