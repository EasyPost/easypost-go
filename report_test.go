package easypost

import (
	"context"
	"errors"
	"net/url"
	"reflect"
	"strings"
	"time"
)

func (c *ClientTests) TestReportCreate() {
	client := c.testClient()
	assert, require := c.Assert(), c.Require()

	report, err := client.CreateReport(
		c.fixture.ReportType(),
		&Report{
			StartDate: c.fixture.ReportDate(),
			EndDate:   c.fixture.ReportDate(),
		},
	)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&Report{}), reflect.TypeOf(report))
	assert.True(strings.HasPrefix(report.ID, "shprep_"))
}

func (c *ClientTests) TestReportCustomColumnsCreate() {
	client := c.testClient()
	assert, require := c.Assert(), c.Require()

	report, err := client.CreateReport(
		c.fixture.ReportType(),
		&Report{
			StartDate: c.fixture.ReportDate(),
			EndDate:   c.fixture.ReportDate(),
			Columns:   []string{"usps_zone"},
		},
	)
	require.NoError(err)

	// Verify parameters by checking VCR cassette for correct URL.
	// Some reports take a long time to generate, so we won't be able to consistently pull the report.
	// There's unfortunately no way to check if the columns were included in the final report without parsing the CSV,
	// so we assume, if we haven't gotten an error by this point, we've made the API calls correctly
	// any failure at this point is a server-side issue
	assert.Equal(reflect.TypeOf(&Report{}), reflect.TypeOf(report))
}

func (c *ClientTests) TestReportCustomAdditionalColumnsCreate() {
	client := c.testClient()
	assert, require := c.Assert(), c.Require()

	report, err := client.CreateReport(
		c.fixture.ReportType(),
		&Report{
			StartDate:         c.fixture.ReportDate(),
			EndDate:           c.fixture.ReportDate(),
			AdditionalColumns: []string{"from_name", "from_company"},
		},
	)
	require.NoError(err)

	// Verify parameters by checking VCR cassette for correct URL.
	// Some reports take a long time to generate, so we won't be able to consistently pull the report.
	// There's unfortunately no way to check if the columns were included in the final report without parsing the CSV,
	// so we assume, if we haven't gotten an error by this point, we've made the API calls correctly
	// any failure at this point is a server-side issue
	assert.Equal(reflect.TypeOf(&Report{}), reflect.TypeOf(report))
}

func (c *ClientTests) TestReportRetrieve() {
	client := c.testClient()
	assert, require := c.Assert(), c.Require()

	report, err := client.CreateReport(
		c.fixture.ReportType(),
		&Report{
			StartDate: c.fixture.ReportDate(),
			EndDate:   c.fixture.ReportDate(),
		},
	)
	require.NoError(err)

	retrievedReport, err := client.GetReport(c.fixture.ReportType(), report.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&Report{}), reflect.TypeOf(retrievedReport))
	assert.Equal(report.StartDate, retrievedReport.StartDate)
	assert.Equal(report.EndDate, retrievedReport.EndDate)
}

func (c *ClientTests) TestReportGetNextPageKeepsDateFilters() {
	assert, require := c.Assert(), c.Require()

	client := New("not-a-real-key")
	client.MockRequests = []MockRequest{{
		MatchRule:    MockRequestMatchRule{Method: "GET", URLRegexPattern: "v2/reports/shipment"},
		ResponseInfo: MockRequestResponseInfo{StatusCode: 200, Body: `{"reports":[{"id":"shprep_1"}],"has_more":true}`},
	}}
	var query url.Values
	client.Hooks.AddRequestEventSubscriber(RequestHookEventSubscriber{
		Callback: func(_ context.Context, event RequestHookEvent) error {
			query = event.URL.Query()
			return nil
		},
	})

	start := NewDateTime(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := NewDateTime(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	firstPage, err := client.ListReports("shipment", &ListOptions{StartDateTime: &start, EndDateTime: &end})
	require.NoError(err)

	secondPage, err := client.GetNextReportPage(firstPage)
	require.NoError(err)
	assert.Equal("shprep_1", query.Get("before_id"))
	assert.Equal("2024-01-01T00:00:00Z", query.Get("start_datetime"))
	assert.Equal("2024-02-01T00:00:00Z", query.Get("end_datetime"))

	// The filters must still be there for the page after that
	query = nil
	_, err = client.GetNextReportPage(secondPage)
	require.NoError(err)
	assert.Equal("2024-01-01T00:00:00Z", query.Get("start_datetime"))
	assert.Equal("2024-02-01T00:00:00Z", query.Get("end_datetime"))
}

func (c *ClientTests) TestReportAll() {
	client := c.testClient()
	assert, require := c.Assert(), c.Require()

	reports, err := client.ListReports(
		c.fixture.ReportType(),
		&ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	reportsList := reports.Reports

	assert.LessOrEqual(len(reportsList), c.fixture.pageSize())
	assert.NotNil(reports.HasMore)
	for _, report := range reportsList {
		assert.Equal(reflect.TypeOf(&Report{}), reflect.TypeOf(report))
	}
}

func (c *ClientTests) TestReportsGetNextPage() {
	client := c.testClient()
	assert, require := c.Assert(), c.Require()

	firstPage, err := client.ListReports(
		c.fixture.ReportType(),
		&ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	nextPage, err := client.GetNextReportPageWithPageSize(firstPage, c.fixture.pageSize())
	defer func() {
		if err == nil {
			assert.True(len(nextPage.Reports) <= c.fixture.pageSize())

			lastIDOfFirstPage := firstPage.Reports[len(firstPage.Reports)-1].ID
			firstIDOfSecondPage := nextPage.Reports[0].ID

			assert.NotEqual(lastIDOfFirstPage, firstIDOfSecondPage)
		}
	}()
	if err != nil {
		var endOfPaginationErr *EndOfPaginationError
		if errors.As(err, &endOfPaginationErr) {
			assert.Equal(err.Error(), endOfPaginationErr.Error())
			return
		}
		require.NoError(err)
	}
}
