package main

import (
	"encoding/json"
	"fmt"
	"github.com/EasyPost/easypost-go"
	"os"
)

func main() {
	apiKey := os.Getenv("EASYPOST_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "missing API key")
		os.Exit(1)
		return
	}
	client := easypost.New(apiKey)

	// Create a batch
	batch, err := client.CreateBatch(
		&easypost.Shipment{ID: "shp_100"},
		&easypost.Shipment{ID: "shp_101"},
		&easypost.Shipment{ID: "shp_102"},
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating batch:", err)
		os.Exit(1)
		return
	}

	prettyJSON, err := json.MarshalIndent(batch, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating JSON:", err)
	}
	fmt.Printf("%s\n", string(prettyJSON))
}
