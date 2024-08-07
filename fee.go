package easypost

// Fee objects are used to represent the breakdown of charges made when
// purchasing an item on EasyPost.
type Fee struct {
	Object   string `json:"object,omitempty" url:"object,omitempty"`
	Type     string `json:"type,omitempty" url:"type,omitempty"`
	Amount   string `json:"amount,omitempty" url:"amount,omitempty"`
	Charged  bool   `json:"charged,omitempty" url:"charged,omitempty"`
	Refunded bool   `json:"refunded,omitempty" url:"refunded,omitempty"`
}
