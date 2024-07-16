package easypost_test

import (
	"github.com/EasyPost/easypost-go/v4"
)

func (c *ClientTests) TestEstimateDeliveryDateForZipPair() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	carrier := c.fixture.USPS()

	params := &easypost.EstimateDeliveryDateForZipPairParams{
		FromZip:         c.fixture.CaAddress1().Zip,
		ToZip:           c.fixture.CaAddress2().Zip,
		Carriers:        []string{carrier},
		PlannedShipDate: c.fixture.PlannedShipDate(),
	}

	estimates, err := client.EstimateDeliveryDateForZipPair(params)
	require.NoError(err)

	assert.Equal(estimates.FromZip, params.FromZip)
	assert.Equal(estimates.ToZip, params.ToZip)
	assert.True(len(estimates.Results) > 0)
	for _, entry := range estimates.Results {
		assert.NotNil(entry.Carrier)
		assert.NotNil(entry.Service)
		assert.NotNil(entry.EasyPostTimeInTransitData.EasyPostEstimatedDeliveryDate)
		assert.NotNil(entry.EasyPostTimeInTransitData.DaysInTransit)
	}
}

func (c *ClientTests) TestRecommendShipDateForZipPair() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	carrier := c.fixture.USPS()

	params := &easypost.RecommendShipDateForZipPairParams{
		FromZip:             c.fixture.CaAddress1().Zip,
		ToZip:               c.fixture.CaAddress2().Zip,
		Carriers:            []string{carrier},
		DesiredDeliveryDate: c.fixture.DesiredDeliveryDate(),
	}

	recommendations, err := client.RecommendShipDateForZipPair(params)
	require.NoError(err)

	assert.Equal(recommendations.FromZip, params.FromZip)
	assert.Equal(recommendations.ToZip, params.ToZip)
	assert.True(len(recommendations.Results) > 0)
	for _, entry := range recommendations.Results {
		assert.NotNil(entry.Carrier)
		assert.NotNil(entry.Service)
		assert.NotNil(entry.EasyPostTimeInTransitData.EasyPostRecommendedShipDate)
		assert.NotNil(entry.EasyPostTimeInTransitData.DaysInTransit)
		assert.NotNil(entry.EasyPostTimeInTransitData.DeliveryDateConfidence)
		assert.NotNil(entry.EasyPostTimeInTransitData.EstimatedTransitDays)
	}

}
