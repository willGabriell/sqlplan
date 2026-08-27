package render

import (
	"bytes"
	"testing"

	"github.com/fatih/color"

	"sqlplan/explain"
)

func TestSummary(t *testing.T) {
	color.NoColor = true

	bs := []explain.Bottleneck{
		{Node: &explain.PlanNode{ID: 1, NodeType: "Hash Join"}, SelfTime: 1.203, Percent: 50.1},
		{Node: &explain.PlanNode{ID: 2, NodeType: "Seq Scan", RelationName: "orders"}, SelfTime: 0.891, Percent: 37.1},
	}

	want := "" +
		"Top 2 gargalos (de 2.400ms totais):\n" +
		"  1. [1] Hash Join                       1.203ms  (50.1%)\n" +
		"  2. [2] Seq Scan on orders              0.891ms  (37.1%)\n" +
		"\n"

	var buf bytes.Buffer
	Summary(&buf, bs, 2.4)

	if buf.String() != want {
		t.Fatalf("got:\n%s\nwant:\n%s", buf.String(), want)
	}
}

func TestSummaryIndexScanLabel(t *testing.T) {
	color.NoColor = true

	bs := []explain.Bottleneck{
		{Node: &explain.PlanNode{ID: 1, NodeType: "Index Scan", IndexName: "users_pkey", RelationName: "users"}, SelfTime: 0.5, Percent: 100},
	}

	var buf bytes.Buffer
	Summary(&buf, bs, 0.5)

	want := "" +
		"Top 1 gargalos (de 0.500ms totais):\n" +
		"  1. [1] Index Scan using users_pkey on users    0.500ms  (100.0%)\n" +
		"\n"

	if buf.String() != want {
		t.Fatalf("got:\n%s\nwant:\n%s", buf.String(), want)
	}
}

func TestSummaryEmpty(t *testing.T) {
	var buf bytes.Buffer
	Summary(&buf, nil, 10)

	if buf.String() != "" {
		t.Fatalf("esperava output vazio pra lista de bottlenecks vazia, veio %q", buf.String())
	}
}
