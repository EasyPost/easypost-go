package easypost

import (
	"context"
	"strings"
)

// CarrierMetadata represents all the metadata for a carrier.
type CarrierMetadata struct {
	Name               string                       `json:"name,omitempty"`
	HumanReadable      string                       `json:"human_readable,omitempty"`
	ServiceLevels      []*MetadataServiceLevel      `json:"service_levels,omitempty"`
	PredefinedPackages []*MetadataPredefinedPackage `json:"predefined_packages,omitempty"`
	ShipmentOptions    []*MetadataShipmentOption    `json:"shipment_options,omitempty"`
	SupportedFeatures  []*MetadataSupportedFeature  `json:"supported_features,omitempty"`
}

// MetadataServiceLevel represents an available service level of a carrier.
type MetadataServiceLevel struct {
	Name          string   `json:"name,omitempty"`
	Carrier       string   `json:"carrier,omitempty"`
	HumanReadable string   `json:"human_readable,omitempty"`
	Description   string   `json:"description,omitempty"`
	Dimensions    []string `json:"dimensions,omitempty"`
	MaxWeight     float64  `json:"max_weight,omitempty"`
}

// MetadataPredefinedPackage represents an available predefined package of a carrier.
type MetadataPredefinedPackage struct {
	Name          string   `json:"name,omitempty"`
	Carrier       string   `json:"carrier,omitempty"`
	HumanReadable string   `json:"human_readable,omitempty"`
	Description   string   `json:"description,omitempty"`
	Dimensions    []string `json:"dimensions,omitempty"`
	MaxWeight     float64  `json:"max_weight,omitempty"`
}

// MetadataShipmentOption represents an available shipment option of a carrier.
type MetadataShipmentOption struct {
	Name          string `json:"name,omitempty"`
	Carrier       string `json:"carrier,omitempty"`
	HumanReadable string `json:"human_readable,omitempty"`
	Description   string `json:"description,omitempty"`
	Deprecated    bool   `json:"deprecated,omitempty"`
	Type          string `json:"type,omitempty"`
}

// MetadataSupportedFeature represents a supported feature of a carrier.
type MetadataSupportedFeature struct {
	Name        string `json:"name,omitempty"`
	Carrier     string `json:"carrier,omitempty"`
	Description string `json:"description,omitempty"`
	Supported   bool   `json:"supported,omitempty"`
}

// GetCarrierMetadata retrieves all metadata for all carriers on the EasyPost platform.
func (c *Client) GetCarrierMetadata() (out []*CarrierMetadata, err error) {
	return c.GetCarrierMetadataWithContext(context.Background(), nil, nil)
}

// GetCarrierMetadataWithCarriers retrieves carrier metadata for a list of carriers.
func (c *Client) GetCarrierMetadataWithCarriers(carriers []string) (out []*CarrierMetadata, err error) {
	return c.GetCarrierMetadataWithContext(context.Background(), carriers, nil)
}

// GetCarrierMetadataWithTypes retrieves carrier metadata for a list of carriers.
func (c *Client) GetCarrierMetadataWithTypes(types []string) (out []*CarrierMetadata, err error) {
	return c.GetCarrierMetadataWithContext(context.Background(), nil, types)
}

// GetCarrierMetadataWithTypes retrieves carrier metadata for a list of carriers and a list of types.
func (c *Client) GetCarrierMetadataWithCarriersAndTypes(carriers []string, types []string) (out []*CarrierMetadata, err error) {
	return c.GetCarrierMetadataWithContext(context.Background(), carriers, types)
}

// GetCarrierMetadataWithContext performs the same operation as
// GetCarrierMetadata, but allows specifying a context that can interrupt the
// request.
func (c *Client) GetCarrierMetadataWithContext(ctx context.Context, carriers []string, types []string) (out []*CarrierMetadata, err error) {
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
		CarrierMetadata *[]*CarrierMetadata `json:"carriers,omitempty"`
	}{CarrierMetadata: &out}
	err = c.get(ctx, url, &res)
	return
}
