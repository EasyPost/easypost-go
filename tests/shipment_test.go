package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v4"
)

func (c *ClientTests) TestShipmentCreate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.FullShipment())
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Shipment{}), reflect.TypeOf(shipment))
	assert.True(strings.HasPrefix(shipment.ID, "shp_"))
	assert.NotNil(shipment.Rates)
	assert.Equal("PNG", shipment.Options.LabelFormat)
	assert.Equal("123", shipment.Options.InvoiceNumber)
	assert.Equal("123", shipment.Reference)
}

func (c *ClientTests) TestShipmentRetrieve() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.FullShipment())
	require.NoError(err)

	retrievedShipment, err := client.GetShipment(shipment.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Shipment{}), reflect.TypeOf(retrievedShipment))
	assert.Equal(shipment, retrievedShipment)
}

func (c *ClientTests) TestShipmentAll() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipments, err := client.ListShipments(
		&easypost.ListShipmentsOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	shipmentsList := shipments.Shipments

	assert.LessOrEqual(len(shipmentsList), c.fixture.pageSize())
	assert.NotNil(shipments.HasMore)
	for _, shipment := range shipmentsList {
		assert.Equal(reflect.TypeOf(&easypost.Shipment{}), reflect.TypeOf(shipment))
	}
}

func (c *ClientTests) TestShipmentBuy() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.FullShipment())
	require.NoError(err)

	lowestRate, err := client.LowestShipmentRate(shipment)
	require.NoError(err)

	boughtShipment, err := client.BuyShipment(shipment.ID, &lowestRate, "")
	require.NoError(err)

	assert.NotNil(boughtShipment.PostageLabel)
}

func (c *ClientTests) TestShipmentRegenerateRates() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.FullShipment())
	require.NoError(err)

	rates, err := client.RerateShipment(shipment.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf([]*easypost.Rate{}), reflect.TypeOf(rates))
	for _, rate := range rates {
		assert.Equal(reflect.TypeOf(&easypost.Rate{}), reflect.TypeOf(rate))
	}
}

func (c *ClientTests) TestShipmentConvertLabel() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	require.NoError(err)

	shipmentWithNewLabel, err := client.GetShipmentLabel(shipment.ID, "ZPL")
	require.NoError(err)

	assert.NotNil(shipmentWithNewLabel.PostageLabel.LabelZPLURL)
}

// If the shipment was purchased with a USPS rate, it must have had its insurance set to `0` when bought
// so that USPS doesn't automatically insure it. So we could manually insure it here.
func (c *ClientTests) TestShipmentInsure() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipmentData := c.fixture.OneCallBuyShipment()
	// Set to 0 so USPS doesn't insure this automatically and we can insure the shipment manually
	shipmentData.Insurance = "0"

	shipment, err := client.CreateShipment(shipmentData)
	require.NoError(err)

	insuredShipment, err := client.InsureShipment(shipment.ID, "100")
	require.NoError(err)

	assert.Equal("100.00", insuredShipment.Insurance)
}

// Refunding a test shipment must happen within seconds of the shipment being created as test shipments naturally
// follow a flow of created -> delivered to cycle through tracking events in test mode - as such anything older
// than a few seconds in test mode may not be refundable.
func (c *ClientTests) TestShipmentRefund() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	require.NoError(err)

	refundShipment, err := client.RefundShipment(shipment.ID)
	require.NoError(err)

	assert.Equal("submitted", refundShipment.RefundStatus)
}

func (c *ClientTests) TestShipmentSmartrate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.BasicShipment())
	require.NoError(err)

	assert.NotNil(shipment.Rates)

	smartrates, err := client.GetShipmentSmartrates(shipment.ID)
	require.NoError(err)

	assert.Equal(shipment.Rates[0].ID, smartrates[0].ID)
	assert.NotNil(smartrates[0].TimeInTransit.Percentile50)
	assert.NotNil(smartrates[0].TimeInTransit.Percentile75)
	assert.NotNil(smartrates[0].TimeInTransit.Percentile85)
	assert.NotNil(smartrates[0].TimeInTransit.Percentile90)
	assert.NotNil(smartrates[0].TimeInTransit.Percentile95)
	assert.NotNil(smartrates[0].TimeInTransit.Percentile97)
	assert.NotNil(smartrates[0].TimeInTransit.Percentile99)
}

