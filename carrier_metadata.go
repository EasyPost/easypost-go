package easypost

import (
	"context"
	"strings"
)

// BetaCarrierMetadata represents all the metadata for a carrier.
type BetaCarrierMetadata struct {
	Name               string                       `json:"name,omitempty"`
	HumanReadable      string                       `json:"human_readable,omitempty"`
	ServiceLevels      []*BetaMetadataServiceLevel      `json:"service_levels,omitempty"`
	PredefinedPackages []*BetaMetadataPredefinedPackage `json:"predefined_packages,omitempty"`
	ShipmentOptions    []*BetaMetadataShipmentOption    `json:"shipment_options,omitempty"`
	SupportedFeatures  []*BetaMetadataSupportedFeature  `json:"supported_features,omitempty"`
}

// BetaMetadataServiceLevel represents an available service level of a carrier.
type BetaMetadataServiceLevel struct {
	Name          string   `json:"name,omitempty"`
	Carrier       string   `json:"carrier,omitempty"`
	HumanReadable string   `json:"human_readable,omitempty"`
	Description   string   `json:"description,omitempty"`
	Dimensions    []string `json:"dimensions,omitempty"`
	MaxWeight     float64  `json:"max_weight,omitempty"`
}

// BetaMetadataPredefinedPackage represents an available predefined package of a carrier.
type BetaMetadataPredefinedPackage struct {
	Name          string   `json:"name,omitempty"`
	Carrier       string   `json:"carrier,omitempty"`
	HumanReadable string   `json:"human_readable,omitempty"`
	Description   string   `json:"description,omitempty"`
	Dimensions    []string `json:"dimensions,omitempty"`
	MaxWeight     float64  `json:"max_weight,omitempty"`
}

// BetaMetadataShipmentOption represents an available shipment option of a carrier.
type BetaMetadataShipmentOption struct {
	Name          string `json:"name,omitempty"`
	Carrier       string `json:"carrier,omitempty"`
	HumanReadable string `json:"human_readable,omitempty"`
	Description   string `json:"description,omitempty"`
	Deprecated    bool   `json:"deprecated,omitempty"`
	Type          string `json:"type,omitempty"`
}

// BetaMetadataSupportedFeature represents a supported feature of a carrier.
type BetaMetadataSupportedFeature struct {
	Name        string `json:"name,omitempty"`
	Carrier     string `json:"carrier,omitempty"`
	Description string `json:"description,omitempty"`
	Supported   bool   `json:"supported,omitempty"`
}

// BetaGetCarrierMetadata retrieves all metadata for all carriers on the EasyPost platform.
func (c *Client) BetaGetCarrierMetadata() (out []*BetaCarrierMetadata, err error) {
	return c.BetaGetCarrierMetadataWithContext(context.Background(), nil, nil)
}

// BetaGetCarrierMetadataWithCarriers retrieves carrier metadata for a list of carriers.
func (c *Client) BetaGetCarrierMetadataWithCarriers(carriers []string) (out []*BetaCarrierMetadata, err error) {
	return c.BetaGetCarrierMetadataWithContext(context.Background(), carriers, nil)
}

// BetaGetCarrierMetadataWithTypes retrieves carrier metadata for a list of carriers.
func (c *Client) BetaGetCarrierMetadataWithTypes(types []string) (out []*BetaCarrierMetadata, err error) {
	return c.BetaGetCarrierMetadataWithContext(context.Background(), nil, types)
}

// BetaGetCarrierMetadataWithTypes retrieves carrier metadata for a list of carriers and a list of types.
func (c *Client) BetaGetCarrierMetadataWithCarriersAndTypes(carriers []string, types []string) (out []*BetaCarrierMetadata, err error) {
	return c.BetaGetCarrierMetadataWithContext(context.Background(), carriers, types)
}

// BetaGetCarrierMetadataWithContext performs the same operation as
// BetaGetCarrierMetadata, but allows specifying a context that can interrupt the
// request.
func (c *Client) BetaGetCarrierMetadataWithContext(ctx context.Context, carriers []string, types []string) (out []*BetaCarrierMetadata, err error) {
	url := "/beta/metadata"
	if carriers != nil && types != nil {
		url = url + "?"
	}
	if carriers != nil {
		url = url + "carriers=" + strings.Join(carriers[:], ",")
	}
	if carriers != nil && types != nil {
		url = url + "&"
	}
	if types != nil {
		url = url + "types=" + strings.Join(types[:], ",")
	}

	res := struct {
		BetaCarrierMetadata *[]*BetaCarrierMetadata `json:"carriers,omitempty"`
	}{BetaCarrierMetadata: &out}
	err = c.get(ctx, url, &res)
	return
}
