package easypost

import "time"

type Refund struct {
	ID                 string     `json:"id,omitempty"`
	Object             string     `json:"object,omitempty"`
	CreatedAt          *time.Time `json:"created_at,omitempty"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
	TrackingCode       string     `json:"tracking_code,omitempty"`
	ConfirmationNumber string     `json:"confirmation_number,omitempty"`
	Status             string     `json:"status,omitempty"`
	Carrier            string     `json:"carrier,omitempty"`
	ShipmentID         string     `json:"shipment_id,omitempty"`
}