func (c *ClientTests) TestShipmentCreateEmptyObjects() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipmentData := c.fixture.BasicShipment()
	shipmentData.Options = nil
	shipmentData.TaxIdentifiers = nil
	shipmentData.Reference = ""
	shipmentData.CustomsInfo = &easypost.CustomsInfo{}
	shipmentData.CustomsInfo.CustomsItems = []*easypost.CustomsItem{}

	shipment, err := client.CreateShipment(shipmentData)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Shipment{}), reflect.TypeOf(shipment))
	assert.True(strings.HasPrefix(shipment.ID, "shp_"))
	assert.NotNil(shipment.Options) // The EasyPost API populates some default values here
	assert.Nil(shipment.CustomsInfo)
	assert.Nil(shipment.TaxIdentifiers)
	assert.NotNil(shipment.Reference)
}

func (c *ClientTests) TestShipmentCreateTaxIdentifier() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipmentData := c.fixture.BasicShipment()
	shipmentData.TaxIdentifiers = []*easypost.TaxIdentifier{c.fixture.TaxIdentifier()}

	shipment, err := client.CreateShipment(shipmentData)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Shipment{}), reflect.TypeOf(shipment))
	assert.True(strings.HasPrefix(shipment.ID, "shp_"))
	assert.Equal("IOSS", shipment.TaxIdentifiers[0].TaxIdType)
}

func (c *ClientTests) TestShipmentCreateWithIds() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	fromAddress, err := client.CreateAddress(c.fixture.CaAddress1(), nil)
	require.NoError(err)
	toAddress, err := client.CreateAddress(c.fixture.CaAddress2(), nil)
	require.NoError(err)
	parcel, err := client.CreateParcel(c.fixture.BasicParcel())
	require.NoError(err)

	shipment, err := client.CreateShipment(
		&easypost.Shipment{
			FromAddress: &easypost.Address{ID: fromAddress.ID},
			ToAddress:   &easypost.Address{ID: toAddress.ID},
			Parcel:      &easypost.Parcel{ID: parcel.ID},
		},
	)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Shipment{}), reflect.TypeOf(shipment))
	assert.True(strings.HasPrefix(shipment.ID, "shp_"))
	assert.True(strings.HasPrefix(shipment.FromAddress.ID, "adr_"))
	assert.True(strings.HasPrefix(shipment.ToAddress.ID, "adr_"))
	assert.True(strings.HasPrefix(shipment.Parcel.ID, "prcl_"))
	assert.Equal("388 Townsend St", shipment.FromAddress.Street1)
}

func (c *ClientTests) TestShipmentInstanceLowestSmartRate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.BasicShipment())
	require.NoError(err)
	smartrate, err := client.LowestSmartrate(shipment, 3, "percentile_90")
	require.NoError(err)

	// Test lowest smartrate with valid filters
	assert.Equal("GroundAdvantage", smartrate.Service)
	assert.Equal(5.93, smartrate.Rate)
	assert.Equal("USPS", smartrate.Carrier)

	// Test lowest smartrate with invalid filters (should error due to strict delivery days)
	smartrate, err = client.LowestSmartrate(shipment, 0, "percentile_90")
	assert.Error(err)

	// Test lowest smartrate with invalid filters (should error due to invalid delivery accuracy)
	smartrate, err = client.LowestSmartrate(shipment, 1, "BAD_ACCURACY")
	assert.Error(err)
}

func (c *ClientTests) TestShipmentLowestRate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.FullShipment())
	require.NoError(err)

	// Test lowest rate with no filters
	lowestRate, err := client.LowestShipmentRate(shipment)
	require.NoError(err)
	assert.Equal("GroundAdvantage", lowestRate.Service)
	assert.Equal("5.93", lowestRate.Rate)
	assert.Equal("USPS", lowestRate.Carrier)

	// Test lowest rate with service filter (this rate is higher than the lowest but should filter)
	lowestRate, err = client.LowestShipmentRateWithCarrierAndService(shipment, nil, []string{"Priority"})
	require.NoError(err)
	assert.Equal("Priority", lowestRate.Service)
	assert.Equal("6.90", lowestRate.Rate)
	assert.Equal("USPS", lowestRate.Carrier)

	// Test lowest rate with carrier filter (should error due to bad carrier)
	lowestRate, err = client.LowestShipmentRateWithCarrier(shipment, []string{"BAD_CARRIER"})
	assert.Error(err)
}

