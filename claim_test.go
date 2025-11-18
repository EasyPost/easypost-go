package easypost

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

func prepareInsuredShipment(client *Client, shipmentCreateParams *Shipment, claimAmount float64) *Shipment {
	shipment, err := client.CreateShipment(shipmentCreateParams)
	if err != nil {
		panic(err)
	}
	rate, err := client.LowestShipmentRate(shipment)
	if err != nil {
		panic(err)
	}
	claimAmountStr := strconv.FormatFloat(claimAmount, 'f', 2, 64)
	purchasedShipment, err := client.BuyShipment(shipment.ID, &rate, claimAmountStr)
	if err != nil {
		panic(err)
	}

	return purchasedShipment
}

func (c *ClientTests) TestClaimCreate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	claimAmount := 100.0
	claimAmountStr := strconv.FormatFloat(claimAmount, 'f', 2, 64)

	createShipmentParams := c.fixture.FullShipment()
	insuredShipment := prepareInsuredShipment(client, createShipmentParams, claimAmount)

	createClaimParams := c.fixture.BasicClaim()
	createClaimParams.TrackingCode = insuredShipment.TrackingCode
	createClaimParams.Amount = claimAmount

	claim, err := client.CreateClaim(createClaimParams)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&Claim{}), reflect.TypeOf(claim))
	assert.True(strings.HasPrefix(claim.ID, "clm_"))
	assert.Equal(createClaimParams.TrackingCode, claim.TrackingCode)
	assert.Equal(claimAmountStr, claim.RequestedAmount)
}

func (c *ClientTests) TestClaimRetrieve() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	claimAmount := 100.0

	createShipmentParams := c.fixture.FullShipment()
	insuredShipment := prepareInsuredShipment(client, createShipmentParams, claimAmount)

	createClaimParams := c.fixture.BasicClaim()
	createClaimParams.TrackingCode = insuredShipment.TrackingCode
	createClaimParams.Amount = claimAmount

	claim, err := client.CreateClaim(createClaimParams)
	require.NoError(err)

	retrievedClaim, err := client.GetClaim(claim.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&Claim{}), reflect.TypeOf(retrievedClaim))
	assert.Equal(claim.ID, retrievedClaim.ID)
}

func (c *ClientTests) TestClaimAll() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	listClaimsParams := &ListClaimsParameters{
		PageSize: c.fixture.pageSize(),
	}

	claims, err := client.ListClaims(listClaimsParams)
	require.NoError(err)

	claimsList := claims.Claims

	assert.LessOrEqual(len(claimsList), c.fixture.pageSize())
	assert.NotNil(claims.HasMore)
	for _, claim := range claimsList {
		assert.Equal(reflect.TypeOf(&Claim{}), reflect.TypeOf(claim))
	}
}

func (c *ClientTests) TestClaimGetNextPage() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	listClaimsParams := &ListClaimsParameters{
		PageSize: c.fixture.pageSize(),
	}

	firstPage, err := client.ListClaims(listClaimsParams)
	require.NoError(err)

	nextPage, err := client.GetNextClaimPageWithPageSize(firstPage, c.fixture.pageSize())
	defer func() {
		if err == nil {
			assert.True(len(nextPage.Claims) <= c.fixture.pageSize())

			lastIDOfFirstPage := firstPage.Claims[len(firstPage.Claims)-1].ID
			firstIdOfSecondPage := nextPage.Claims[0].ID

			assert.NotEqual(lastIDOfFirstPage, firstIdOfSecondPage)
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

func (c *ClientTests) TestClaimCancel() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	claimAmount := 100.0

	createShipmentParams := c.fixture.FullShipment()
	insuredShipment := prepareInsuredShipment(client, createShipmentParams, claimAmount)

	createClaimParams := c.fixture.BasicClaim()
	createClaimParams.TrackingCode = insuredShipment.TrackingCode
	createClaimParams.Amount = claimAmount

	claim, err := client.CreateClaim(createClaimParams)
	require.NoError(err)

	cancelledClaim, err := client.CancelClaim(claim.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&Claim{}), reflect.TypeOf(cancelledClaim))
	assert.True(strings.HasPrefix(cancelledClaim.ID, "clm_"))
	assert.Equal("cancelled", cancelledClaim.Status)
}
