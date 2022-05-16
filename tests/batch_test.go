package easypost_test

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestBatchCreate() {
	client := c.TestClient()
	assert := c.Assert()

	batch, err := client.CreateBatch(
		c.fixture.BasicShipment(),
		nil,
	)
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(reflect.TypeOf(&easypost.Batch{}), reflect.TypeOf(batch))
	assert.True(strings.HasPrefix(batch.ID, "batch_"))
	assert.NotNil(batch.Shipments)
}

func (c *ClientTests) TestBatchRetrieve() {
	client := c.TestClient()
	assert := c.Assert()

	batch, err := client.CreateBatch(
		c.fixture.BasicShipment(),
		nil,
	)
	if err != nil {
		c.T().Error(err)
		return
	}

	retrievedBatch, err := client.GetBatch(batch.ID)
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(reflect.TypeOf(&easypost.Batch{}), reflect.TypeOf(retrievedBatch))
	assert.Equal(batch.ID, retrievedBatch.ID)
}

func (c *ClientTests) TestBatchAll() {
	client := c.TestClient()
	assert := c.Assert()

	batches, err := client.ListBatches(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.LessOrEqual(len(batches.Batch), c.fixture.pageSize())
	assert.NotNil(batches.HasMore)
	for _, batch := range batches.Batch {
		assert.Equal(reflect.TypeOf(&easypost.Batch{}), reflect.TypeOf(batch))
	}
}

func (c *ClientTests) TestBatchCreateAndBuy() {
	client := c.TestClient()
	assert := c.Assert()

	batch, err := client.CreateAndBuyBatch(c.fixture.OneCallBuyShipment(), c.fixture.OneCallBuyShipment())
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(reflect.TypeOf(&easypost.Batch{}), reflect.TypeOf(batch))
	assert.True(strings.HasPrefix(batch.ID, "batch_"))
	assert.Equal(2, batch.NumShipments)
}

func (c *ClientTests) TestBatchBuy() {
	client := c.TestClient()
	assert := c.Assert()

	batch, err := client.CreateBatch(
		c.fixture.OneCallBuyShipment(),
		nil,
	)
	if err != nil {
		c.T().Error(err)
		return
	}

	boughtBatch, err := client.BuyBatch(batch.ID)
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(reflect.TypeOf(&easypost.Batch{}), reflect.TypeOf(boughtBatch))
	assert.Equal(1, boughtBatch.NumShipments)
}

func (c *ClientTests) TestBatchCreateScanForm() {
	client := c.TestClient()
	assert := c.Assert()

	batch, err := client.CreateBatch(
		c.fixture.OneCallBuyShipment(),
		nil,
	)
	if err != nil {
		c.T().Error(err)
		return
	}

	boughtBatch, err := client.BuyBatch(batch.ID)
	if err != nil {
		c.T().Error(err)
		return
	}

	currentDir, _ := os.Getwd()
	cassettePath := filepath.Join(currentDir, "cassettes", "TestBatchCreateScanForm.yaml")
	if _, err := os.Stat(cassettePath); errors.Is(err, os.ErrNotExist) {
		time.Sleep(5 * time.Second) // Wait enough time for the batch to process buying the shipment
	}

	batchWithScanForm, err := client.CreateBatchScanForms(boughtBatch.ID, "")
	if err != nil {
		c.T().Error(err)
		return
	}

	// We can't assert anything meaningful here because the scanform gets queued for generation and may not be immediately available
	assert.Equal(reflect.TypeOf(&easypost.Batch{}), reflect.TypeOf(batchWithScanForm))
}

func (c *ClientTests) TestBatchAddRemoveShipment() {
	client := c.TestClient()
	assert := c.Assert()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	if err != nil {
		c.T().Error(err)
		return
	}
	shipment2, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	if err != nil {
		c.T().Error(err)
		return
	}

	batch, err := client.CreateBatch()
	if err != nil {
		c.T().Error(err)
		return
	}

	batchWithShipment, err := client.AddShipmentsToBatch(batch.ID, shipment.ID, shipment2.ID)
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(2, batchWithShipment.NumShipments)

	batchWithoutShipment, err := client.RemoveShipmentsFromBatch(batch.ID, shipment.ID, shipment2.ID)
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(0, batchWithoutShipment.NumShipments)
}

func (c *ClientTests) TestBatchLabel() {
	client := c.TestClient()
	assert := c.Assert()

	batch, err := client.CreateBatch(
		c.fixture.OneCallBuyShipment(),
		nil,
	)
	if err != nil {
		c.T().Error(err)
		return
	}

	boughtBatch, err := client.BuyBatch(batch.ID)
	if err != nil {
		c.T().Error(err)
		return
	}

	currentDir, _ := os.Getwd()
	cassettePath := filepath.Join(currentDir, "cassettes", "TestBatchLabel.yaml")
	if _, err := os.Stat(cassettePath); errors.Is(err, os.ErrNotExist) {
		time.Sleep(5 * time.Second) // Wait enough time for the batch to process buying the shipment
	}

	batchWithLabel, err := client.GetBatchLabels(boughtBatch.ID, "ZPL")
	if err != nil {
		c.T().Error(err)
		return
	}

	// We can't assert anything meaningful here because the label gets queued for generation and may not be immediately available
	assert.Equal(reflect.TypeOf(&easypost.Batch{}), reflect.TypeOf(batchWithLabel))
}
