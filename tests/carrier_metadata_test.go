package easypost_test

func (c *ClientTests) TestBetaGetCarrierMetadata() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	carrierMetadata, err := client.GetCarrierMetadata()
	require.NoError(err)

	// Assert we get multiple carriers
	uspsFound := false
	fedexFound := false
	for _, carrier := range carrierMetadata {
		if carrier.Name == "usps" {
			uspsFound = true
		}
		if carrier.Name == "fedex" {
			fedexFound = true
		}
		if uspsFound && fedexFound {
			break
		}
	}
	assert.True(uspsFound)
	assert.True(fedexFound)
}

func (c *ClientTests) TestBetaGetCarrierMetadataWithCarriersAndTypes() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	carrierMetadata, err := client.GetCarrierMetadataWithCarriersAndTypes([]string{"usps"}, []string{"service_levels", "predefined_packages"})
	require.NoError(err)

	// Assert we get the single carrier we asked for and only the types we asked for
	uspsFound := false
	for _, carrier := range carrierMetadata {
		if carrier.Name == "usps" {
			uspsFound = true
			break
		}
	}
	assert.True(uspsFound)
	assert.Equal(1, len(carrierMetadata))
	assert.NotNil(carrierMetadata[0].ServiceLevels)
	assert.NotNil(carrierMetadata[0].PredefinedPackages)
	assert.Nil(carrierMetadata[0].SupportedFeatures)
}
