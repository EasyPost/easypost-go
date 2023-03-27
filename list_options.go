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

func nextPageParameters(hasMore bool, lastId string) (out *ListOptions, err error) {
	if !hasMore {
		err = raiseEndOfPaginationError()
		return
	}
	return &ListOptions{
		AfterID: lastId,
	}, nil
}

func nextPageParametersWithPageSize(hasMore bool, lastId string, pageSize int) (out *ListOptions, err error) {
	params, err := nextPageParameters(hasMore, lastId)
	if pageSize > 0 {
		params.PageSize = pageSize
	}
	return params, err
}
