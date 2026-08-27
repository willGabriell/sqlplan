package render

import (
	"bytes"
	"testing"

	"github.com/fatih/color"

	"sqlplan/explain"
)

func TestTree(t *testing.T) {
	color.NoColor = true

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
		"[1] Hash Join (cost=200.00..400.00 rows=500) (actual rows=480 time=5.120ms, 100.0% do total)\n" +
		"├── [2] Seq Scan on orders (cost=0.00..180.00 rows=8000) (actual rows=8000 time=1.900ms, 37.1% do total) ⚠ seq scan retornando muitas linhas\n" +
		"└── [3] Hash (cost=100.00..100.00 rows=500) (actual rows=500 time=0.800ms, 15.6% do total)\n"

	explain.AssignIDs(&n, new(int))

	var buf bytes.Buffer
	Tree(&buf, &n)

	if buf.String() != want {
		t.Fatalf("got:\n%s\nwant:\n%s", buf.String(), want)
	}
}

func TestTreeIndexScan(t *testing.T) {
	color.NoColor = true

	n := explain.PlanNode{
		NodeType:     "Index Scan",
		IndexName:    "users_pkey",
		RelationName: "users",
	}

	explain.AssignIDs(&n, new(int))

	var buf bytes.Buffer
	Tree(&buf, &n)

	want := "[1] Index Scan using users_pkey on users (cost=0.00..0.00 rows=0) (actual rows=0 time=0.000ms, 0.0% do total)\n"
	if buf.String() != want {
		t.Fatalf("got:\n%s\nwant:\n%s", buf.String(), want)
	}
}

// TestTreeReturnsDistinctNodeTypes garante que o retorno de Tree traz cada
// NodeType uma única vez, na ordem de aparição — usado pelo --glossary.
func TestTreeReturnsDistinctNodeTypes(t *testing.T) {
	color.NoColor = true

	n := explain.PlanNode{
		NodeType: "Nested Loop",
		Plans: []explain.PlanNode{
			{NodeType: "Seq Scan", RelationName: "a"},
			{NodeType: "Seq Scan", RelationName: "b"},
		},
	}

	explain.AssignIDs(&n, new(int))

	var buf bytes.Buffer
	types := Tree(&buf, &n)

	want := []string{"Nested Loop", "Seq Scan"}
	if len(types) != len(want) {
		t.Fatalf("got %v, want %v", types, want)
	}
	for i := range want {
		if types[i] != want[i] {
			t.Fatalf("got %v, want %v", types, want)
		}
	}
}

// TestTreeBranches trava o desenho de galhos (`├──`/`└──`/`│`) em 3+ níveis
// com múltiplos irmãos no mesmo nível.
func TestTreeBranches(t *testing.T) {
	color.NoColor = true

	n := explain.PlanNode{
		NodeType:        "Hash Join",
		ActualTotalTime: 10,
		Plans: []explain.PlanNode{
			{
				NodeType:        "Hash Join",
				ActualTotalTime: 5,
				Plans: []explain.PlanNode{
					{NodeType: "Seq Scan", RelationName: "a", ActualTotalTime: 2},
					{NodeType: "Seq Scan", RelationName: "b", ActualTotalTime: 3},
				},
			},
			{NodeType: "Hash", ActualTotalTime: 4},
		},
	}

	want := "" +
		"[1] Hash Join (cost=0.00..0.00 rows=0) (actual rows=0 time=10.000ms, 100.0% do total)\n" +
		"├── [2] Hash Join (cost=0.00..0.00 rows=0) (actual rows=0 time=5.000ms, 50.0% do total)\n" +
		"│   ├── [3] Seq Scan on a (cost=0.00..0.00 rows=0) (actual rows=0 time=2.000ms, 20.0% do total)\n" +
		"│   └── [4] Seq Scan on b (cost=0.00..0.00 rows=0) (actual rows=0 time=3.000ms, 30.0% do total)\n" +
		"└── [5] Hash (cost=0.00..0.00 rows=0) (actual rows=0 time=4.000ms, 40.0% do total)\n"

	explain.AssignIDs(&n, new(int))

	var buf bytes.Buffer
	Tree(&buf, &n)

	if buf.String() != want {
		t.Fatalf("got:\n%s\nwant:\n%s", buf.String(), want)
	}
}
