package easypost

import (
	"encoding/json"
)

// UnmarshalJSONObject attempts to unmarshal an easypost object from JSON data.
// An error is only returned if data is non-zero length and JSON decoding fails.
// If data contains a valid JSON object with an "object" key/value that matches
// a known object type, it will return a pointer to the corresponding object
// type from this package, e.g. *Address, *Batch, *Refund, etc. If the decoded
// JSON doesn't match, the return value will be the default type stored by
// json.Unmarshal.
func UnmarshalJSONObject(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, nil
	}

	var obj struct {
		Object string `json:"object"`
	}

	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	var result interface{}
	switch obj.Object {
	case "Address":
		result = new(Address)
	case "ApiKey":
		result = new(APIKey)
	case "Batch":
		result = new(Batch)
	case "CarrierAccount":
		result = new(CarrierAccount)
	case "CarrierDetail":
		result = new(TrackingCarrierDetail)
	case "CarrierType":
		result = new(CarrierType)
	case "CustomsInfo":
		result = new(CustomsInfo)
	case "CustomsItem":
		result = new(CustomsItem)
	case "Event":
		result = new(Event)
	case "Form":
		result = new(Form)
	case "Insurance":
		result = new(Insurance)
	case "Order":
		result = new(Order)
	case "Parcel":
		result = new(Parcel)
	case "Payload":
		result = new(EventPayload)
	case "PaymentLog":
		result = new(PaymentLog)
	case "PaymentLogReport":
		result = new(Report)
	case "Pickup":
		result = new(Pickup)
	case "PostageLabel":
		result = new(PostageLabel)
	case "Rate":
		result = new(Rate)
	case "Refund":
		result = new(Refund)
	case "RefundReport":
		result = new(Report)
	case "ScanForm":
		result = new(ScanForm)
	case "Shipment":
		result = new(Shipment)
	case "ShipmentInvoiceReport":
		result = new(Report)
	case "ShipmentReport":
		result = new(Report)
	case "Tracker":
		result = new(Tracker)
	case "TrackerReport":
		result = new(Report)
	case "TrackingDetail":
		result = new(TrackingDetail)
	case "TrackingLocation":
		result = new(TrackingLocation)
	case "User":
		result = new(User)
	case "Webhook":
		result = new(Webhook)
	}

	// If 'obj' didn't match one of the types above, result will be a nil
	// interface type. This will cause Unmarshal to use the default Go type for
	// the JSON value, like map[string]interface{}.
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}
