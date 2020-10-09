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

    // Create a ScanForm for a Batch
    batch, err := client.CreateBatchScanForms(
        "batch_12345",
        "string",
    )
    if err != nil {
        fmt.Fprintln(os.Stderr, "error creating ScanForm for Batch:", err)
        os.Exit(1)
        return
    }

    prettyJSON, err := json.MarshalIndent(batch, "", "    ")
    if err != nil {
        fmt.Fprintln(os.Stderr, "error creating JSON:", err)
        os.Exit(1)
        return
    }
    fmt.Println(string(prettyJSON))
}
