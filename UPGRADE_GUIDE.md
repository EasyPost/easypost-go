# Upgrade Guide

Use the following guide to assist in the upgrade process of the `easypost-go` library between major versions.

## Upgrading from 1.x to 2.0

### 2.0 High Impact Changes

* [Importing the Library](#20-importing-the-library)
* [Updating Dependencies](#20-updating-dependencies)

### 2.0 Medium Impact Changes

* [Default Timeouts for HTTP Requests](#20-default-timeouts-for-http-requests)
* [Removal of `GetShipmentRates()` Shipment Method](#20-removal-of-getshipmentrates-shipment-method)

## 2.0 Importing the Library

Likelihood of Impact: High

With the transition to `v2`, this library must now be imported as follows:

```go
import (
    "github.com/EasyPost/easypost-go/v2"
)
```

## 2.0 Updating Dependencies

Likelihood of Impact: High

**Go 1.15 Required**

easypost-go now requires Go 1.15 or greater.

**Dependencies**

All dependencies had minor version bumps.

## 2.0 Default Timeouts for HTTP Requests

Likelihood of Impact: Medium

Default timeouts for all HTTP requests are now set to 60 seconds. If you require longer timeouts, you can set them by overriding the defaults:

```go
client := &easypost.Client{
    Timeout: 60000,
}
```

## 2.0 Removal of GetShipmentRates() Shipment Method

Likelihood of Impact: Medium

The HTTP method used for the `GetShipmentRates` endpoint at the API level has changed from `POST` to `GET` and will only retrieve rates for a shipment instead of regenerating them. A new `/rerate` endpoint has been introduced to replace this functionality; In this library, you can now call the `RerateShipment` method to regenerate rates. Due to the logic change, the `GetShipmentRates` method has been removed since a Shipment inherently already has rates associated.
