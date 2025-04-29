package easypost_test

import (
	"errors"

	"github.com/EasyPost/easypost-go/v5"
)

func (c *ClientTests) TestGetNextPage() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	addresses, err := client.ListAddresses(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	nextPage, err := client.GetNextAddressPage(addresses)
	defer func() {
		if err == nil {
			lastIDOfFirstPage := addresses.Addresses[len(addresses.Addresses)-1].ID
			firstIdOfSecondPage := nextPage.Addresses[0].ID

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

func (c *ClientTests) TestGetNextPageWithPageSize() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	addresses, err := client.ListAddresses(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	nextPage, err := client.GetNextAddressPageWithPageSize(addresses, c.fixture.pageSize())
	defer func() {
		if err == nil {
			assert.True(len(nextPage.Addresses) <= c.fixture.pageSize())

			lastIDOfFirstPage := addresses.Addresses[len(addresses.Addresses)-1].ID
			firstIdOfSecondPage := nextPage.Addresses[0].ID

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

func (c *ClientTests) TestGetNextPageReachEnd() {
	mockRequests := []easypost.MockRequest{
		{
			MatchRule: easypost.MockRequestMatchRule{
				Method:          "GET",
				UrlRegexPattern: "v2\\/addresses",
			},
			ResponseInfo: easypost.MockRequestResponseInfo{
				StatusCode: 200,
				Body:       `{"addresses": [], "has_more": false}`,
			},
		},
	}

	client := c.MockClient(mockRequests)
	assert, require := c.Assert(), c.Require()

	addresses, err := client.ListAddresses(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	hitEnd := false

	for {
		addresses, err = client.GetNextAddressPage(addresses)
		if err != nil {
			var endOfPaginationErr *easypost.EndOfPaginationError
			if errors.As(err, &endOfPaginationErr) {
				hitEnd = true
			} else {
				require.NoError(err)
			}
		}
		if hitEnd {
			break
		}
	}

	assert.True(hitEnd)
}
