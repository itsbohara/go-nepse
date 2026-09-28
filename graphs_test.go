package nepse

import (
	"encoding/json"
	"os"
	"testing"
)

// NEPSE's scrip graph sends {"time","contractRate","contractQuantity"} rows (recorded 2026-09-28,
// NABIL); decoding must not drop the price.
func TestGraphDataPointScripGraphContractRate(t *testing.T) {
	raw, err := os.ReadFile("testdata/scrip_graph_nabil.json")
	if err != nil {
		t.Fatal(err)
	}
	var pts []GraphDataPoint
	if err := json.Unmarshal(raw, &pts); err != nil {
		t.Fatal(err)
	}
	if len(pts) != 171 {
		t.Fatalf("got %d points, want 171", len(pts))
	}
	for i, p := range pts {
		if p.Value == 0 {
			t.Fatalf("point %d has value 0: %+v", i, p)
		}
	}
	if first, last := pts[0], pts[len(pts)-1]; first.Timestamp != 1790571660 || first.Value != 569 || last.Value != 567 {
		t.Fatalf("first %+v last %+v, want 569 at 1790571660 … 567", first, last)
	}
}

func TestGraphDataPointFormats(t *testing.T) {
	cases := []struct {
		name, in string
		ts       int64
		v        float64
	}{
		{"index array", `[1790571660, 2613.33]`, 1790571660, 2613.33},
		{"legacy value object", `{"time": 1790571660, "value": 12.5}`, 1790571660, 12.5},
		{"contractRate wins", `{"time": 1, "contractRate": 569, "value": 1}`, 1, 569},
		{"null contractRate falls back to value", `{"time": 1, "contractRate": null, "value": 7}`, 1, 7},
	}
	for _, c := range cases {
		var p GraphDataPoint
		if err := json.Unmarshal([]byte(c.in), &p); err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if p.Timestamp != c.ts || p.Value != c.v {
			t.Errorf("%s: got %+v, want {%d %v}", c.name, p, c.ts, c.v)
		}
	}
}
