package easypost

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// FedExAccountValidationResponse represents the response from FedEx account validation endpoints.
type FedExAccountValidationResponse struct {
	// If the response contains the following, one must complete pin or invoice validation next
	EmailAddress string   `json:"email_address,omitempty" url:"email_address,omitempty"`
	Options      []string `json:"options,omitempty" url:"options,omitempty"`
	PhoneNumber  string   `json:"phone_number,omitempty" url:"phone_number,omitempty"`

	// If the response contains the following, pre-validation has been completed
	ID          string            `json:"id,omitempty" url:"id,omitempty"`
	Object      string            `json:"object,omitempty" url:"object,omitempty"`
	Type        string            `json:"type,omitempty" url:"type,omitempty"`
	Credentials map[string]string `json:"credentials,omitempty" url:"credentials,omitempty"`
}

// FedExRequestPinResponse represents the response from requesting a PIN.
type FedExRequestPinResponse struct {
	Message string `json:"message,omitempty" url:"message,omitempty"`
}

// RegisterFedExAddress registers the billing address for a FedEx account.
// Advanced method for custom parameter structures.
func (c *Client) RegisterFedExAddress(fedexAccountNumber string, params map[string]interface{}) (out *FedExAccountValidationResponse, err error) {
	return c.RegisterFedExAddressWithContext(context.Background(), fedexAccountNumber, params)
}

// RegisterFedExAddressWithContext performs the same operation as RegisterFedExAddress, but allows specifying a context that can interrupt the request.
func (c *Client) RegisterFedExAddressWithContext(ctx context.Context, fedexAccountNumber string, params map[string]interface{}) (out *FedExAccountValidationResponse, err error) {
	wrappedParams := wrapAddressValidation(params)
	endpoint := fmt.Sprintf("fedex_registrations/%s/address", fedexAccountNumber)
	err = c.do(ctx, http.MethodPost, endpoint, wrappedParams, &out)
	return
}

// RequestFedExPin requests a PIN for FedEx account verification.
func (c *Client) RequestFedExPin(fedexAccountNumber string, pinMethodOption string) (out *FedExRequestPinResponse, err error) {
	return c.RequestFedExPinWithContext(context.Background(), fedexAccountNumber, pinMethodOption)
}

// RequestFedExPinWithContext performs the same operation as RequestFedExPin, but allows specifying a context that can interrupt the request.
func (c *Client) RequestFedExPinWithContext(ctx context.Context, fedexAccountNumber string, pinMethodOption string) (out *FedExRequestPinResponse, err error) {
	wrappedParams := map[string]interface{}{
		"pin_method": map[string]interface{}{
			"option": pinMethodOption,
		},
	}
	endpoint := fmt.Sprintf("fedex_registrations/%s/pin", fedexAccountNumber)
	err = c.do(ctx, http.MethodPost, endpoint, wrappedParams, &out)
	return
}

// ValidateFedExPin validates the PIN entered by the user for FedEx account verification.
func (c *Client) ValidateFedExPin(fedexAccountNumber string, params map[string]interface{}) (out *FedExAccountValidationResponse, err error) {
	return c.ValidateFedExPinWithContext(context.Background(), fedexAccountNumber, params)
}

// ValidateFedExPinWithContext performs the same operation as ValidateFedExPin, but allows specifying a context that can interrupt the request.
func (c *Client) ValidateFedExPinWithContext(ctx context.Context, fedexAccountNumber string, params map[string]interface{}) (out *FedExAccountValidationResponse, err error) {
	wrappedParams := wrapPinValidation(params)
	endpoint := fmt.Sprintf("fedex_registrations/%s/pin/validate", fedexAccountNumber)
	err = c.do(ctx, http.MethodPost, endpoint, wrappedParams, &out)
	return
}

// SubmitFedExInvoice submits invoice information to complete FedEx account registration.
func (c *Client) SubmitFedExInvoice(fedexAccountNumber string, params map[string]interface{}) (out *FedExAccountValidationResponse, err error) {
	return c.SubmitFedExInvoiceWithContext(context.Background(), fedexAccountNumber, params)
}

// SubmitFedExInvoiceWithContext performs the same operation as SubmitFedExInvoice, but allows specifying a context that can interrupt the request.
func (c *Client) SubmitFedExInvoiceWithContext(ctx context.Context, fedexAccountNumber string, params map[string]interface{}) (out *FedExAccountValidationResponse, err error) {
	wrappedParams := wrapInvoiceValidation(params)
	endpoint := fmt.Sprintf("fedex_registrations/%s/invoice", fedexAccountNumber)
	err = c.do(ctx, http.MethodPost, endpoint, wrappedParams, &out)
	return
}

// wrapAddressValidation wraps address validation parameters and ensures the "name" field exists.
// If not present, generates a UUID (with hyphens removed) as the name.
func wrapAddressValidation(params map[string]interface{}) map[string]interface{} {
	wrappedParams := make(map[string]interface{})

	if addressValidation, ok := params["address_validation"].(map[string]interface{}); ok {
		addressValidationCopy := make(map[string]interface{})
		for k, v := range addressValidation {
			addressValidationCopy[k] = v
		}
		ensureNameField(addressValidationCopy)
		wrappedParams["address_validation"] = addressValidationCopy
	}

	if easypostDetails, ok := params["easypost_details"]; ok {
		wrappedParams["easypost_details"] = easypostDetails
	}

	return wrappedParams
}

// wrapPinValidation wraps PIN validation parameters and ensures the "name" field exists.
// If not present, generates a UUID (with hyphens removed) as the name.
func wrapPinValidation(params map[string]interface{}) map[string]interface{} {
	wrappedParams := make(map[string]interface{})

	if pinValidation, ok := params["pin_validation"].(map[string]interface{}); ok {
		pinValidationCopy := make(map[string]interface{})
		for k, v := range pinValidation {
			pinValidationCopy[k] = v
		}
		ensureNameField(pinValidationCopy)
		wrappedParams["pin_validation"] = pinValidationCopy
	}

	if easypostDetails, ok := params["easypost_details"]; ok {
		wrappedParams["easypost_details"] = easypostDetails
	}

	return wrappedParams
}

// wrapInvoiceValidation wraps invoice validation parameters and ensures the "name" field exists.
// If not present, generates a UUID (with hyphens removed) as the name.
func wrapInvoiceValidation(params map[string]interface{}) map[string]interface{} {
	wrappedParams := make(map[string]interface{})

	if invoiceValidation, ok := params["invoice_validation"].(map[string]interface{}); ok {
		invoiceValidationCopy := make(map[string]interface{})
		for k, v := range invoiceValidation {
			invoiceValidationCopy[k] = v
		}
		ensureNameField(invoiceValidationCopy)
		wrappedParams["invoice_validation"] = invoiceValidationCopy
	}

	if easypostDetails, ok := params["easypost_details"]; ok {
		wrappedParams["easypost_details"] = easypostDetails
	}

	return wrappedParams
}

// ensureNameField ensures the "name" field exists in the provided map.
// If not present, generates a UUID (with hyphens removed) as the name.
// This follows the pattern used in the web UI implementation.
func ensureNameField(m map[string]interface{}) {
	if _, exists := m["name"]; !exists || m["name"] == nil {
		uuidStr := strings.ReplaceAll(uuid.New().String(), "-", "")
		m["name"] = uuidStr
	}
}
