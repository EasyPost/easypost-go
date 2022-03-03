package easypost_test

import (
	"fmt"
	"github.com/EasyPost/easypost-go/v2"
	"reflect"
	"strings"
)

func (c *ClientTests) TestBatchCreate() {
	client := c.TestClient()
	assert := c.Assert()

	batch, _ := client.CreateBatch(
		c.fixture.BasicShipment(),
		nil,
	)

	assert.Equal(reflect.TypeOf(&easypost.Batch{}), reflect.TypeOf(batch))
	assert.True(strings.HasPrefix(batch.ID, "batch_"))
	assert.NotNil(batch.Shipments)
}

func (c *ClientTests) TestBatchRetrieve() {
	client := c.TestClient()
	assert := c.Assert()

	batch, _ := client.CreateBatch(
		c.fixture.BasicShipment(),
		nil,
	)

	retrievedBatch, _ := client.GetBatch(batch.ID)

	assert.Equal(reflect.TypeOf(&easypost.Batch{}), reflect.TypeOf(retrievedBatch))
	assert.Equal(batch.ID, retrievedBatch.ID)
}

// This test is currently skipped because we need to fix the typo of ListBatchesResult from `insurance` to `batch` in next PR
func (c *ClientTests) TestBatchAll() {
	c.T().Skip("This test is currently skipped because we need to fix the typo of ListBatchesResult from `insurance` to `batch` in a separate PR")

	client := c.TestClient()
	assert := c.Assert()

	batches, _ := client.ListBatches(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)

	for _, batch := range batches.Insurances {
		assert.Equal(reflect.TypeOf(&easypost.Batch{}), reflect.TypeOf(batch))
	}
}

func (c *ClientTests) TestBatchCreateAndBuy() {
	client := c.TestClient()
	assert := c.Assert()

	batch, _ := client.CreateAndBuyBatch(c.fixture.OneCallBuyShipment(), c.fixture.OneCallBuyShipment())

	assert.Equal(reflect.TypeOf(&easypost.Batch{}), reflect.TypeOf(batch))
	assert.True(strings.HasPrefix(batch.ID, "batch_"))
	assert.Equal(2, batch.NumShipments)
}

func (c *ClientTests) TestBatchBuy() {
	client := c.TestClient()
	assert := c.Assert()

	batch, _ := client.CreateBatch(
		c.fixture.OneCallBuyShipment(),
		nil,
	)

	boughtBatch, _ := client.BuyBatch(batch.ID)

	assert.Equal(reflect.TypeOf(&easypost.Batch{}), reflect.TypeOf(boughtBatch))
	assert.Equal(1, boughtBatch.NumShipments)
}

func (c *ClientTests) TestBatchCreateScanForm() {
	client := c.TestClient()
	assert := c.Assert()

	batch, _ := client.CreateBatch(
		c.fixture.OneCallBuyShipment(),
		nil,
	)

	boughtBatch, _ := client.BuyBatch(batch.ID)

	batchWithScanForm, _ := client.CreateBatchScanForms(boughtBatch.ID, "")

	// We can't assert anything meaningful here because the scanform gets queued for generation and may not be immediately available
	assert.Equal(reflect.TypeOf(&easypost.Batch{}), reflect.TypeOf(batchWithScanForm))
}

func (c *ClientTests) TestBatchAddRemoveShipment() {
	client := c.TestClient()
	assert := c.Assert()

	shipment, _ := client.CreateShipment(c.fixture.OneCallBuyShipment())
	shipment2, _ := client.CreateShipment(c.fixture.OneCallBuyShipment())

	batch, _ := client.CreateBatch()

	batchWithShipment, err := client.AddShipmentsToBatch(batch.ID, shipment.ID, shipment2.ID)

	fmt.Println(err)
	assert.Equal(2, batchWithShipment.NumShipments)

	batchWithoutShipment, err := client.RemoveShipmentsFromBatch(batch.ID, shipment.ID, shipment2.ID)

	fmt.Println(err)
	assert.Equal(0, batchWithoutShipment.NumShipments)
}

func (c *ClientTests) TestBatchLabel() {
	client := c.TestClient()
	assert := c.Assert()

	batch, _ := client.CreateBatch(
		c.fixture.OneCallBuyShipment(),
		nil,
	)

	boughtBatch, _ := client.BuyBatch(batch.ID)

	batchWithLabel, _ := client.GetBatchLabels(boughtBatch.ID, "ZPL")

	// We can't assert anything meaningful here because the label gets queued for generation and may not be immediately available
	assert.Equal(reflect.TypeOf(&easypost.Batch{}), reflect.TypeOf(batchWithLabel))
}
