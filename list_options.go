package easypost

// ListOptions is used to specify query parameters for listing EasyPost objects.
type ListOptions struct {
	BeforeID      string    `url:"before_id,omitempty"`
	AfterID       string    `url:"after_id,omitempty"`
	StartDateTime *DateTime `url:"start_datetime,omitempty"`
	EndDateTime   *DateTime `url:"end_datetime,omitempty"`
	PageSize      int       `url:"page_size,omitempty"`
}

// nextPageParameters returns the next page of a paginated collection.
// If pageSize is 0, it will use the default page size
func nextPageParameters(hasMore bool, lastID string, pageSize int) (out *ListOptions, err error) {
	if !hasMore {
		err = newEndOfPaginationError()
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

// nextChildUserPageParameters returns the next page of a paginated collection of child users.
// If pageSize is 0, it will use the default page size
func nextChildUserPageParameters(hasMore bool, lastID string, pageSize int) (out *ListOptions, err error) {
	if !hasMore {
		err = newEndOfPaginationError()
		return
	}
	out = &ListOptions{
		AfterID: lastID,
	}
	if pageSize > 0 {
		out.PageSize = pageSize
	}
	return
}
