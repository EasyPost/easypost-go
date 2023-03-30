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

// nextPageParameters returns the next page of a paginated collection.
// If pageSize is 0, it will use the default page size
func nextPageParameters(hasMore bool, lastID string, pageSize int) (out *ListOptions, err error) {
	if !hasMore {
		err = EndOfPaginationError
		return
	}
	out = &ListOptions{
		BeforeID: lastID,
	}
	if pageSize > 0 {
		out.PageSize = pageSize
	}
	return
}
