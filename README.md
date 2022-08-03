# EasyPost Go Client Library

[![CI](https://github.com/EasyPost/easypost-go/workflows/CI/badge.svg)](https://github.com/EasyPost/easypost-go/actions?query=workflow%3ACI)
[![GitHub version](https://badge.fury.io/gh/EasyPost%2Feasypost-go.svg)](https://badge.fury.io/gh/EasyPost%2Feasypost-go)
[![GoDoc](https://godoc.org/github.com/EasyPost/easypost-go?status.svg)](https://pkg.go.dev/github.com/EasyPost/easypost-go)

EasyPost, the simple shipping solution. You can sign up for an account at <https://easypost.com>.

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

API Documentation can be found at: <https://easypost.com/docs/api>. Alternatively, you can view the [GoDoc](https://pkg.go.dev/github.com/EasyPost/easypost-go).

Upgrading major versions of this project? Refer to the [Upgrade Guide](UPGRADE_GUIDE.md).

## Development

```bash
# Install dependencies
make install

# Lint project (requires `golangci-lint` to be installed - not included)
make lint

# Run tests
EASYPOST_TEST_API_KEY=123... EASYPOST_PROD_API_KEY=123... make test

# Generate test coverage
EASYPOST_TEST_API_KEY=123... EASYPOST_PROD_API_KEY=123... make coverage

# Run security analysis on the project (requires `gosec` to be installed - not included)
make scan
```

### Testing

The test suite in this project was specifically built to produce consistent results on every run, regardless of when they run or who is running them. This project uses [VCR](https://github.com/dnaeon/go-vcr) to record and replay HTTP requests and responses via "cassettes". When the suite is run, the HTTP requests and responses for each test function will be saved to a cassette if they do not exist already and replayed from this saved file if they do, which saves the need to make live API calls on every test run.

**Sensitive Data:** We've made every attempt to include scrubbers for sensitive data when recording cassettes so that PII or sensitive info does not persist in version control; however, please ensure when recording or re-recording cassettes that prior to committing your changes, no PII or sensitive information gets persisted by inspecting the cassette.

**Making Changes:** If you make an addition to this project, the request/response will get recorded automatically for you. When making changes to this project, you'll need to re-record the associated cassette to force a new live API call for that test which will then record the request/response used on the next run.

**Test Data:** The test suite has been populated with various helpful fixtures that are available for use, each completely independent of a particular user **with the exception of the USPS carrier account ID** (see [Unit Test API Keys](#unit-test-api-keys) for more information) which has a fallback value of our internal testing user's ID. Some fixtures use hard-coded dates that may need to be incremented if cassettes get re-recorded (such as reports or pickups).

#### Unit Test API Keys

The following are required on every test run:

- `EASYPOST_TEST_API_KEY`
- `EASYPOST_PROD_API_KEY`

The following are required when you need to re-record cassettes for applicable tests (fallback values are used otherwise):

- `USPS_CARRIER_ACCOUNT_ID` (eg: one-call buying a shipment for non-EasyPost employees)
- `REFERRAL_USER_PROD_API_KEY` (eg: adding a credit card to a referral user)

Some tests may require a user with a particular set of enabled features such as a `Partner` user when creating referrals. We have attempted to call out these functions in their respective docstrings.
