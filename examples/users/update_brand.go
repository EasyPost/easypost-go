package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/EasyPost/easypost-go"
)

func main() {
	apiKey := os.Getenv("EASYPOST_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "missing API key")
		os.Exit(1)
		return
	}
	client := easypost.New(apiKey)

	// Get current user
	user, err := client.RetrieveMe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error retrieving current user:", err)
		os.Exit(1)
		return
	}

	params := map[string]interface{}{"color": "#aaaaaa"}

	// Update the user brand
	brand, err := client.UpdateBrand(params, user.ID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error to update brand:", err)
		os.Exit(1)
		return
	}

	prettyJSON, err := json.MarshalIndent(brand, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating JSON:", err)
	}
	fmt.Printf("%s\n", string(prettyJSON))
}
