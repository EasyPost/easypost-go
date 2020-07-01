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

	// Retrieve a tracker
	tracker, err := client.GetTracker("EZ1000000001")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error retrieving tracker:", err)
		os.Exit(1)
		return
	}

	prettyJSON, err := json.MarshalIndent(tracker, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating JSON:", err)
	}
	fmt.Printf("%s\n", string(prettyJSON))
}
