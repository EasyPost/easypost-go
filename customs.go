package easypost

import (
	"context"
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
}

// A CustomsItem object describes goods for international shipment.
type CustomsItem struct {
	ID             string     `json:"id,omitempty"`
	Object         string     `json:"object,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	Description    string     `json:"description,omitempty"`
	Quantity       float64    `json:"quantity,omitempty"`
	Value          float64    `json:"value,omitempty,string"`
	Weight         float64    `json:"weight,omitempty"`
	HSTariffNumber string     `json:"hs_tariff_number,omitempty"`
	Code           string     `json:"code,omitempty"`
	OriginCountry  string     `json:"origin_country,omitempty"`
	Currency       string     `json:"currency,omitempty"`
}

type createCustomsInfoRequest struct {
	CustomsInfo *CustomsInfo `json:"customs_info,omitempty"`
}

// CreateCustomsInfo creates a new CustomsInfo object.
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
//					HSTarriffNumber: "654321",
//					OriginCountry:   "US",
//				},
//			},
//		},
//	)
func (c *Client) CreateCustomsInfo(in *CustomsInfo) (out *CustomsInfo, err error) {
	req := &createCustomsInfoRequest{CustomsInfo: in}
	err = c.post(context.Background(), "customs_infos", req, &out)
	return
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
	err = c.get(context.Background(), "customs_infos/"+customsInfoID, &out)
	return
}

// GetCustomsInfoWithContext performs the same operation as GetCustomsInfo, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetCustomsInfoWithContext(ctx context.Context, customsInfoID string) (out *CustomsInfo, err error) {
	err = c.get(ctx, "customs_infos/"+customsInfoID, &out)
	return
}

type createCustomsItemRequest struct {
	CustomsItem *CustomsItem `json:"customs_item,omitempty"`
}

// CreateCustomsItem creates a new CustomsItem object.
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
	req := &createCustomsItemRequest{CustomsItem: in}
	err = c.post(context.Background(), "customs_items", req, &out)
	return
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
	err = c.get(context.Background(), "customs_items/"+customsItemID, &out)
	return
}

// GetCustomsItemWithContext performs the same operation as GetCustomsItem, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetCustomsItemWithContext(ctx context.Context, customsItemID string) (out *CustomsItem, err error) {
	err = c.get(ctx, "customs_items/"+customsItemID, &out)
	return
}
