package easypost_test

import (
	"github.com/EasyPost/easypost-go/v4"
)

func (c *ClientTests) TestEstimateDeliveryDateForZipPair() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	carrier := c.fixture.USPS()

	params := &easypost.EstimateDeliveryDateForZipPairParams{
		OriginPostalCode:      c.fixture.CaAddress1().Zip,
		DestinationPostalCode: c.fixture.CaAddress2().Zip,
		Carriers:              []string{carrier},
		PlannedShipDate:       c.fixture.PlannedShipDate(),
	}

	estimates, err := client.EstimateDeliveryDateForZipPair(params)
	require.NoError(err)

	assert.Equal(estimates.OriginPostalCode, params.OriginPostalCode)
	assert.Equal(estimates.DestinationPostalCode, params.DestinationPostalCode)
	assert.True(len(estimates.Estimates) > 0)
	for _, entry := range estimates.Estimates {
		assert.NotNil(entry.Carrier)
		assert.NotNil(entry.Service)
		assert.NotNil(entry.EasyPostTimeInTransitData.EasyPostEstimatedDeliveryDate)
		assert.NotNil(entry.EasyPostTimeInTransitData.TimeInTransitPercentiles)
	}
}

func (c *ClientTests) TestRecommendShipDateForZipPair() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	carrier := c.fixture.USPS()

	params := &easypost.RecommendShipDateForZipPairParams{
		OriginPostalCode:      c.fixture.CaAddress1().Zip,
		DestinationPostalCode: c.fixture.CaAddress2().Zip,
		Carriers:              []string{carrier},
		DesiredDeliveryDate:   c.fixture.DesiredDeliveryDate(),
	}

	recommendations, err := client.RecommendShipDateForZipPair(params)
	require.NoError(err)

	assert.Equal(recommendations.OriginPostalCode, params.OriginPostalCode)
	assert.Equal(recommendations.DestinationPostalCode, params.DestinationPostalCode)
	assert.True(len(recommendations.Estimates) > 0)
	for _, entry := range recommendations.Estimates {
		assert.NotNil(entry.Carrier)
		assert.NotNil(entry.Service)
		assert.NotNil(entry.EasyPostTimeInTransitData.EasyPostRecommendedShipDate)
		assert.NotNil(entry.EasyPostTimeInTransitData.TimeInTransitPercentiles)
		assert.NotNil(entry.EasyPostTimeInTransitData.DeliveryDateConfidence)
		assert.NotNil(entry.EasyPostTimeInTransitData.EstimatedTransitDays)
	}

}
