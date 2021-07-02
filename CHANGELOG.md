# CHANGELOG

## NEXT

* Remove experimental undocumented methods `ListTrackersUpdated` and `ListTrackersUpdatedWithContext`

## v1.3.1 2021-06-23

* Corrects `CODAmount` from a `float64` to a `string`

## v1.3.0 2021-05-27

* Adds `Smartrate` functionality to the `Shipments` object (available by calling `GetShipmentSmartrates()`)

## 1.2.0 2021-04-06

 * Fix batch scan form functionality; replace GetBatchScanForms with
   CreateBatchScanForms

## 1.1.0 2020-06-10

 * Update example code.
 * Add GetRate method and example.
 * Remove ListCustomsInfos, ListCustomsItems, ListOrders and ListPickups
   methods.
 * Add GetEvent, ListEvents and ListEventPayloads methods. Add example code
   for using these methods as well as using the Event type in a webhook
   handler.

## 1.0.3 2020-05-30

 * Fix issue in List actions that take query parameters

## 1.0.2 2020-05-08

 * Add shipment options for certified and registered mail.
 * Changes to unit tests.

## 1.0.1 2020-03-17

 * Update URL in installation instructions.
 * Properly format the carrier_accounts parameter in create shipment requests.
   Prior to this change, specifying a carrier account would result in a failed
   request. This removes the undocumented parameter from the CreateShipment and
   CreateShipmentWithContext methods, and adds a CarrierAccountIDs field to
   the Shipment struct type.
 * Add proper struct tag to CarrierAccount.Fields so that it has the proper
   key name when serialized to JSON. Prior to this change, values in the Fields
   field were not recognized by the API.

## 1.0.0 2019-09-03

 * Initial release.
