package easypost

import (
	"context"
	"net/http"
)

// Report represents a CSV-formatted file that is a log of all the objects
// created within a certain time frame.
type Report struct {
	ID                string    `json:"id,omitempty"`
	Object            string    `json:"object,omitempty"`
	Mode              string    `json:"mode,omitempty"`
	CreatedAt         *DateTime `json:"created_at,omitempty"`
	UpdatedAt         *DateTime `json:"updated_at,omitempty"`
	Status            string    `json:"status,omitempty"`
	StartDate         string    `json:"start_date,omitempty"`
	EndDate           string    `json:"end_date,omitempty"`
	IncludeChildren   bool      `json:"include_children,omitempty"`
	URL               string    `json:"url,omitempty"`
	URLExpiresAt      *DateTime `json:"url_expires_at,omitempty"`
	SendEmail         bool      `json:"send_email,omitempty"`
	Columns           []string  `json:"columns,omitempty"`
	AdditionalColumns []string  `json:"additional_columns,omitempty"`
}

// ListReportsResult holds the results from the list reports API.
type ListReportsResult struct {
	Reports []*Report `json:"reports,omitempty"`
	Type    string    `json:"type,omitempty"`
	PaginatedCollection
}

// CreateReport generates a new report. Valid Fields for input are StartDate,
// EndDate and SendEmail. A new Report object is returned. Once the Status is
// available, the report can be downloaded from the provided URL for 30 seconds.
//
//	c := easypost.New(MyEasyPostAPIKey)
//	c.CreateReport(
//		"payment_log",
//		&easypost.Report{StartDate: "2016-10-01", EndDate: "2016-10-31"},
//	)
func (c *Client) CreateReport(typ string, in *Report) (out *Report, err error) {
	return c.CreateReportWithContext(context.Background(), typ, in)
}

// CreateReportWithContext performs the same operation as CreateReport, but
// allows specifying a context that can interrupt the request.
func (c *Client) CreateReportWithContext(ctx context.Context, typ string, in *Report) (out *Report, err error) {
	err = c.post(ctx, "reports/"+typ, in, &out)
	return
}

// ListReports provides a paginated result of Report objects of the given type.
func (c *Client) ListReports(typ string, opts *ListOptions) (out *ListReportsResult, err error) {
	return c.ListReportsWithContext(context.Background(), typ, opts)
}

// ListReportsWithContext performs the same operation as ListReports, but allows
// specifying a context that can interrupt the request.
func (c *Client) ListReportsWithContext(ctx context.Context, typ string, opts *ListOptions) (out *ListReportsResult, err error) {
	err = c.do(ctx, http.MethodGet, "reports/"+typ, c.convertOptsToURLValues(opts), &out)
	// Store the original query parameters for reuse when getting the next page
	out.Type = typ
	return
}

// GetNextReportPage returns the next page of reports
func (c *Client) GetNextReportPage(collection *ListReportsResult) (out *ListReportsResult, err error) {
	return c.GetNextReportPageWithContext(context.Background(), collection)
}

// GetNextReportPageWithPageSize returns the next page of reports with a specific page size
func (c *Client) GetNextReportPageWithPageSize(collection *ListReportsResult, pageSize int) (out *ListReportsResult, err error) {
	return c.GetNextReportPageWithPageSizeWithContext(context.Background(), collection, pageSize)
}

// GetNextReportPageWithContext performs the same operation as GetNextReportPage, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetNextReportPageWithContext(ctx context.Context, collection *ListReportsResult) (out *ListReportsResult, err error) {
	return c.GetNextReportPageWithPageSizeWithContext(ctx, collection, 0)
}

// GetNextReportPageWithPageSizeWithContext performs the same operation as GetNextReportPageWithPageSize, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetNextReportPageWithPageSizeWithContext(ctx context.Context, collection *ListReportsResult, pageSize int) (out *ListReportsResult, err error) {
	if len(collection.Reports) == 0 {
		err = EndOfPaginationError
		return
	}
	lastID := collection.Reports[len(collection.Reports)-1].ID
	params, err := nextPageParameters(collection.HasMore, lastID, pageSize)
	if err != nil {
		return
	}
	reportParams := &ListOptions{
		BeforeID: params.BeforeID,
	}
	if pageSize > 0 {
		reportParams.PageSize = pageSize
	}
	return c.ListReportsWithContext(ctx, collection.Type, reportParams)
}

// GetReport fetches a Report object by report type and ID.
func (c *Client) GetReport(typ, reportID string) (out *Report, err error) {
	return c.GetReportWithContext(context.Background(), typ, reportID)
}

// GetReportWithContext performs the same operation as GetReport, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetReportWithContext(ctx context.Context, typ, reportID string) (out *Report, err error) {
	err = c.get(ctx, "reports/"+typ+"/"+reportID, &out)
	return
}
