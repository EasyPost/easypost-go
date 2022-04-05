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
# Run tests (must be run from the `tests` directory)
go test

# Generate test coverage
go test ./tests -coverprofile=covprofile -coverpkg=./... && go tool cover -html=covprofile

# Lint project (requires golangci-lint to be installed)
golangci-lint run
```
