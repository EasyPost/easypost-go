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

	// Retrieve a rate
	rate, err := client.GetRate("rate_123")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error retrieving rate:", err)
		os.Exit(1)
		return
	}

	prettyJSON, err := json.MarshalIndent(rate, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating JSON:", err)
		os.Exit(1)
		return
	}
	fmt.Println(string(prettyJSON))
}
