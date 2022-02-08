package easypost

import (
	"context"
	"net/http"
	"time"
)

// Report represents a CSV-formatted file that is a log of all the objects
// created within a certain time frame.
type Report struct {
	ID              string     `json:"id,omitempty"`
	Object          string     `json:"object,omitempty"`
	Mode            string     `json:"mode,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
	Status          string     `json:"status,omitempty"`
	StartDate       string     `json:"start_date,omitempty"`
	EndDate         string     `json:"end_date,omitempty"`
	IncludeChildren bool       `json:"include_children,omitempty"`
	URL             string     `json:"url,omitempty"`
	URLExpiresAt    *time.Time `json:"url_expires_at,omitempty"`
	SendEmail       bool       `json:"send_email,omitempty"`
}

// CreateReport generates a new report. Valid Fields for input are StartDate,
// EndDate and SendEmail. A new Report object is returned. Once the Status is
// available, the report can be downloded from the provided URL for 30 seconds.
//	c := easypost.New(MyEasyPostAPIKey)
//	c.CreateReport(
//		"payment_log",
//		&easypost.Report{StartDate: "2016-10-01", EndDate: "2016-10-31"},
//	)
func (c *Client) CreateReport(typ string, in *Report) (out *Report, err error) {
	err = c.post(context.Background(), "reports/"+typ, in, &out)
	return
}

// CreateReportWithContext performs the same operation as CreateReport, but
// allows specifying a context that can interrupt the request.
func (c *Client) CreateReportWithContext(ctx context.Context, typ string, in *Report) (out *Report, err error) {
	err = c.post(ctx, "reports/"+typ, in, &out)
	return
}

// ListReportsOptions is used to specify query parameters for listing Report
// objects.
type ListReportsOptions struct {
	BeforeID  string `url:"before_id,omitempty"`
	AfterID   string `url:"after_id,omitempty"`
	StartDate string `url:"start_datetime,omitempty"`
	EndDate   string `url:"end_datetime,omitempty"`
	PageSize  int    `url:"page_size,omitempty"`
}

// ListReportsResult holds the results from the list reports API.
type ListReportsResult struct {
	Reports []*Report `json:"reports,omitempty"`
	// HasMore indicates if there are more responses to be fetched. If True,
	// additional responses can be fetched by updating the ListReportsOptions
	// parameter's AfterID field with the ID of the last item in this object's
	// Reports field.
	HasMore bool `json:"has_more,omitempty"`
}

// ListReports provides a paginated result of Report objects of the given type.
func (c *Client) ListReports(typ string, opts *ListReportsOptions) (out *ListReportsResult, err error) {
	return c.ListReportsWithContext(context.Background(), typ, opts)
}

// ListReportsWithContext performs the same operation as ListReports, but allows
// specifying a context that can interrupt the request.
func (c *Client) ListReportsWithContext(ctx context.Context, typ string, opts *ListReportsOptions) (out *ListReportsResult, err error) {
	err = c.do(ctx, http.MethodGet, "reports/"+typ, c.convertOptsToURLValues(opts), &out)
	return
}

// GetReport fetches a Report object by report type and ID.
func (c *Client) GetReport(typ, reportID string) (out *Report, err error) {
	err = c.get(context.Background(), "reports/"+typ+"/"+reportID, &out)
	return
}

// GetReportWithContext performs the same operation as GetReport, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetReportWithContext(ctx context.Context, typ, reportID string) (out *Report, err error) {
	err = c.get(ctx, "reports/"+typ+"/"+reportID, &out)
	return
}
