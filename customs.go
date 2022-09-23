package easypost

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// CustomsInfo objects contain CustomsItem objects and all necessary information
// for the generation of customs forms required for international shipping.
type CustomsInfo struct {
	ID                  string         `json:"id,omitempty"`
	Object              string         `json:"object,omitempty"`
	CreatedAt           *time.Time     `json:"created_at,omitempty"`
	UpdatedAt           *time.Time     `json:"updated_at,omitempty"`
	EELPFC              string         `json:"eel_pfc,omitempty"`
	ContentsType        string         `json:"contents_type,omitempty"`
	ContentsExplanation string         `json:"contents_explanation,omitempty"`
	CustomsCertify      bool           `json:"customs_certify,omitempty"`
	CustomsSigner       string         `json:"customs_signer,omitempty"`
	NonDeliveryOption   string         `json:"non_delivery_option,omitempty"`
	RestrictionType     string         `json:"restriction_type,omitempty"`
	CustomsItems        []*CustomsItem `json:"customs_items,omitempty"`
	Declaration         string         `json:"declaration,omitempty"`
}

// A CustomsItem object describes goods for international shipment.
type CustomsItem struct {
	ID             string     `json:"id,omitempty"`
	Object         string     `json:"object,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	Description    string     `json:"description,omitempty"`
	Quantity       float64    `json:"quantity,omitempty"`
	Value          float64    `json:"value,omitempty"`
	Weight         float64    `json:"weight,omitempty"`
	HSTariffNumber string     `json:"hs_tariff_number,omitempty"`
	Code           string     `json:"code,omitempty"`
	OriginCountry  string     `json:"origin_country,omitempty"`
	Currency       string     `json:"currency,omitempty"`
}

type createCustomsInfoRequest struct {
	CustomsInfo *CustomsInfo `json:"customs_info,omitempty"`
}

type createCustomsItemRequest struct {
	CustomsItem *CustomsItem `json:"customs_item,omitempty"`
}

func (ci *CustomsItem) MarshalJSON() ([]byte, error) {
	// convert a CustomsItem to JSON data byte array

	// needed later to avoid infinite recursion
	type Alias CustomsItem

	return json.Marshal(&struct {
		Value string `json:"value,omitempty"`
		*Alias
	}{
		// always convert the Value into a string when serializing to JSON
		Value: fmt.Sprintf("%f", ci.Value),
		Alias: (*Alias)(ci),
	})
}

func (ci *CustomsItem) UnmarshalJSON(b []byte) (err error) {
	// convert a JSON data byte array to CustomsItem

	// needed later to avoid infinite recursion
	type Alias CustomsItem

	// make a raw copy of the data
	var rawData map[string]interface{}
	if err = json.Unmarshal(b, &rawData); err != nil {
		return
	}

	// tweak values
	for key, value := range rawData {
		switch key {
		case "value":
			// convert Value property to a float64 if it is a string (via API response)
			if valueString, noErr := value.(string); noErr {
				if valueFloat, err := strconv.ParseFloat(valueString, 64); err == nil {
					rawData[key] = valueFloat
				}
				return
			} else {
				// keep Value property a float64 if it is already one (via Fixture data)
				rawData[key] = value.(float64)
			}
		default:
			// no tweak needed, keep the value as is
			rawData[key] = value
		}
	}

	// convert the tweaked data back into a byte array
	rawDataCorrectTypesBytes, err := json.Marshal(rawData)
	if err != nil {
		return
	}

	// now convert the tweaked byte array into a CustomsItem via an Alias
	tmp := &struct {
		*Alias
	}{
		Alias: (*Alias)(ci),
	}
	err = json.Unmarshal(rawDataCorrectTypesBytes, &tmp)

	return
}

// CreateCustomsInfo creates a new CustomsInfo object.
//
//	c := easypost.New(MyEasyPostAPIKey)
//	out, err := c.CreateCustomsInfo(
//		&easypost.CustomsInfo{
//			CustomsCertify:  true,
//			CustomsSigner:   "Steve Brule",
//			CustomsType:     "merchandise",
//			RestrictionType: "none",
//			EELPFC:          "NOEEI 30.37(a)",
//			CustomsItems:    []*easypost.CustomsItem{
//				&easypost.CustomsItem{ID: "cstitem_2002"},
//				&easypost.CustomsItem{
//					Description:     "Sweet shirts",
//					Quantity:        2,
//					Value:           23,
//					Weight:          11,
//					HSTariffNumber: "654321",
//					OriginCountry:   "US",
//				},
//			},
//		},
//	)
func (c *Client) CreateCustomsInfo(in *CustomsInfo) (out *CustomsInfo, err error) {
	return c.CreateCustomsInfoWithContext(context.Background(), in)
}

// CreateCustomsInfoWithContext performs the same operation as
// CreateCustomsInfo, but allows specifying a context that can interrupt the
// request.
func (c *Client) CreateCustomsInfoWithContext(ctx context.Context, in *CustomsInfo) (out *CustomsInfo, err error) {
	req := &createCustomsInfoRequest{CustomsInfo: in}
	err = c.post(ctx, "customs_infos", req, &out)
	return
}

// GetCustomsInfo returns the CustomsInfo object with the given ID or reference.
func (c *Client) GetCustomsInfo(customsInfoID string) (out *CustomsInfo, err error) {
	return c.GetCustomsInfoWithContext(context.Background(), customsInfoID)
}

// GetCustomsInfoWithContext performs the same operation as GetCustomsInfo, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetCustomsInfoWithContext(ctx context.Context, customsInfoID string) (out *CustomsInfo, err error) {
	err = c.get(ctx, "customs_infos/"+customsInfoID, &out)
	return
}

// CreateCustomsItem creates a new CustomsItem object.
//
//	c := easypost.New(MyEasyPostAPIKey)
//	out, err := c.CreateCustomsItem(
//		&easypost.CustomsItem{
//			Description:    "T-shirt",
//			Quantity:       1,
//			Weight:         5,
//			Value:          10,
//			HSTariffNumber: "123456",
//			OriginCountry:  "US",
//		},
//	)
func (c *Client) CreateCustomsItem(in *CustomsItem) (out *CustomsItem, err error) {
	return c.CreateCustomsItemWithContext(context.Background(), in)
}

// CreateCustomsItemWithContext performs the same operation as
// CreateCustomsItem, but allows specifying a context that can interrupt the
// request.
func (c *Client) CreateCustomsItemWithContext(ctx context.Context, in *CustomsItem) (out *CustomsItem, err error) {
	req := &createCustomsItemRequest{CustomsItem: in}
	err = c.post(ctx, "customs_items", req, &out)
	return
}

// GetCustomsItem returns the CustomsInfo object with the given ID or reference.
func (c *Client) GetCustomsItem(customsItemID string) (out *CustomsItem, err error) {
	return c.GetCustomsItemWithContext(context.Background(), customsItemID)
}

// GetCustomsItemWithContext performs the same operation as GetCustomsItem, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetCustomsItemWithContext(ctx context.Context, customsItemID string) (out *CustomsItem, err error) {
	err = c.get(ctx, "customs_items/"+customsItemID, &out)
	return
}
