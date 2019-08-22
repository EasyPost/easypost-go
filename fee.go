package easypost

// Fee objects are used to represent the breakdown of charges made when
// purchasing an item on EasyPost.
type Fee struct {
	Object   string `json:"object,omitempty"`
	Type     string `json:"type,omitempty"`
	Amount   string `json:"amount,omitempty"`
	Charged  bool   `json:"charged,omitempty"`
	Refunded bool   `json:"refunded,omitempty"`
}
