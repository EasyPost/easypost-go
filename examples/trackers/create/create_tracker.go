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

	// Create a tracker
	tracker, err := client.CreateTracker(
		&easypost.CreateTrackerOptions{
			TrackingCode: "EZ1000000001",
			Carrier:      "USPS",
		},
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating tracker:", err)
		os.Exit(1)
		return
	}

	prettyJSON, err := json.MarshalIndent(tracker, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating JSON:", err)
	}
	fmt.Printf("%s\n", string(prettyJSON))
}
