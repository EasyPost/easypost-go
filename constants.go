package easypost

import "strings"

func getCarrierAccountCreateEndpointTypes() map[string][]string {
	return map[string][]string{
		"carrier_accounts/open": {
			"FedexDefaultAccount",
			"GofoWalletAccount",
			"UpsDapAccount",
			"UniUniWalletAccount",
			"VehoWalletAccount",
		},
		"carrier_accounts/register": {
			"FedexAccount",
			"FedexSmartpostAccount",
		},
		"carrier_accounts/register_oauth": {
			"AmazonShippingAccount",
			"UpsAccount",
			"UpsMailInnovationsAccount",
			"UpsSurepostAccount",
			"UspsShipAccount",
		},
	}
}

func getCarrierAccountTypesForEndpoint(endpoint string) []string {
	carrierAccountTypes, ok := getCarrierAccountCreateEndpointTypes()[endpoint]
	if !ok {
		return []string{}
	}

	return append([]string(nil), carrierAccountTypes...)
}

// By the nature of Golang and to ensure backwards compatibility, we cannot
// consolidate this list into the `getCarrierAccountTypesWithCustomOauth` as
// we've done in the other libraries since the structs differ.
func getUpsCarrierAccountTypes() []string {
	var upsCarrierAccountTypes []string
	for _, carrierAccountType := range getCarrierAccountTypesForEndpoint("carrier_accounts/register_oauth") {
		if strings.HasPrefix(carrierAccountType, "Ups") {
			upsCarrierAccountTypes = append(upsCarrierAccountTypes, carrierAccountType)
		}
	}

	return upsCarrierAccountTypes
}

func getCarrierAccountTypesWithCustomOauth() []string {
	var oauthCarrierAccountTypes []string
	for _, carrierAccountType := range getCarrierAccountTypesForEndpoint("carrier_accounts/register_oauth") {
		if !strings.HasPrefix(carrierAccountType, "Ups") {
			oauthCarrierAccountTypes = append(oauthCarrierAccountTypes, carrierAccountType)
		}
	}

	return oauthCarrierAccountTypes
}

var ApiDidNotReturnErrorDetails = "API did not return error details"
var ApiErrorDetailsParsingError = "RESPONSE.PARSE_ERROR"
var InvalidParameter = "Invalid parameter: "
var JsonDeserializationErrorMessage = "Error deserializing JSON into object of type "
var JsonNoDataErrorMessage = "No data was provided to serialize"
var JsonSerializationErrorMessage = "Error serializing object of type "
var MismatchWebhookSignature = "Webhook received did not originate from EasyPost or had a webhook secret mismatch"
var MissingProperty = "Missing property: "
var MissingRequiredParameter = "Missing required parameter: "
var MissingWebhookSignature = "Webhook does not contain a valid HMAC signature."
var NoMatchingPaymentMethod = "No matching payment method type found"
var NoPagesLeftToRetrieve = "There are no more pages to retrieve"
var NoPaymentMethods = "No payment methods are set up. Please add a payment method and try again."
var NoRatesFoundMatchingFilters = "No rates found matching the given filters"
var NoUserFoundForId = "No user found for the given ID"
var PaymentMethodNotSetUp = "The chosen payment method is not set up yet"
