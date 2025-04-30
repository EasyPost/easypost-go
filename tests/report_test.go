package easypost_test

import (
	"errors"
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v5"
)

func (c *ClientTests) TestReportCreate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	report, err := client.CreateReport(
		c.fixture.ReportType(),
		&easypost.Report{
			StartDate: c.fixture.ReportDate(),
			EndDate:   c.fixture.ReportDate(),
		},
	)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Report{}), reflect.TypeOf(report))
	assert.True(strings.HasPrefix(report.ID, "shprep_"))
}

func (c *ClientTests) TestReportCustomColumnsCreate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	report, err := client.CreateReport(
		c.fixture.ReportType(),
		&easypost.Report{
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
	assert.Equal(reflect.TypeOf(&easypost.Report{}), reflect.TypeOf(report))
}

func (c *ClientTests) TestReportCustomAdditionalColumnsCreate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	report, err := client.CreateReport(
		c.fixture.ReportType(),
		&easypost.Report{
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
	assert.Equal(reflect.TypeOf(&easypost.Report{}), reflect.TypeOf(report))
}

func (c *ClientTests) TestReportRetrieve() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	report, err := client.CreateReport(
		c.fixture.ReportType(),
		&easypost.Report{
			StartDate: c.fixture.ReportDate(),
			EndDate:   c.fixture.ReportDate(),
		},
	)
	require.NoError(err)

	retrievedReport, err := client.GetReport(c.fixture.ReportType(), report.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Report{}), reflect.TypeOf(retrievedReport))
	assert.Equal(report.StartDate, retrievedReport.StartDate)
	assert.Equal(report.EndDate, retrievedReport.EndDate)
}

func (c *ClientTests) TestReportAll() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	reports, err := client.ListReports(
		c.fixture.ReportType(),
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	reportsList := reports.Reports

	assert.LessOrEqual(len(reportsList), c.fixture.pageSize())
	assert.NotNil(reports.HasMore)
	for _, report := range reportsList {
		assert.Equal(reflect.TypeOf(&easypost.Report{}), reflect.TypeOf(report))
	}
}

func (c *ClientTests) TestReportsGetNextPage() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	firstPage, err := client.ListReports(
		c.fixture.ReportType(),
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	nextPage, err := client.GetNextReportPageWithPageSize(firstPage, c.fixture.pageSize())
	defer func() {
		if err == nil {
			assert.True(len(nextPage.Reports) <= c.fixture.pageSize())

			lastIDOfFirstPage := firstPage.Reports[len(firstPage.Reports)-1].ID
			firstIdOfSecondPage := nextPage.Reports[0].ID

			assert.NotEqual(lastIDOfFirstPage, firstIdOfSecondPage)
		}
	}()
	if err != nil {
		var endOfPaginationErr *easypost.EndOfPaginationError
		if errors.As(err, &endOfPaginationErr) {
			assert.Equal(err.Error(), endOfPaginationErr.Error())
			return
		}
		require.NoError(err)
	}
}
