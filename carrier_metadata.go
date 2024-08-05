package easypost

import (
	"context"
	"net/http"
	"strings"
)

// CarrierMetadata represents all the metadata for a carrier.
type CarrierMetadata struct {
	Name               string                       `json:"name,omitempty" url:"name,omitempty"`
	HumanReadable      string                       `json:"human_readable,omitempty" url:"human_readable,omitempty"`
	ServiceLevels      []*MetadataServiceLevel      `json:"service_levels,omitempty" url:"service_levels,omitempty"`
	PredefinedPackages []*MetadataPredefinedPackage `json:"predefined_packages,omitempty" url:"predefined_packages,omitempty"`
	ShipmentOptions    []*MetadataShipmentOption    `json:"shipment_options,omitempty" url:"shipment_options,omitempty"`
	SupportedFeatures  []*MetadataSupportedFeature  `json:"supported_features,omitempty" url:"supported_features,omitempty"`
}

// MetadataServiceLevel represents an available service level of a carrier.
type MetadataServiceLevel struct {
	Name          string   `json:"name,omitempty" url:"name,omitempty"`
	Carrier       string   `json:"carrier,omitempty" url:"carrier,omitempty"`
	HumanReadable string   `json:"human_readable,omitempty" url:"human_readable,omitempty"`
	Description   string   `json:"description,omitempty" url:"description,omitempty"`
	Dimensions    []string `json:"dimensions,omitempty" url:"dimensions,omitempty"`
	MaxWeight     float64  `json:"max_weight,omitempty" url:"max_weight,omitempty"`
}

// MetadataPredefinedPackage represents an available predefined package of a carrier.
type MetadataPredefinedPackage struct {
	Name          string   `json:"name,omitempty" url:"name,omitempty"`
	Carrier       string   `json:"carrier,omitempty" url:"carrier,omitempty"`
	HumanReadable string   `json:"human_readable,omitempty" url:"human_readable,omitempty"`
	Description   string   `json:"description,omitempty" url:"description,omitempty"`
	Dimensions    []string `json:"dimensions,omitempty" url:"dimensions,omitempty"`
	MaxWeight     float64  `json:"max_weight,omitempty" url:"max_weight,omitempty"`
}

// MetadataShipmentOption represents an available shipment option of a carrier.
type MetadataShipmentOption struct {
	Name          string `json:"name,omitempty" url:"name,omitempty"`
	Carrier       string `json:"carrier,omitempty" url:"carrier,omitempty"`
	HumanReadable string `json:"human_readable,omitempty" url:"human_readable,omitempty"`
	Description   string `json:"description,omitempty" url:"description,omitempty"`
	Deprecated    bool   `json:"deprecated,omitempty" url:"deprecated,omitempty"`
	Type          string `json:"type,omitempty" url:"type,omitempty"`
}

// MetadataSupportedFeature represents a supported feature of a carrier.
type MetadataSupportedFeature struct {
	Name        string `json:"name,omitempty" url:"name,omitempty"`
	Carrier     string `json:"carrier,omitempty" url:"carrier,omitempty"`
	Description string `json:"description,omitempty" url:"description,omitempty"`
	Supported   bool   `json:"supported,omitempty" url:"supported,omitempty"`
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

// GetCarrierMetadataWithCarriersAndTypes retrieves carrier metadata for a list of carriers and a list of types.
func (c *Client) GetCarrierMetadataWithCarriersAndTypes(carriers []string, types []string) (out []*CarrierMetadata, err error) {
	return c.GetCarrierMetadataWithContext(context.Background(), carriers, types)
}

// GetCarrierMetadataWithContext performs the same operation as GetCarrierMetadata, but allows specifying a context that can interrupt the request.
func (c *Client) GetCarrierMetadataWithContext(ctx context.Context, carriers []string, types []string) (out []*CarrierMetadata, err error) {
	params := struct {
		Carriers string `json:"carriers,omitempty" url:"carriers,omitempty"`
		Types    string `json:"types,omitempty" url:"types,omitempty"`
	}{}

	url := "/v2/metadata/carriers"
	if carriers != nil {
		params.Carriers = strings.Join(carriers[:], ",")
	}
	if types != nil {
		params.Types = strings.Join(types[:], ",")
	}

	res := struct {
		CarrierMetadata *[]*CarrierMetadata `json:"carriers,omitempty" url:"carriers,omitempty"`
	}{CarrierMetadata: &out}

	err = c.do(ctx, http.MethodGet, url, params, &res)
	return
}
