package easypost

import (
	"context"
	"time"
)

// A Parcel objects represent a physical container being shipped.
type Parcel struct {
	ID                string     `json:"id,omitempty"`
	Object            string     `json:"object,omitempty"`
	Mode              string     `json:"mode,omitempty"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
	Length            float64    `json:"length,omitempty"`
	Width             float64    `json:"width,omitempty"`
	Height            float64    `json:"height,omitempty"`
	PredefinedPackage string     `json:"predefined_package,omitempty"`
	Weight            float64    `json:"weight,omitempty"`
}

type createParcelRequest struct {
	Parcel *Parcel `json:"parcel,omitempty"`
}

// CreateParcel creates a new Parcel object.
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
	err = c.post(context.Background(), "parcels", &createParcelRequest{Parcel: in}, &out)
	return
}

// CreateParcelWithContext performs the same operation as CreateParcel, but
// allows specifying a context that can interrupt the request.
func (c *Client) CreateParcelWithContext(ctx context.Context, in *Parcel) (out *Parcel, err error) {
	err = c.post(ctx, "parcels", &createParcelRequest{Parcel: in}, &out)
	return
}

// GetParcel retrieves an existing Parcel object by ID.
func (c *Client) GetParcel(parcelID string) (out *Parcel, err error) {
	err = c.get(context.Background(), "parcels/"+parcelID, &out)
	return
}

// GetParcelWithContext performs the same operation as GetParcel, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetParcelWithContext(ctx context.Context, parcelID string) (out *Parcel, err error) {
	err = c.get(ctx, "parcels/"+parcelID, &out)
	return
}
