package easypost

import (
	"context"
	"net/http"
)

// A Parcel objects represent a physical container being shipped.
type Parcel struct {
	ID                string    `json:"id,omitempty" url:"id,omitempty"`
	Object            string    `json:"object,omitempty" url:"object,omitempty"`
	Mode              string    `json:"mode,omitempty" url:"mode,omitempty"`
	CreatedAt         *DateTime `json:"created_at,omitempty" url:"created_at,omitempty"`
	UpdatedAt         *DateTime `json:"updated_at,omitempty" url:"updated_at,omitempty"`
	Length            float64   `json:"length,omitempty" url:"length,omitempty"`
	Width             float64   `json:"width,omitempty" url:"width,omitempty"`
	Height            float64   `json:"height,omitempty" url:"height,omitempty"`
	PredefinedPackage string    `json:"predefined_package,omitempty" url:"predefined_package,omitempty"`
	Weight            float64   `json:"weight,omitempty" url:"weight,omitempty"`
}

type createParcelRequest struct {
	Parcel *Parcel `json:"parcel,omitempty" url:"parcel,omitempty"`
}

// CreateParcel creates a new Parcel object.
//
//	c := easypost.New(MyEasyPostAPIKey)
//	out, err := c.CreateParcel(
//		&easypost.Parcel{
//			Length: 20.2,
//			Width:  10.9,
//			Height: 5,
//			Weight: 65.9,
//		},
//	)
func (c *Client) CreateParcel(in *Parcel) (out *Parcel, err error) {
	return c.CreateParcelWithContext(context.Background(), in)
}

// CreateParcelWithContext performs the same operation as CreateParcel, but
// allows specifying a context that can interrupt the request.
func (c *Client) CreateParcelWithContext(ctx context.Context, in *Parcel) (out *Parcel, err error) {
	err = c.do(ctx, http.MethodPost, "parcels", &createParcelRequest{Parcel: in}, &out)
	return
}

// GetParcel retrieves an existing Parcel object by ID.
func (c *Client) GetParcel(parcelID string) (out *Parcel, err error) {
	return c.GetParcelWithContext(context.Background(), parcelID)
}

// GetParcelWithContext performs the same operation as GetParcel, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetParcelWithContext(ctx context.Context, parcelID string) (out *Parcel, err error) {
	err = c.do(ctx, http.MethodGet, "parcels/"+parcelID, nil, &out)
	return
}
