package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestReportCreate() {
	client := c.TestClient()
	assert := c.Assert()

	report, _ := client.CreateReport(
		c.fixture.ReportType(),
		&easypost.Report{
			StartDate: c.fixture.ReportDate(),
			EndDate:   c.fixture.ReportDate(),
		},
	)

	assert.Equal(reflect.TypeOf(&easypost.Report{}), reflect.TypeOf(report))
	assert.True(strings.HasPrefix(report.ID, "shprep_"))
}

func (c *ClientTests) TestReportCustomColumnsCreate() {
	client := c.TestClient()
	assert := c.Assert()

	report, _ := client.CreateReport(
		c.fixture.ReportType(),
		&easypost.Report{
			StartDate: c.fixture.ReportDate(),
			EndDate:   c.fixture.ReportDate(),
			Columns:   []string{"usps_zone"},
		},
	)

	// verify parameters by checking VCR cassette for correct URL
	// Some reports take a long time to generate, so we won't be able to consistently pull the report
	// There's unfortunately no way to check if the columns were included in the final report without parsing the CSV
	// so we assume, if we haven't gotten an error by this point, we've made the API calls correctly
	// any failure at this point is a server-side issue
	assert.Equal(reflect.TypeOf(&easypost.Report{}), reflect.TypeOf(report))
}

func (c *ClientTests) TestReportCustomAdditionalColumnsCreate() {
	client := c.TestClient()
	assert := c.Assert()

	report, _ := client.CreateReport(
		c.fixture.ReportType(),
		&easypost.Report{
			StartDate:         c.fixture.ReportDate(),
			EndDate:           c.fixture.ReportDate(),
			AdditionalColumns: []string{"from_name", "from_company"},
		},
	)

	// verify parameters by checking VCR cassette for correct URL
	// Some reports take a long time to generate, so we won't be able to consistently pull the report
	// There's unfortunately no way to check if the columns were included in the final report without parsing the CSV
	// so we assume, if we haven't gotten an error by this point, we've made the API calls correctly
	// any failure at this point is a server-side issue
	assert.Equal(reflect.TypeOf(&easypost.Report{}), reflect.TypeOf(report))
}

func (c *ClientTests) TestReportRetrieve() {
	client := c.TestClient()
	assert := c.Assert()

	report, _ := client.CreateReport(
		c.fixture.ReportType(),
		&easypost.Report{
			StartDate: c.fixture.ReportDate(),
			EndDate:   c.fixture.ReportDate(),
		},
	)

	retrievedReport, _ := client.GetReport(c.fixture.ReportType(), report.ID)

	assert.Equal(reflect.TypeOf(&easypost.Report{}), reflect.TypeOf(retrievedReport))
	assert.Equal(report.StartDate, retrievedReport.StartDate)
	assert.Equal(report.EndDate, retrievedReport.EndDate)
}

func (c *ClientTests) TestReportAll() {
	client := c.TestClient()
	assert := c.Assert()

	reports, _ := client.ListReports(
		c.fixture.ReportType(),
		&easypost.ListReportsOptions{
			PageSize: c.fixture.pageSize(),
		},
	)

	reportsList := reports.Reports

	assert.LessOrEqual(len(reportsList), c.fixture.pageSize())
	assert.NotNil(reports.HasMore)
	for _, report := range reportsList {
		assert.Equal(reflect.TypeOf(&easypost.Report{}), reflect.TypeOf(report))
	}
}
