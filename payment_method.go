package easypost

type PaymentMethodPriority int64

const (
	PrimaryPaymentMethodPriority PaymentMethodPriority = iota
	SecondaryPaymentMethodPriority
)

type PaymentMethodType int64

const (
	CreditCardPaymentType PaymentMethodType = iota
	BankAccountPaymentType
)

type PaymentMethod struct {
	ID                     string               `json:"id,omitempty" url:"id,omitempty"`
	Object                 string               `json:"object,omitempty" url:"object,omitempty"`
	PrimaryPaymentMethod   *PaymentMethodObject `json:"primary_payment_method,omitempty" url:"primary_payment_method,omitempty"`
	SecondaryPaymentMethod *PaymentMethodObject `json:"secondary_payment_method,omitempty" url:"secondary_payment_method,omitempty"`
}

type PaymentMethodObject struct {
	BankName                  string `json:"bank_name,omitempty" url:"bank_name,omitempty"`                                     // bank account
	Brand                     string `json:"brand,omitempty" url:"brand,omitempty"`                                             // credit card
	Country                   string `json:"country,omitempty" url:"country,omitempty"`                                         // bank account
	DisabledAt                string `json:"disabled_at,omitempty" url:"disabled_at,omitempty"`                                 // both
	ExpirationMonth           int    `json:"exp_month,omitempty" url:"exp_month,omitempty"`                                     // credit card
	ExpirationYear            int    `json:"exp_year,omitempty" url:"exp_year,omitempty"`                                       // credit card
	ID                        string `json:"id,omitempty" url:"id,omitempty"`                                                   // both
	Last4                     string `json:"last4,omitempty" url:"last4,omitempty"`                                             // both
	Name                      string `json:"name,omitempty" url:"name,omitempty"`                                               // credit card
	Object                    string `json:"object,omitempty" url:"object,omitempty"`                                           // both
	Verified                  bool   `json:"verified,omitempty" url:"verified,omitempty"`                                       // bank account
	RequiresMandateCollection bool   `json:"requires_mandate_collection,omitempty" url:"requires_mandate_collection,omitempty"` // both
}

// getPaymentMethodObjectType returns the PaymentMethodType enum of a PaymentMethodObject.
func (c *Client) getPaymentMethodObjectType(object *PaymentMethodObject) (out PaymentMethodType, err error) {
	if object.Object == "CreditCard" {
		out = CreditCardPaymentType
	} else if object.Object == "BankAccount" {
		out = BankAccountPaymentType
	} else {
		return out, newInvalidObjectError(NoMatchingPaymentMethod)
	}
	return
}

// getPaymentMethodEndpoint returns the associated endpoint for a PaymentMethodObject based on its type.
func (c *Client) getPaymentMethodEndpoint(object *PaymentMethodObject) (out string, err error) {
	paymentMethodType, _ := c.getPaymentMethodObjectType(object)
	switch paymentMethodType {
	case CreditCardPaymentType:
		out = "credit_cards"
	case BankAccountPaymentType:
		out = "bank_accounts"
	default:
		return out, newInvalidObjectError(NoMatchingPaymentMethod)
	}
	return
}

// getPaymentMethodByPriority returns the PaymentMethodObject associated with the given PaymentMethodPriority.
func (c *Client) getPaymentMethodByPriority(priority PaymentMethodPriority) (out *PaymentMethodObject, err error) {
	paymentMethods, err := c.RetrievePaymentMethods()
	if err != nil {
		return out, err
	}

	switch priority {
	case PrimaryPaymentMethodPriority:
		out = paymentMethods.PrimaryPaymentMethod
	case SecondaryPaymentMethodPriority:
		out = paymentMethods.SecondaryPaymentMethod
	}

	if out == nil {
		return out, newInvalidObjectError(PaymentMethodNotSetUp)
	}

	return out, nil
}

// GetPaymentEndpointByPrimaryOrSecondary gets the payment priority based on primaryOrSecondary parameter.
func (c *Client) GetPaymentEndpointByPrimaryOrSecondary(primaryOrSecondary PaymentMethodPriority) (out string) {
	switch primaryOrSecondary {
	case PrimaryPaymentMethodPriority:
		out = "primary"
	case SecondaryPaymentMethodPriority:
		out = "secondary"
	}

	return
}
