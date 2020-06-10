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

	// Create and verify an address
	toAddress, err := client.CreateAddress(
		&easypost.Address{
			Name:    "Bugs Bunny",
			Street1: "4000 Warner Blvd",
			City:    "Burbank",
			State:   "CA",
			Zip:     "91522",
		},
		&easypost.CreateAddressOptions{Verify: []string{"delivery"}},
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating address:", err)
		os.Exit(1)
		return
	}

	prettyJSON, err := json.MarshalIndent(toAddress, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating JSON:", err)
		os.Exit(1)
		return
	}
	fmt.Println(string(prettyJSON))
}
