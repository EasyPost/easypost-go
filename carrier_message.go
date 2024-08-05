package easypost

// CarrierMessage contains additional status data that is provided by some
// carriers for certain operations.
type CarrierMessage struct {
	Carrier          string `json:"carrier,omitempty" url:"carrier,omitempty"`
	Type             string `json:"type,omitempty" url:"type,omitempty"`
	Message          string `json:"message,omitempty" url:"message,omitempty"`
	CarrierAccountID string `json:"carrier_account_id,omitempty" url:"carrier_account_id,omitempty"`
}
