package easypost

// CarbonOffset objects contain carbon offset details for a rate.
type CarbonOffset struct {
	Object   string `json:"object,omitempty"`
	Currency string `json:"currency,omitempty"`
	Grams    int64  `json:"grams,omitempty"`
	Price    string `json:"price,omitempty"`
}
