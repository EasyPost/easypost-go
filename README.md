# EasyPost Go Client Library

[![CI](https://github.com/EasyPost/easypost-go/workflows/CI/badge.svg)](https://github.com/EasyPost/easypost-go/actions?query=workflow%3ACI)
[![GitHub version](https://badge.fury.io/gh/EasyPost%2Feasypost-go.svg)](https://badge.fury.io/gh/EasyPost%2Feasypost-go)
[![GoDoc](https://godoc.org/github.com/EasyPost/easypost-go?status.svg)](https://pkg.go.dev/github.com/EasyPost/easypost-go)

EasyPost, the simple shipping solution. You can sign up for an account at https://easypost.com.

## Install

```bash
go get -u github.com/EasyPost/easypost-go/v2
```

```go
// Import the library
import (
    "github.com/EasyPost/easypost-go/v2"
)
```

## Usage

A simple create & buy shipment example:

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/EasyPost/easypost-go/v2"
)

func main() {
	apiKey := os.Getenv("EASYPOST_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "missing API key")
		os.Exit(1)
		return
	}
	client := easypost.New(apiKey)

    fromAddress := &easypost.Address{
		Company: "EasyPost",
		Street1: "One Montgomery St",
		Street2: "Ste 400",
		City:    "San Francisco",
		State:   "CA",
		Zip:     "94104",
		Phone:   "415-555-1212",
	}

	toAddress := &easypost.Address{
		Name:    "Bugs Bunny",
		Street1: "4000 Warner Blvd",
		City:    "Burbank",
		State:   "CA",
		Zip:     "91522",
		Phone:   "8018982020",
	}

	parcel := &easypost.Parcel{
		Length: 10.2,
		Width:  7.8,
		Height: 4.3,
		Weight: 21.2,
	}

	shipment, err := client.CreateShipment(
		&easypost.Shipment{
			FromAddress: fromAddress,

			ToAddress:   toAddress,
			Parcel:      parcel,
		},
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating shipment:", err)
		os.Exit(1)
		return
	}

	lowestRate, err := client.LowestRate(shipment)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error getting lowest rate:", err)
		os.Exit(1)
		return
	}

    shipment, err = client.BuyShipment(shipment.ID, &easypost.Rate{ID: lowestRate.ID}, "")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error buying shipment:", err)
		os.Exit(1)
	}

	prettyJSON, err := json.MarshalIndent(shipment, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating JSON:", err)
		os.Exit(1)
		return
	}

    fmt.Println(string(prettyJSON))
}
```

## Documentation

See the above example, or view the [GoDoc](https://pkg.go.dev/github.com/EasyPost/easypost-go) for more details.

## Development

```bash
# Run tests
EASYPOST_TEST_API_KEY=123... EASYPOST_PROD_API_KEY=123... make test

# Generate test coverage
EASYPOST_TEST_API_KEY=123... EASYPOST_PROD_API_KEY=123... make coverage

# Lint project (requires golangci-lint to be installed)
make lint
```

### Testing

The test suite in this project was specifically built to produce consistent results on every run, regardless of when they run or who is running them. This project uses [VCR](https://github.com/dnaeon/go-vcr) to record and replay HTTP requests and responses via "cassettes". When the suite is run, the HTTP requests and responses for each test function will be saved to a cassette if they do not exist already and replayed from this saved file if they do, which saves the need to make live API calls on every test run.

If you make an addition to this project, the request/response will get recorded automatically for you. When making changes to this project, you'll need to re-record the associated cassette to force a new live API call for that test which will then record the request/response used on the next run.

The test suite has been populated with various helpful fixtures that are available for use, each completely independent from a particular user **with the exception of the USPS carrier account ID** which has a fallback value to our internal testing user's ID. If you are a non-EasyPost employee and are re-recording cassettes, you may need to provide the `USPS_CARRIER_ACCOUNT_ID` environment variable with the ID associated with your USPS account (which will be associated with your API keys in use) for tests that use this fixture.

**Note on dates:** Some fixtures use hard-coded dates that may need to be incremented if cassettes get re-recorded (such as reports or pickups).
