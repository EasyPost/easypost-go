package easypost

func getCarrierAccountTypesWithCustomWorkflows() []string {
	return []string{"FedexAccount", "UpsAccount"}
}

var ApiErrorDetailsParsingError = "RESPONSE.PARSE_ERROR"
var ApiDidNotReturnErrorDetails = "API did not return error details"
var NoPagesLeftToRetrieve = "There are no more pages to retrieve"
var MismatchWebhookSignature = "Webhook received did not originate from EasyPost or had a webhook secret mismatch"
var MissingWebhookSignature = "Webhook does not contain a valid HMAC signature."
var NoPaymentMethods = "No payment methods are set up. Please add a payment method and try again."
var InvalidParameter = "Invalid parameter: "
var JsonDeserializationErrorMessage = "Error deserializing JSON into object of type "
var JsonSerializationErrorMessage = "Error serializing object of type "
var JsonNoDataErrorMessage = "No data was provided to serialize"
var MissingRequiredParameter = "Missing required parameter: "
var MissingProperty = "Missing property: "
var NoMatchingPaymentMethod = "No matching payment method type found"
var NoRatesFoundMatchingFilters = "No rates found matching the given filters"
var PaymentMethodNotSetUp = "The chosen payment method is not set up yet"
