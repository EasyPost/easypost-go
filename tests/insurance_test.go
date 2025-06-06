package easypost_test

import (
	"errors"
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v5"
)

func (c *ClientTests) TestInsuranceCreate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	require.NoError(err)

	insuranceData := c.fixture.BasicInsurance()
	insuranceData.TrackingCode = shipment.TrackingCode

	insurance, err := client.CreateInsurance(insuranceData)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Insurance{}), reflect.TypeOf(insurance))
	assert.True(strings.HasPrefix(insurance.ID, "ins_"))
	assert.Equal("100.00000", insurance.Amount)
}

func (c *ClientTests) TestInsuranceRetrieve() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	insurances, err := client.ListInsurances(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	retrievedInsurance, err := client.GetInsurance(insurances.Insurances[0].ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Insurance{}), reflect.TypeOf(retrievedInsurance))
	assert.Equal(insurances.Insurances[0], retrievedInsurance)
}

func (c *ClientTests) TestInsuranceAll() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	insurances, err := client.ListInsurances(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	insurancesList := insurances.Insurances

	assert.LessOrEqual(len(insurancesList), c.fixture.pageSize())
	assert.NotNil(insurances.HasMore)
	for _, insurance := range insurancesList {
		assert.Equal(reflect.TypeOf(&easypost.Insurance{}), reflect.TypeOf(insurance))
	}
}

func (c *ClientTests) TestInsuranceGetNextPage() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	firstPage, err := client.ListInsurances(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	nextPage, err := client.GetNextInsurancePageWithPageSize(firstPage, c.fixture.pageSize())
	defer func() {
		if err == nil {
			assert.True(len(nextPage.Insurances) <= c.fixture.pageSize())

			lastIDOfFirstPage := firstPage.Insurances[len(firstPage.Insurances)-1].ID
			firstIdOfSecondPage := nextPage.Insurances[0].ID

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

func (c *ClientTests) TestInsuranceRefund() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	insuranceData := c.fixture.BasicInsurance()
	insuranceData.TrackingCode = "EZ1000000001"

	insurance, _ := client.CreateInsurance(insuranceData)
	refundInsurance, err := client.RefundInsurance(insurance.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Insurance{}), reflect.TypeOf(refundInsurance))
	assert.True(strings.HasPrefix(refundInsurance.ID, "ins_"))
	assert.Equal("cancelled", refundInsurance.Status)
	assert.Equal("Insurance was cancelled by the user.", refundInsurance.Messages[0])
}
