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

	// Retrieve a parcel
	parcel, err := client.GetParcel("prcl_123")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error retrieving parcel:", err)
		os.Exit(1)
		return
	}

	prettyJSON, err := json.MarshalIndent(parcel, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating JSON:", err)
		os.Exit(1)
		return
	}
	fmt.Println(string(prettyJSON))
}
