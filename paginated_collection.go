package easypost

// PaginatedCollection holds the results of a paginated API response.
type PaginatedCollection struct {
	// HasMore indicates if there are more responses to be fetched. If True,
	// additional responses can be fetched by updating the AfterID parameter
	// with the ID of the last item in this object's collection list.
	HasMore bool `json:"has_more,omitempty" url:"has_more,omitempty"`
}
