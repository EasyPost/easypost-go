# CHANGELOG

## v2.6.0 (2022-08-02)

- Adds Carbon Offset support
  - Adds the ability to create a shipment with carbon offset
  - Adds the ability to buy a shipment with carbon offset
  - Adds the ability to one-call-buy a shipment with carbon offset
  - Adds the ability to re-rate a shipment with carbon offset

## v2.5.0 (2022-07-18)

- Add support for generating shipment forms via `GenerateShipmentForm` function

## v2.4.0 (2022-07-11)

- Adds `RetrievePaymentMethods`, `FundWallet` and `DeletePaymentMethod` functions
- Adds `BillingType` in CarrierAccount and Rate structs
- Adds new lowest rate functions for Shipments, Orders and Pickups
- Adds new lowest smartrate functions for Shipments
- Adds OS details to the User-Agent header
- Adds support for webhook secrets
- Update methods now use `patch` instead of `put` behind the scenes to better match the API expectation and documentation. Update functions should still behave the same as before
- Enforces passing an API key on each request (the library will now fail fast instead of sending an impossible-to-service HTTP request)

## v2.3.0 (2022-04-13)

- Adds `Columns` and `AdditionalColumns` params to `Report` struct
- Adds `Declaration` attribute in `CustomsInfo` struct

## v2.2.0 (2022-03-08)

- Adds missing `CreateAndVerifyAddress` and `CreateAndVerifyAddressWithContext` functions

## v2.1.0 (2022-03-07)

- Adds `CreateRefund`, `CreateRefundWithContext`, `ListRefunds`, `ListRefundsWithContext`, `GetRefund`, and `GetRefundWithContext` in the Refund file
- Adds missing `ID` attribute to `Brand` struct
- Fixes a bug where the `AddShipmentsToBatch`, and `AddShipmentsToBatchWithContext` functions wouldn't allow a Shipment object to be passed by changing the params from `string` to `interface{}`
- Fixes a bug where the `RemoveShipmentsFromBatch`, and `RemoveShipmentsFromBatchWithContext` functions wouldn't allow a Shipment IDs to be passed by changing the params from `*Shipment` to `interface{}`
- Fixes the return type of `ListBatchesResult` from `batch` to `insurance`
- Fixes the HTTP method from `GET` to `POST` for `GetBatchLabels` and `GetBatchLabelsWithContext`
- Adds comprehensive test coverage

## v2.0.1 (2022-02-10)

- Corrects namespace for v2 release from `github.com/EasyPost/easypost-go` to `github.com/EasyPost/easypost-go/v2`

## v2.0.0 (2022-02-09)

**NOTE: Do not use this release, use v2.0.1 or later due to this release being incorrectly packaged.**

Upgrading major versions of this project? Refer to the [Upgrade Guide](UPGRADE_GUIDE.md).

- Bumps minimum Go version from `1.12` to `1.15`
- Bumps dependencies
- Lints the entire project
- adds a `client.RetrieveMe()` function to retrieve the authenticated user
- Adds support to one-call buy an order by passing a `Service` and `CarrierAccount`
- Adds support to update user brand via `client.UpdateBrand()`
- Adds support to create a list of trackers via `client.CreateTrackerList()`
- Adds support for getting the lowest rate of a shipment via `client.LowestRate()`
- Adds support to rerate a shipment via the `client.RerateShipment()` method
- Adds a default timeout of 60 seconds to requests. This can be overridden by setting the `Client.Timeout` option in milliseconds
- Fixed a spelling error for `origin_location` on the Tracker struct
- Removed `GetShipmentRates()` and `GetShipmentRatesWithContext()` methods since the shipment struct already has rates. If you need to get new rates for a shipment, please use the `RerateShipment()` method instead
- Adds a `SmartRate` struct since the structure of the `Rate` and `SmartRate` objects are different (previously the SmartRate object borrowed the Rate struct)

## v1.4.0 (2021-10-06)

- Adds support for `TaxIdentifiers`
- Remove experimental undocumented methods `ListTrackersUpdated` and `ListTrackersUpdatedWithContext`

## v1.3.1 (2021-06-23)

- Corrects `CODAmount` from a `float64` to a `string`

## v1.3.0 (2021-05-27)

- Adds `Smartrate` functionality to the `Shipments` object (available by calling `GetShipmentSmartrates()`)

## v1.2.0 (2021-04-06)

- Fix batch scan form functionality; replace GetBatchScanForms with
  CreateBatchScanForms

## v1.1.0 (2020-06-10)

- Update example code.
- Add GetRate method and example.
- Remove ListCustomsInfos, ListCustomsItems, ListOrders and ListPickups
  methods.
- Add GetEvent, ListEvents and ListEventPayloads methods. Add example code
  for using these methods as well as using the Event type in a webhook
  handler.

## v1.0.3 (2020-05-30)

- Fix issue in List actions that take query parameters

## v1.0.2 (2020-05-08)

- Add shipment options for certified and registered mail.
- Changes to unit tests.

## v1.0.1 (2020-03-17)

- Update URL in installation instructions.
- Properly format the carrier_accounts parameter in create shipment requests.
  Prior to this change, specifying a carrier account would result in a failed
  request. This removes the undocumented parameter from the CreateShipment and
  CreateShipmentWithContext methods, and adds a CarrierAccountIDs field to
  the Shipment struct type.
- Add proper struct tag to CarrierAccount.Fields so that it has the proper
  key name when serialized to JSON. Prior to this change, values in the Fields
  field were not recognized by the API.

## v1.0.0 (2019-09-03)

- Initial release.
