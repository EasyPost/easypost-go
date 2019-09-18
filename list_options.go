package easypost

import "time"

// ListOptions is used to specify query parameters for listing EasyPost objects.
type ListOptions struct {
	BeforeID      string     `json:"before_id,omitempty"`
	AfterID       string     `json:"after_id,omitempty"`
	StartDateTime *time.Time `json:"start_datetime,omitempty"`
	EndDateTime   *time.Time `json:"end_datetime,omitempty"`
	PageSize      int        `json:"page_size,omitempty"`
}
