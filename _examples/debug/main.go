package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/voidarchive/go-nepse"
)

func main() {
	opts := nepse.DefaultOptions()
	opts.TLSVerification = false

	client, err := nepse.NewClient(opts)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer func() { _ = client.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data, err := client.Companies(ctx)
	if err != nil {
		panic(err)
	}

	// Count equity by status
	statusCounts := make(map[string]int)
	for _, c := range data {
		if c.InstrumentType == "Equity" {
			statusCounts[c.Status]++
		}
	}

	fmt.Printf("Total companies: %d\n", len(data))
	fmt.Println("Equity by status:")
	for s, count := range statusCounts {
		fmt.Printf("  Status '%s': %d\n", s, count)
	}
}
