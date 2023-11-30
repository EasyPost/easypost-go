# CHANGELOG

## Next major release

- Removes the undocumented `CreateAndBuy` function from the Batch service. The proper usage is to create a batch first and buy it separately
- Removed `CarbonOffset` parameter from `createShipmentRequest`, `buyShipmentRequest`, and `buyShipmentRequest` structs as EasyPost now offers Carbon Neutral shipments by default for free

## v3.2.0 (2023-10-11)

- Add `GetAPIKeysForUser` under API Key service to retrieve API keys for a specific user ID

## v3.1.0 (2023-09-29)

- Adds `CommercialInvoiceSignature` and `CommercialInvoiceLetterhead` shipment options

## v3.0.1 (2023-09-05)

- Fix endpoint for creating a FedEx Smartpost carrier account

## v3.0.0 (2023-08-03)

- Drop support for Go 1.15
  - Minimum supported version is now Go 1.16
- Adds specific error types for API and local errors
  - All errors inherit the `LibraryError` interface
    - API errors inherit the `APIError` interface, which inherits the `LibraryError` interface
    - Local errors inherit the `LocalError` interface, which inherits the `LibraryError` interface
  - Common error messages are now available as constants for parsing (e.g. regex)
- All uses of `*time.Time` have been replaced with new `easypost.DateTime` class
  - `DateTime` is a wrapper around `time.Time` that handles unexpected date formats from the API
- Previously-marked deprecated functions and structs have been removed:
  - `ListReportOptions` struct -> use `ListOptions` struct instead
  - Beta carrier metadata functions -> use non-beta functions instead
  - `LowestRate` shipment functions -> use `LowestShipmentRate` functions instead
  - `CreateWebhook` function -> use `CreateWebhookWithDetails` function instead
  - `EnableWebhook` function -> use `UpdateWebhook` function instead
- `TrackingCodes` string array on `ListTrackersOptions` struct is now `TrackingCode` single string
- `AddShipmentsToBatch` and `RemoveShipmentsFromBatch` functions now explicitly accept `Shipment` structs instead of generic `interface{}` types
  - Functions will no longer accept solely IDs; users will need to provide whole `Shipment` structs

## v2.20.0 (2023-07-28)

- Adds hooks to introspect the request and response of an API call (see `HTTP Hooks` section of the README for more information)

## v2.19.0 (2023-06-06)

- Migrates Carrier Metadata to GA (now available via `GetCarrierMetadata`)

## v2.18.0 (2023-05-02)

- Adds `GetShipmentEstimatedDeliveryDate` and `GetShipmentEstimatedDeliveryDateWithContext` functions

## v2.17.0 (2023-04-26)

- Adds `SuppressETD` to `ShipmentOptions`

## v2.16.0 (2023-04-26)

- Adds missing `CustomsInfo` to the Order struct

## v2.15.0 (2023-04-18)

- Adds `BetaGetCarrierMetadata`, `BetaGetCarrierMetadataWithCarriers`, `BetaGetCarrierMetadataWithTypes`, and `BetaGetCarrierMetadataWithCarriersAndTypes` function
- The `Message` attribute of an `APIError` is now an `interface{}` rather than a `string`, due to potential inconsistent data structure from the API response
  - Behind-the-scenes, the `message` portion of the JSON response is transformed to a concatenated string. Users should be able to safely cast the `Message` attribute to a string when accessing it via `myApiError.Message.(string)`

## v2.14.0 (2023-04-04)

- Adds `GetNextXPage` functions (e.g. `GetNextShipmentPage`) to retrieve the next page of results for a provided paginated list

## v2.13.1 (2023-03-22)

- Adds missing `StatusDetail` property to `Tracker` and `TrackingDetail` structs
- Removes unusable `FullTestTracker` field from `CreateTrackerOptions`

## v2.13.0 (2023-02-17)

- Adds `StatelessRate` object struct, `BetaGetStatelessRates` function to get ephemeral rates.
- Adds `LowestStatelessRate`, `LowestStatelessRateWithCarrier` and `LowestStatelessRateWithCarrierAndService` functions to get the lowest stateless rate from a list of stateless rates.

## v2.12.0 (2023-01-18)

- Adds function to retrieve a specific payload for a specific event, `GetEventPayload`
- Adds `ListPickups` function to `Pickup` to retrieve all pickups

## v2.11.1 (2023-01-11)

- Removes `BetaAddPaymentMethodWithPrimaryOrSecondary` and adds the `primaryOrSecondary` parameter from that function to the `BetaAddPaymentMethod` function which was the initial intention

## v2.11.0 (2023-01-11)

- Adds new beta billing functionality for ReferralCustomer users
  - `BetaAddPaymentMethod` can add a pre-existing Stripe bank account or credit card to your EasyPost account
  - `BetaRefundByAmount` refunds your wallet by a dollar amount
  - `BetaRefundByPaymentLog` refunds you wallet by a PaymentLog ID
- Adds new `PickupMinDatetime` and `DeliveryMaxDatetime` Shipment options

## v2.10.0 (2022-12-07)

- Routes requests for creating a carrier account with a custom workflow (eg: FedEx, UPS) to the correct endpoint when using the `CreateCarrierAccount` function

## v2.9.0 (2022-11-15)

- Adds `LabelSize` shipment option

## v2.8.0 (2022-09-21)

- Adds `EndShipperID` shipment option
- Adds support to pass an `EndShipperID` when buying a shipment
- Adds White Label support:
  - Create a referral customer
  - Update a referral customer's email address
  - List all referral customers
  - Add a credit card to a referral customer's account

## v2.7.0 (2022-08-25)

- Adds `CreateEndShipper()`, `GetEndShipper()`, `ListEndShippers()`, and `UpdateEndShippers()` functions
- Adds `duty_payment` and `duty_payment_account` to Shipment options

## v2.6.0 (2022-08-02)

- Adds Carbon Offset support
  - Adds the ability to create a shipment with carbon offset
  - Adds the ability to buy a shipment with carbon offset
  - Adds the ability to one-call-buy a shipment with carbon offset
  - Adds the ability to re-rate a shipment with carbon offset
- Adds `ValidateWebhook` function that returns your webhook or raises an error if there is a `webhook_secret` mismatch

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
