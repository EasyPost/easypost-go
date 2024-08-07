package easypost

// TaxIdentifier objects contain tax information used by carriers.
type TaxIdentifier struct {
	Entity         string `json:"entity,omitempty" url:"entity,omitempty"`
	TaxIdType      string `json:"tax_id_type,omitempty" url:"tax_id_type,omitempty"`
	TaxId          string `json:"tax_id,omitempty" url:"tax_id,omitempty"`
	IssuingCountry string `json:"issuing_country,omitempty" url:"issuing_country,omitempty"`
}
