package easypost

import (
	"context"
	"net/http"
)

// LumaPromiseResponse represents the response from the Luma AI Promise API.
type LumaPromiseResponse struct {
	LumaInfo *LumaInfo `json:"luma_info,omitempty" url:"luma_info,omitempty"`
}

// LumaInfo contains information about the Luma AI results for a shipment.
type LumaInfo struct {
	AiResults          []*AiResults `json:"ai_results,omitempty" url:"ai_results,omitempty"`
	MatchingRuleIdx    int          `json:"matching_rule_idx,omitempty" url:"matching_rule_idx,omitempty"`
	RulesetDescription string       `json:"ruleset_description,omitempty" url:"ruleset_description,omitempty"`
	LumaSelectedRate   *Rate        `json:"luma_selected_rate,omitempty" url:"luma_selected_rate,omitempty"`
}

// AiResults contains the results of the AI analysis for a shipment.
type AiResults struct {
	Carrier                  string    `json:"carrier,omitempty" url:"carrier,omitempty"`
	MeetsRulesetRequirements bool      `json:"meets_ruleset_requirements,omitempty" url:"meets_ruleset_requirements,omitempty"`
	PredictedDeliverByDate   *DateTime `json:"predicted_delivery_date,omitempty" url:"predicted_delivery_date,omitempty"`
	PredictedDeliverDays     int       `json:"predicted_delivery_days,omitempty" url:"predicted_delivery_days,omitempty"`
	RateID                   string    `json:"rate_id,omitempty" url:"rate_id,omitempty"`
	RateUSD                  string    `json:"rate_usd,omitempty" url:"rate_usd,omitempty"`
	Service                  string    `json:"service,omitempty" url:"service,omitempty"`
}

// WrappedLumaRequest is the request structure for retrieving Luma AI Promise.
type WrappedLumaRequest struct {
	Shipment *LumaRequest `json:"shipment,omitempty" url:"shipment,omitempty"`
}

// LumaRequest represents the request parameters for the Luma AI Promise API.
type LumaRequest struct {
	Shipment
	RulesetName     string `json:"ruleset_name,omitempty" url:"ruleset_name,omitempty"`
	PlannedShipDate string `json:"planned_ship_date,omitempty" url:"planned_ship_date,omitempty"`
	DeliverByDate   string `json:"deliver_by_date,omitempty" url:"deliver_by_date,omitempty"`
	PersistLabel    bool   `json:"persist_label,omitempty" url:"persist_label,omitempty"`
}

// GetLumaPromise retrieves the Luma AI Promise for a given shipment.
func (c *Client) GetLumaPromise(in *LumaRequest) (out *LumaInfo, err error) {
	return c.GetLumaPromiseWithContext(context.Background(), in)
}

// GetLumaPromiseWithContext performs the same operation as GetLumaPromise, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetLumaPromiseWithContext(ctx context.Context, in *LumaRequest) (out *LumaInfo, err error) {
	res := struct {
		LumaInfo *LumaInfo `json:"luma_info,omitempty" url:"luma_info,omitempty"`
	}{LumaInfo: out}
	req := &WrappedLumaRequest{Shipment: in}
	err = c.do(ctx, http.MethodPost, "luma/promise", &req, &res)
	out = res.LumaInfo
	return
}
