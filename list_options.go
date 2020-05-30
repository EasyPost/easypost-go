package easypost

import "time"

// ListOptions is used to specify query parameters for listing EasyPost objects.
type ListOptions struct {
	BeforeID      string     `url:"before_id,omitempty"`
	AfterID       string     `url:"after_id,omitempty"`
	StartDateTime *time.Time `url:"start_datetime,omitempty"`
	EndDateTime   *time.Time `url:"end_datetime,omitempty"`
	PageSize      int        `url:"page_size,omitempty"`
}
