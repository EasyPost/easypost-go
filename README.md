# EasyPost Go Client Library

[![CI](https://github.com/EasyPost/easypost-go/workflows/CI/badge.svg)](https://github.com/EasyPost/easypost-go/actions?query=workflow%3ACI)
[![Coverage Status](https://coveralls.io/repos/github/EasyPost/easypost-go/badge.svg?branch=master)](https://coveralls.io/github/EasyPost/easypost-go?branch=master)
[![GitHub version](https://badge.fury.io/gh/EasyPost%2Feasypost-go.svg)](https://badge.fury.io/gh/EasyPost%2Feasypost-go)
[![GoDoc](https://godoc.org/github.com/EasyPost/easypost-go?status.svg)](https://pkg.go.dev/github.com/EasyPost/easypost-go)

EasyPost, the simple shipping solution. You can sign up for an account at <https://easypost.com>.

## Install

```bash
go get -u github.com/EasyPost/easypost-go/v3
```

```go
// Import the library
import (
    "github.com/EasyPost/easypost-go/v3"
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

    "github.com/EasyPost/easypost-go/v3"
)

func main() {
    apiKey := os.Getenv("EASYPOST_API_KEY")
    if apiKey == "" {
        fmt.Fprintln(os.Stderr, "missing API key")
        os.Exit(1)
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
    }

    lowestRate, err := client.LowestRate(shipment)
    if err != nil {
        fmt.Fprintln(os.Stderr, "error getting lowest rate:", err)
        os.Exit(1)
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
    }

    fmt.Println(string(prettyJSON))
}
```

## HTTP Hooks

Users can audit the HTTP requests and responses being made by the library by setting the `Hooks` property of a `Client` with a set of event subscriptions. Available subscriptions include:

- `RequestHookEventSubscriber` - Called before an HTTP request is made. A `RequestHookEvent` object, containing details about the request that will be sent to the server, is passed to the subscription's `RequestHookEventSubscriberCallback`
    - Modifying any data in the callback will NOT affect the actual request that is sent to the server.
- `ResponseHookEventSubscriber` - Called after an HTTP request is made. A `ResponseHookEvent` object, containing details about the response that was received from the server, is passed to the subscription's `ResponseHookEventSubscriberCallback`
    - Modifying any data in the callback will NOT affect the actual response that was received from the server, and will NOT affect the data deserialized into the library's models.

Users can interact with these details in their callbacks as they see fit (e.g. logging).

```go
func myRequestHookEventSubscriberCallback(ctx context.Context, event easypost.RequestHookEvent) error {
    // Interact with details about the request here
    fmt.Printf("Making HTTP call to %s\n", event.URL)
	return nil
}

func myResponseHookEventSubscriberCallback(ctx context.Context, event easypost.ResponseHookEvent) error {
    // Interact with details about the response here
    fmt.Printf("Received HTTP response with status code %d\n", event.StatusCode)
    return nil
}

requestSubscriber := easypost.RequestHookEventSubscriber{
    Callback: myRequestHookEventSubscriberCallback,
    HookEventSubscriber: easypost.HookEventSubscriber{
        ID: "my-request-hook",
    },
}

responseSubscriber := easypost.ResponseHookEventSubscriber{
    Callback: myResponseHookEventSubscriberCallback,
    HookEventSubscriber: easypost.HookEventSubscriber{
        ID: "my-response-hook",
    },
}

client := easypost.New("EASYPOST_API_KEY")

client.Hooks.AddRequestEventSubscriber(requestSubscriber)
client.Hooks.AddResponseEventSubscriber(responseSubscriber)
```

Users can unsubscribe from these events at any time by removing the subscription from the `Hooks` property of a client.

```go
client.Hooks.RemoveRequestEventSubscriber(requestSubscriber)
client.Hooks.RemoveResponseEventSubscriber(responseSubscriber)
```

## Documentation

API documentation can be found at: <https://easypost.com/docs/api>.

Library documentation can be found on the web at [GoDoc](https://pkg.go.dev/github.com/EasyPost/easypost-go/v3).

Upgrading major versions of this project? Refer to the [Upgrade Guide](UPGRADE_GUIDE.md).

## Support

New features and bug fixes are released on the latest major release of this library. If you are on an older major release of this library, we recommend upgrading to the most recent release to take advantage of new features, bug fixes, and security patches. Older versions of this library will continue to work and be available as long as the API version they are tied to remains active; however, they will not receive updates and are considered EOL.

For additional support, see our [org-wide support policy](https://github.com/EasyPost/.github/blob/main/SUPPORT.md).

## Development

```bash
# Install dependencies
make install

# Lint project (requires `golangci-lint` to be installed - not included)
make lint

# Run tests
EASYPOST_TEST_API_KEY=123... EASYPOST_PROD_API_KEY=123... make test
EASYPOST_TEST_API_KEY=123... EASYPOST_PROD_API_KEY=123... make coverage

# Run security analysis on the project (requires `gosec` to be installed - not included)
make scan

# Update submodules
make update-examples-submodule
```

### Testing

The test suite in this project was specifically built to produce consistent results on every run, regardless of when they run or who is running them. This project uses [VCR](https://github.com/dnaeon/go-vcr) to record and replay HTTP requests and responses via "cassettes". When the suite is run, the HTTP requests and responses for each test function will be saved to a cassette if they do not exist already and replayed from this saved file if they do, which saves the need to make live API calls on every test run. If you receive errors about a cassette expiring, delete and re-record the cassette to ensure the data is up-to-date.

**Sensitive Data:** We've made every attempt to include scrubbers for sensitive data when recording cassettes so that PII or sensitive info does not persist in version control; however, please ensure when recording or re-recording cassettes that prior to committing your changes, no PII or sensitive information gets persisted by inspecting the cassette.

**Making Changes:** If you make an addition to this project, the request/response will get recorded automatically for you. When making changes to this project, you'll need to re-record the associated cassette to force a new live API call for that test which will then record the request/response used on the next run.

**Test Data:** The test suite has been populated with various helpful fixtures that are available for use, each completely independent of a particular user **with the exception of the USPS carrier account ID** (see [Unit Test API Keys](#unit-test-api-keys) for more information) which has a fallback value of our internal testing user's ID. Some fixtures use hard-coded dates that may need to be incremented if cassettes get re-recorded (such as reports or pickups).

#### Unit Test API Keys

The following are required on every test run:

- `EASYPOST_TEST_API_KEY`
- `EASYPOST_PROD_API_KEY`

Some tests may require an EasyPost user with a particular set of enabled features such as a `Partner` user when creating referrals. We have attempted to call out these functions in their respective docstrings. The following are required when you need to re-record cassettes for applicable tests:

- `USPS_CARRIER_ACCOUNT_ID` (eg: one-call buying a shipment for non-EasyPost employees)
- `PARTNER_USER_PROD_API_KEY` (eg: creating a referral user)
- `REFERRAL_CUSTOMER_PROD_API_KEY` (eg: adding a credit card to a referral user)

#### Mocking

Some of our unit tests require HTTP calls that cannot be easily tested with live/recorded calls (e.g. HTTP calls that
trigger payments or interact with external APIs).

We have implemented a custom, lightweight HTTP mocking functionality in this library that allows us to mock HTTP calls
and responses.

The mock client is the same as a normal client, with a set of mock request-response pairs stored as a property.

At the time to make a real HTTP request, the mock client will instead check which mock request entry matches the queued
request (matching by HTTP method and a regex pattern for the URL), and will return the corresponding mock response (HTTP
status code and body).

**NOTE**: If a client is configured with any mock request entries, it will ONLY make mock requests. If it attempts to
make a request that does not match any of the configured mock requests, the request will fail and trigger an exception.

To use the mock client:

```golang
import "github.com/EasyPost/easypost-go/v3"

// create  a list of mock request-response pairs
mockRequests := []easypost.MockRequest{
    {
        MatchRule: easypost.MockRequestMatchRule{
            // HTTP method and regex pattern for the URL must both pass for the request to match
            Method:          "POST",
            UrlRegexPattern: "v2\\/bank_accounts\\/\\S*\\/charges$",
        },
        ResponseInfo: easypost.MockRequestResponseInfo{
            // HTTP status code and body to return when this request is matched
            StatusCode: 200,
            Body:       `{}`,
        },
    }
}

// create a new client with the mock requests
client := easypost.Client{
    APIKey:       "some_key",
    Client:       &http.Client{Transport: c.recorder},
    MockRequests: mockRequests,
}

// use the client as normal
```
