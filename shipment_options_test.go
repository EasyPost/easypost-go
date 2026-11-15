package easypost

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/google/go-querystring/query"
)

// A letter marked not machinable must send machinable=false. encoding/json and
// go-querystring both treat a false bool as empty when the tag says omitempty,
// so the field never leaves the client and USPS rates the letter as machinable.
func TestShipmentOptionsMachinableFalseIsEncoded(t *testing.T) {
	var opts ShipmentOptions
	if err := json.Unmarshal([]byte(`{"machinable":false}`), &opts); err != nil {
		t.Fatalf("unmarshal options: %v", err)
	}

	body, err := json.Marshal(opts)
	if err != nil {
		t.Fatalf("marshal options: %v", err)
	}
	if !bytes.Contains(body, []byte(`"machinable":false`)) {
		t.Fatalf("machinable false omitted from options JSON: %s", body)
	}

	values, err := query.Values(opts)
	if err != nil {
		t.Fatalf("encode options: %v", err)
	}
	if values.Get("machinable") != "false" {
		t.Fatalf("machinable false omitted from options query: %s", values.Encode())
	}

	// CreateShipment posts this wrapper (client.setParameters json.Marshal).
	var shipment Shipment
	if err := json.Unmarshal([]byte(`{"options":{"machinable":false}}`), &shipment); err != nil {
		t.Fatalf("unmarshal shipment: %v", err)
	}
	reqBody, err := json.Marshal(createShipmentRequest{Shipment: &shipment})
	if err != nil {
		t.Fatalf("marshal create shipment: %v", err)
	}
	if !bytes.Contains(reqBody, []byte(`"machinable":false`)) {
		t.Fatalf("machinable false omitted from create-shipment JSON: %s", reqBody)
	}

	trueBody, err := json.Marshal(mustShipmentOptions(t, `{"machinable":true}`))
	if err != nil {
		t.Fatalf("marshal true: %v", err)
	}
	if !bytes.Contains(trueBody, []byte(`"machinable":true`)) {
		t.Fatalf("machinable true omitted: %s", trueBody)
	}

	unsetBody, err := json.Marshal(ShipmentOptions{})
	if err != nil {
		t.Fatalf("marshal unset: %v", err)
	}
	if bytes.Contains(unsetBody, []byte(`machinable`)) {
		t.Fatalf("unset machinable should stay omitted, got %s", unsetBody)
	}
}

func mustShipmentOptions(t *testing.T, raw string) ShipmentOptions {
	t.Helper()
	var opts ShipmentOptions
	if err := json.Unmarshal([]byte(raw), &opts); err != nil {
		t.Fatalf("unmarshal %s: %v", raw, err)
	}
	return opts
}
