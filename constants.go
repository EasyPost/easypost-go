package easypost

func getCarrierAccountTypesWithCustomWorkflows() []string {
	return []string{"FedexAccount", "FedexSmartpostAccount"}
}

// By the nature of Golang and to ensure backwards compatibility, we cannot
// consolidate this list into the `getCarrierAccountTypesWithCustomOauth` as
// we've done in the other libraries since the structs differ.
func getUpsCarrierAccountTypes() []string {
	return []string{"UpsAccount", "UpsMailInnovationsAccount", "UpsSurepostAccount"}
}

func getCarrierAccountTypesWithCustomOauth() []string {
	return []string{"AmazonShippingAccount"}
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
