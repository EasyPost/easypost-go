package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestReportPaymentLogCreate() {
	client := c.TestClient()
	assert := c.Assert()

	report, _ := client.CreateReport(
		"payment_log",
		&easypost.Report{
			StartDate: c.fixture.ReportStartDate(),
			EndDate:   c.fixture.ReportEndDate(),
		},
	)

	assert.Equal(reflect.TypeOf(&easypost.Report{}), reflect.TypeOf(report))
	assert.True(strings.HasPrefix(report.ID, "plrep_"))
}

func (c *ClientTests) TestReportRefundCreate() {
	client := c.TestClient()
	assert := c.Assert()

	report, _ := client.CreateReport(
		"refund",
		&easypost.Report{
			StartDate: c.fixture.ReportStartDate(),
			EndDate:   c.fixture.ReportEndDate(),
		},
	)

	assert.Equal(reflect.TypeOf(&easypost.Report{}), reflect.TypeOf(report))
	assert.True(strings.HasPrefix(report.ID, "refrep_"))
}

func (c *ClientTests) TestReportShipmentLogCreate() {
	client := c.TestClient()
	assert := c.Assert()

	report, _ := client.CreateReport(
		"shipment",
		&easypost.Report{
			StartDate: c.fixture.ReportStartDate(),
			EndDate:   c.fixture.ReportEndDate(),
		},
	)

	assert.Equal(reflect.TypeOf(&easypost.Report{}), reflect.TypeOf(report))
	assert.True(strings.HasPrefix(report.ID, "shprep_"))
}

func (c *ClientTests) TestReportShipmentInvoiceCreate() {
	client := c.TestClient()
	assert := c.Assert()

	report, _ := client.CreateReport(
		"shipment_invoice",
		&easypost.Report{
			StartDate: c.fixture.ReportStartDate(),
			EndDate:   c.fixture.ReportEndDate(),
		},
	)

	assert.Equal(reflect.TypeOf(&easypost.Report{}), reflect.TypeOf(report))
	assert.True(strings.HasPrefix(report.ID, "shpinvrep_"))
}

func (c *ClientTests) TestReportTrackerCreate() {
	client := c.TestClient()
	assert := c.Assert()

	report, _ := client.CreateReport(
		"tracker",
		&easypost.Report{
			StartDate: c.fixture.ReportStartDate(),
			EndDate:   c.fixture.ReportEndDate(),
		},
	)

	assert.Equal(reflect.TypeOf(&easypost.Report{}), reflect.TypeOf(report))
	assert.True(strings.HasPrefix(report.ID, "trkrep_"))
}

func (c *ClientTests) TestReportPaymentLogRetrieve() {
	client := c.TestClient()
	assert := c.Assert()

	report, _ := client.CreateReport(
		"payment_log",
		&easypost.Report{
			StartDate: c.fixture.ReportStartDate(),
			EndDate:   c.fixture.ReportEndDate(),
		},
	)

	retrievedReport, _ := client.GetReport("payment_log", report.ID)

	assert.Equal(reflect.TypeOf(&easypost.Report{}), reflect.TypeOf(retrievedReport))
	assert.Equal(report.StartDate, retrievedReport.StartDate)
	assert.Equal(report.EndDate, retrievedReport.EndDate)
}

func (c *ClientTests) TestReportRefundLogRetrieve() {
	client := c.TestClient()
	assert := c.Assert()

	report, _ := client.CreateReport(
		"refund",
		&easypost.Report{
			StartDate: c.fixture.ReportStartDate(),
			EndDate:   c.fixture.ReportEndDate(),
		},
	)

	retrievedReport, _ := client.GetReport("refund", report.ID)

	assert.Equal(reflect.TypeOf(&easypost.Report{}), reflect.TypeOf(retrievedReport))
	assert.Equal(report.StartDate, retrievedReport.StartDate)
	assert.Equal(report.EndDate, retrievedReport.EndDate)
}

func (c *ClientTests) TestReportShipmentLogRetrieve() {
	client := c.TestClient()
	assert := c.Assert()

	report, _ := client.CreateReport(
		"shipment",
		&easypost.Report{
			StartDate: c.fixture.ReportStartDate(),
			EndDate:   c.fixture.ReportEndDate(),
		},
	)

	retrievedReport, _ := client.GetReport("shipment", report.ID)

	assert.Equal(reflect.TypeOf(&easypost.Report{}), reflect.TypeOf(retrievedReport))
	assert.Equal(report.StartDate, retrievedReport.StartDate)
	assert.Equal(report.EndDate, retrievedReport.EndDate)
}

func (c *ClientTests) TestReportShipmentInvoiceLogRetrieve() {
	client := c.TestClient()
	assert := c.Assert()

	report, _ := client.CreateReport(
		"shipment_invoice",
		&easypost.Report{
			StartDate: c.fixture.ReportStartDate(),
			EndDate:   c.fixture.ReportEndDate(),
		},
	)

	retrievedReport, _ := client.GetReport("shipment_invoice", report.ID)

	assert.Equal(reflect.TypeOf(&easypost.Report{}), reflect.TypeOf(retrievedReport))
	assert.Equal(report.StartDate, retrievedReport.StartDate)
	assert.Equal(report.EndDate, retrievedReport.EndDate)
}

func (c *ClientTests) TestReportTrackerLogRetrieve() {
	client := c.TestClient()
	assert := c.Assert()

	report, _ := client.CreateReport(
		"tracker",
		&easypost.Report{
			StartDate: c.fixture.ReportStartDate(),
			EndDate:   c.fixture.ReportEndDate(),
		},
	)

	retrievedReport, _ := client.GetReport("tracker", report.ID)

	assert.Equal(reflect.TypeOf(&easypost.Report{}), reflect.TypeOf(retrievedReport))
	assert.Equal(report.StartDate, retrievedReport.StartDate)
	assert.Equal(report.EndDate, retrievedReport.EndDate)
}

func (c *ClientTests) TestReportAll() {
	client := c.TestClient()
	assert := c.Assert()

	reports, _ := client.ListReports(
		"shipment",
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

func (c *ClientTests) TestReportCreateNoType() {
	client := c.TestClient()
	require := c.Require()

	_, err := client.CreateReport("", nil)

	require.Error(err)
}

func (c *ClientTests) TestReportAllNoType() {
	client := c.TestClient()
	require := c.Require()

	_, err := client.ListReports(
		"",
		&easypost.ListReportsOptions{
			PageSize: c.fixture.pageSize(),
		},
	)

	require.Error(err)
}
