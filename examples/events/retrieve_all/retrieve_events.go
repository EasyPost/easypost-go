package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/EasyPost/easypost-go/v2"
)

func main() {
	apiKey := os.Getenv("EASYPOST_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "missing API key")
		os.Exit(1)
		return
	}
	client := easypost.New(apiKey)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	// Retrieve events from the past day.
	yesterday := time.Now().Add(-24 * time.Hour)
	opts := &easypost.ListOptions{StartDateTime: &yesterday}

	results, err := &easypost.ListEventsResult{HasMore: true}, error(nil)
	for results.HasMore && err == nil {
		if results, err = client.ListEvents(opts); err == nil {
			for i := range results.Events {
				_ = enc.Encode(results.Events[i])
				// If a webhook is registered, payloads can be examined to
				// obtain the event result.
				payloads, _ := client.ListEventPayloads(results.Events[i].ID)
				if len(payloads) != 0 {
					event, _ := payloads[0].RequestBody.(*easypost.Event)
					if event != nil {
						_ = enc.Encode(event.Result)
					}
				}
			}
			if results.HasMore {
				// Update BeforeID in order to fetch additional pages.
				opts.BeforeID = results.Events[len(results.Events)-1].ID
			}
		}
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "error retrieving events:", err)
		os.Exit(1)
	}
}
