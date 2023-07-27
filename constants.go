package easypost

func getCarrierAccountTypesWithCustomWorkflows() []string {
	return []string{"FedexAccount", "UpsAccount"}
}

var ApiErrorDetailsParsingError = "RESPONSE.PARSE_ERROR"
var ApiDidNotReturnErrorDetails = "API did not return error details"
var NoPagesLeftToRetrieve = "There are no more pages to retrieve"
var MismatchWebhookSignature = "Webhook received did not originate from EasyPost or had a webhook secret mismatch"
var MissingWebhookSignature = "Webhook does not contain a valid HMAC signature."
