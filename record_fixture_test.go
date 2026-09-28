package nepse

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

// TestRecordScripGraphFixture refreshes testdata/scrip_graph_nabil.json from the live API.
// Opt-in: NEPSE_RECORD=1 go test -run TestRecordScripGraphFixture .
func TestRecordScripGraphFixture(t *testing.T) {
	if os.Getenv("NEPSE_RECORD") != "1" {
		t.Skip("set NEPSE_RECORD=1 to call the live API")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	c, err := NewClient(DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	sec, err := c.findSecurityBySymbol(ctx, "NABIL")
	if err != nil {
		t.Fatal(err)
	}
	id, err := c.computeScripGraphPayloadID(ctx)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := c.apiPostRequestRaw(ctx, fmt.Sprintf("%s/%d", c.config.Endpoints.CompanyDailyGraph, sec.ID), graphPostPayload{ID: id})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile("testdata/scrip_graph_nabil.json", raw, 0o644); err != nil {
		t.Fatal(err)
	}
	t.Logf("wrote %d bytes", len(raw))
}
