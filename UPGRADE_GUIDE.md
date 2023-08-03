# Upgrade Guide

Use the following guide to assist in the upgrade process of the `easypost-go` library between major versions.

* [Upgrading from 2.x to 3.0](#upgrading-from-2x-to-30)
* [Upgrading from 1.x to 2.0](#upgrading-from-1x-to-20)

## Upgrading from 2.x to 3.0

### 3.0 High Impact Changes

* [Importing the Library](#30-importing-the-library)
* [Updating Dependencies](#30-updating-dependencies)
* [New `DateTime` Type](#30-new-datetime-type)
* [Deprecated Methods and Structs Removed](#30-deprecated-methods-and-structs-removed)
* [Error Types](#30-error-types)

### 3.0 Medium Impact Changes

* [`TrackingCodes` Parameter Renamed](#30-trackingcodes-parameter-renamed)
* [Add/Remove Shipments to Batch Parameter Change](#30-addremove-shipments-to-batch-parameter-change)

## 3.0 Importing the Library

*Likelihood of Impact: **High***

With the transition to `v3`, this library must now be imported as follows:

```go
import (
    "github.com/EasyPost/easypost-go/v3"
)
```

## 3.0 Updating Dependencies

*Likelihood of Impact: **High***

**Go 1.16 Required**

easypost-go now requires Go 1.16 or greater.

**Dependencies**

Some dependencies had minor version bumps.

## 3.0 New `DateTime` Type

*Likelihood of Impact: **High***

All uses of `time.Time` in the library (including in parameter and API response structs) have been replaced with a new `DateTime` type.

This type is a wrapper around `time.Time` that handles unexpected datetime formats returned by the API.

`DateTime` should be a drop-in replacement for `time.Time` in most cases:

```go
// Before
var t time.Time
t := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

// After
var t easypost.DateTime
t := easypost.DateTime(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))
// or
t := easypost.DateTimeFromTime(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))
// or
t := easypost.NewDateTime(2021, 1, 1, 0, 0, 0, 0, time.UTC)
```

`DateTime` objects can be converted to `time.Time` objects using the `AsTime()` method:

```go
var t easypost.DateTime
t := easypost.NewDateTime(2021, 1, 1, 0, 0, 0, 0, time.UTC)

t.AsTime() // time.Time
```

## 3.0 Deprecated Methods and Structs Removed

*Likelihood of Impact: **High***

All methods and structs marked as deprecated in `v2` have been removed.

The following is a list of items that have been removed, along with their suggested replacements:

| Deprecated                                | Replacement                           |
|-------------------------------------------|---------------------------------------|
| `ListReportOptions` struct                | `ListOptions` struct                  |
| `BetaCarrierMetadata` methods and structs | `CarrierMetadata` methods and structs |
| `LowestRate` shipment method              | `LowestShipmentRate` method           |
| `CreateWebhook` method                    | `CreateWebhookWithDetails` method     |
| `EnableWebhook` method                    | `UpdateWebhook` method                |

## 3.0 Error Types

*Likelihood of Impact: **High***

Error handling has been overhauled and a number of specific exception types have been introduced.

All exceptions inherit from `easypost.LibraryError`.

Subclasses of `easypost.LibraryError` are grouped into two categories:

- `easypost.APIError` for errors returned by the API. Users will need to anticipate one or multiple of the following errors that inherit this class:
  - `easypost.BadRequestError` - thrown when the API returns a 400 status code
  - `easypost.UnauthorizedError` - thrown when the API returns a 401 status code
  - `easypost.PaymentError` - thrown when the API returns a 402 status code
  - `easypost.ForbiddenError` - thrown when the API returns a 403 status code
  - `easypost.NotFoundError` - thrown when the API returns a 404 status code
  - `easypost.MethodNotAllowedError` - thrown when the API returns a 405 status code
  - `easypost.ProxyError` - thrown when the API returns a 407 status code
  - `easypost.TimeoutError` - thrown when the API returns a 408 status code
  - `easypost.InvalidRequestError` - thrown when the API returns a 422 status code
  - `easypost.RateLimitError` - thrown when the API returns a 429 status code
  - `easypost.InternalServerError` - thrown when the API returns a 500 status code
  - `easypost.GatewayTimeoutError` - thrown when the API returns a 502 or 504 status code
  - `easypost.ServiceUnavailableError` - thrown when the API returns a 503 status code
  - `easypost.RedirectError` - thrown when the API returns a 3xx status code
  - `easypost.RetryError` - thrown when the API returns a 1xx status code
  - `easypost.UnknownError` - thrown when the API returns an unexpected status code
  - `easypost.SSLError` - thrown when there is an SSL error
- `easypost.LocalError` for errors raised by local operations, such as validation, JSON de/serialization or list filtering.
  - `easypost.EndOfPaginationError` - thrown when the end of a paginated list is reached
  - `easypost.FilteringError` - thrown when there is an issue running a filtering operation
  - `easypost.InvalidObjectError` - thrown when an object is invalid
  - `easypost.MissingPropertyError` - thrown when a required property is missing
  - `easypost.MissingWebhookSignatureError` - thrown when a webhook signature is missing
  - `easypost.MismatchWebhookSignatureError` - thrown when a webhook signature does not match
  - `easypost.ExternalApiError` - thrown when there is an error with an external API (e.g. Stripe)

Any error type can be pretty-printed as a string using the `Error()` method:

```go
_, err := client.CreateAddress(&easypost.Address{
  Name: "Invalid Address",
})

fmt.Println(err.Error()) // "PARAMETER.REQUIRED Missing required parameter"
```

Any `easypost.APIError`-derived error type will have the following properties:

- `Code` - An error code string referring to the specific error that occurred server-side
- `StatusCode` - The HTTP status code returned by the API
- `Message` - A human-readable message describing the error
- `Errors` - A slice of `easypost.Error` structs, representing specific details of server-side errors encountered

Any `easypost.LocalError`-derived error type will have the following properties:

- `Message` - A human-readable message describing the error

Common strings and templates used to construct error messages are available as constants (e.g. `easypost.NoRatesMatchingFilters`). These can be used to perform regex-based evaluation of error messages.

**Note**: An `easypost.LibraryError` and its subclasses represent errors that occur within this library. This is different from `easypost.Error`, which is a struct representing a server-side error returned by a failed API call.

## 3.0 TrackingCodes Parameter Renamed

*Likelihood of Impact: **Medium***

The `TrackingCodes` parameter for the `ListTrackers` method has been renamed to `TrackingCode` to match the API.

Instead of passing a slice of strings, you should now pass a single string:

```go
// Before
trackers, err := client.ListTrackers(&easypost.ListTrackersOptions{
	TrackingCodes: []string{"EZ1000000001", "EZ1000000002"},
})

// After
trackers, err := client.ListTrackers(&easypost.ListTrackersOptions{
	TrackingCode: "EZ1000000001",
})
```

## 3.0 Add/Remove Shipments to Batch Parameter Change

*Likelihood of Impact: **Medium***

The `AddShipmentsToBatch` and `RemoveShipmentsFromBatch` methods now explicitly accept `Shipment` structs instead of generic `interface{}` types

These methods will no longer accept solely IDs; users will need to provide whole `Shipment` structs

```go
// Before
batch, err := client.AddShipmentsToBatch("batch_123", shipment1.ID, shipment2.ID)

// After
batch, err := client.AddShipmentsToBatch("batch_123", shipment1, shipment2)
```

## Upgrading from 1.x to 2.0

### 2.0 High Impact Changes

* [Importing the Library](#20-importing-the-library)
* [Updating Dependencies](#20-updating-dependencies)

### 2.0 Medium Impact Changes

* [Default Timeouts for HTTP Requests](#20-default-timeouts-for-http-requests)
* [Removal of `GetShipmentRates()` Shipment Method](#20-removal-of-getshipmentrates-shipment-method)

## 2.0 Importing the Library

*Likelihood of Impact: **High***

With the transition to `v2`, this library must now be imported as follows:

```go
import (
    "github.com/EasyPost/easypost-go/v2"
)
```

## 2.0 Updating Dependencies

*Likelihood of Impact: **High***

**Go 1.15 Required**

easypost-go now requires Go 1.15 or greater.

**Dependencies**

All dependencies had minor version bumps.

## 2.0 Default Timeouts for HTTP Requests

*Likelihood of Impact: **Medium***

Default timeouts for all HTTP requests are now set to 60 seconds. If you require longer timeouts, you can set them by overriding the defaults:

```go
client := &easypost.Client{
    Timeout: 60000,
}
```

## 2.0 Removal of GetShipmentRates() Shipment Method

*Likelihood of Impact: **Medium***

The HTTP method used for the `GetShipmentRates` endpoint at the API level has changed from `POST` to `GET` and will only retrieve rates for a shipment instead of regenerating them. A new `/rerate` endpoint has been introduced to replace this functionality; In this library, you can now call the `RerateShipment` method to regenerate rates. Due to the logic change, the `GetShipmentRates` method has been removed since a Shipment inherently already has rates associated.
