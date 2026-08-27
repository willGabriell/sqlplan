package render

import (
	"bytes"
	"testing"

	"sqlplan/explain"
)

func TestTree(t *testing.T) {
	n := explain.PlanNode{
		NodeType:        "Hash Join",
		StartupCost:     200,
		TotalCost:       400,
		PlanRows:        500,
		ActualRows:      480,
		ActualTotalTime: 5.12,
		Plans: []explain.PlanNode{
			{
				NodeType:        "Seq Scan",
				RelationName:    "orders",
				StartupCost:     0,
				TotalCost:       180,
				PlanRows:        8000,
				ActualRows:      8000,
				ActualTotalTime: 1.9,
			},
			{
				NodeType:        "Hash",
				StartupCost:     100,
				TotalCost:       100,
				PlanRows:        500,
				ActualRows:      500,
				ActualTotalTime: 0.8,
			},
		},
	}

	want := "" +
		"Hash Join (cost=200.00..400.00 rows=500) (actual rows=480 time=5.120ms)\n" +
		"  Seq Scan on orders (cost=0.00..180.00 rows=8000) (actual rows=8000 time=1.900ms)\n" +
		"  Hash (cost=100.00..100.00 rows=500) (actual rows=500 time=0.800ms)\n"

	var buf bytes.Buffer
	Tree(&buf, &n, 0)

	if buf.String() != want {
		t.Fatalf("got:\n%s\nwant:\n%s", buf.String(), want)
	}
}

func TestTreeIndexScan(t *testing.T) {
	n := explain.PlanNode{
		NodeType:     "Index Scan",
		IndexName:    "users_pkey",
		RelationName: "users",
	}

	var buf bytes.Buffer
	Tree(&buf, &n, 0)

	want := "Index Scan using users_pkey on users (cost=0.00..0.00 rows=0) (actual rows=0 time=0.000ms)\n"
	if buf.String() != want {
		t.Fatalf("got:\n%s\nwant:\n%s", buf.String(), want)
	}
}
