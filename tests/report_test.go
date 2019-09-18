package easypost_test

import (
	"testing"

	"github.com/EasyPost/easypost-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShipmentReport(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	report, err := TestClient.CreateReport(
		"shipment",
		&easypost.Report{
			StartDate: "2012-12-01",
			EndDate:   "2013-01-01",
		},
	)

	if err, ok := err.(*easypost.APIError); ok && err.Code == "REPORT.ALREADY_EXISTS" {
		reports, err := TestClient.ListReports(
			"shipment",
			&easypost.ListReportsOptions{
				StartDate: "2012-12-01",
				EndDate:   "2013-01-01",
			},
		)
		require.NoError(err)
		require.NotEmpty(reports.Reports)
		report = reports.Reports[0]
	} else {
		require.NoError(err)
	}

	assert.True(report.Status == "available" || report.Status == "new")

	report2, err := TestClient.GetReport("shipment", report.ID)
	require.NoError(err)
	assert.Equal(report.ID, report2.ID)

	reports, err := TestClient.ListReports("shipment", nil)
	require.NoError(err)
	assert.NotEmpty(reports.Reports)
	assert.Equal(report.ID, reports.Reports[0].ID)
}