func (c *ClientTests) TestShipmentGenerateForm() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	require.NoError(err)

	formType := "return_packing_slip"

	// Test shipment form with no additional options
	shipmentWithForm, err := client.GenerateShipmentForm(shipment.ID, formType)
	require.NoError(err)
	assert.Len(shipmentWithForm.Forms, 1)
	form := shipmentWithForm.Forms[0]
	assert.Equal(formType, form.FormType)
	assert.NotNil(form.FormURL)

	// Test shipment form with additional options
	shipmentWithForm, err = client.GenerateShipmentFormWithOptions(shipment.ID, formType, c.fixture.RmaFormOptions())
	require.NoError(err)
	assert.Len(shipmentWithForm.Forms, 1)
	form = shipmentWithForm.Forms[0]
	assert.Equal(formType, form.FormType)
	assert.NotNil(form.FormURL)
}

func (c *ClientTests) TestBuyShipmentWithEndShipper() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	endShipper, err := client.CreateEndShipper(c.fixture.CaAddress1())
	require.NoError(err)

	shipment, err := client.CreateShipment(c.fixture.BasicShipment())
	require.NoError(err)

	lowestRate, err := client.LowestShipmentRate(shipment)
	require.NoError(err)

	boughtShipment, err := client.BuyShipmentWithEndShipper(shipment.ID, &lowestRate, "", endShipper.ID)
	require.NoError(err)

	assert.NotNil(boughtShipment.PostageLabel)
}

func (c *ClientTests) TestShipmentsGetNextPage() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	firstPage, err := client.ListShipments(
		&easypost.ListShipmentsOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	nextPage, err := client.GetNextShipmentPageWithPageSize(firstPage, c.fixture.pageSize())
	defer func() {
		if err == nil {
			assert.True(len(nextPage.Shipments) <= c.fixture.pageSize())

			lastIDOfFirstPage := firstPage.Shipments[len(firstPage.Shipments)-1].ID
			firstIdOfSecondPage := nextPage.Shipments[0].ID

			assert.NotEqual(lastIDOfFirstPage, firstIdOfSecondPage)
		}
	}()
	if err != nil {
		assert.Equal(err.Error(), easypost.NoPagesLeftToRetrieve)
		return
	}
}

// Tests that we retrieve time-in-transit data for each of the Rates of a Shipment.
func (c *ClientTests) TestShipmentGetShipmentEstimatedDeliveryDate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.BasicShipment())
	require.NoError(err)

	estimates, err := client.GetShipmentEstimatedDeliveryDate(shipment.ID, c.fixture.PlannedShipDate())
	require.NoError(err)

	assert.True(len(estimates) > 0)
	for _, entry := range estimates {
		assert.NotNil(entry.EasyPostTimeInTransitData.EasyPostEstimatedDeliveryDate)
		assert.NotNil(entry.EasyPostTimeInTransitData.DaysInTransit)
		assert.NotNil(entry.EasyPostTimeInTransitData.PlannedShipDate)
	}
}

func (c *ClientTests) TestShipmentRecommendShipDate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.BasicShipment())
	require.NoError(err)

	recommendations, err := client.RecommendShipDateForShipment(shipment.ID, c.fixture.DesiredDeliveryDate())
	require.NoError(err)

	assert.True(len(recommendations) > 0)
	for _, entry := range recommendations {
		assert.NotNil(entry.EasyPostTimeInTransitData.EasyPostRecommendedShipDate)
		assert.NotNil(entry.EasyPostTimeInTransitData.DeliveryDateConfidence)
		assert.NotNil(entry.EasyPostTimeInTransitData.EstimatedTransitDays)
		assert.NotNil(entry.EasyPostTimeInTransitData.DaysInTransit)
		assert.NotNil(entry.EasyPostTimeInTransitData.DesiredDeliveryDate)
	}
}
