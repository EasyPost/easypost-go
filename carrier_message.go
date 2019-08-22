package easypost

// CarrierMessage contains additional status data that is provided by some
// carriers for certain operations.
type CarrierMessage struct {
	Carrier          string `json:"carrier,omitempty"`
	Type             string `json:"type,omitempty"`
	Message          string `json:"message,omitempty"`
	CarrierAccountID string `json:"carrier_account_id,omitempty"`
}
