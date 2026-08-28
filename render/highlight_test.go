package render

import (
	"testing"

	"github.com/willGabriell/sqlplan/explain"
)

func TestEvaluateNode(t *testing.T) {
	tests := []struct {
		name string
		node explain.PlanNode
		want []string
	}{
		{
			name: "seq scan abaixo do limiar não destaca",
			node: explain.PlanNode{NodeType: "Seq Scan", ActualRows: 999, PlanRows: 999},
			want: nil,
		},
		{
			name: "seq scan acima do limiar destaca",
			node: explain.PlanNode{NodeType: "Seq Scan", ActualRows: 1001, PlanRows: 1001},
			want: []string{"seq scan retornando muitas linhas"},
		},
		{
			name: "index scan com muitas linhas não é seq scan, não destaca por essa regra",
			node: explain.PlanNode{NodeType: "Index Scan", ActualRows: 50000, PlanRows: 50000},
			want: nil,
		},
		{
			name: "estimativa 10x pra mais destaca",
			node: explain.PlanNode{NodeType: "Hash Join", ActualRows: 1000, PlanRows: 10},
			want: []string{"estimativa muito distante da real"},
		},
		{
			name: "estimativa 10x pra menos destaca",
			node: explain.PlanNode{NodeType: "Hash Join", ActualRows: 10, PlanRows: 1000},
			want: []string{"estimativa muito distante da real"},
		},
		{
			name: "estimativa dentro da razão não destaca",
			node: explain.PlanNode{NodeType: "Hash Join", ActualRows: 100, PlanRows: 500},
			want: nil,
		},
		{
			name: "actual rows zero não explode por divisão por zero",
			node: explain.PlanNode{NodeType: "Hash Join", ActualRows: 0, PlanRows: 10000},
			want: []string{"estimativa muito distante da real"},
		},
		{
			name: "plan rows zero não explode por divisão por zero",
			node: explain.PlanNode{NodeType: "Seq Scan", ActualRows: 10000, PlanRows: 0},
			want: []string{"seq scan retornando muitas linhas", "estimativa muito distante da real"},
		},
		{
			name: "ambos zero não destaca",
			node: explain.PlanNode{NodeType: "Hash Join", ActualRows: 0, PlanRows: 0},
			want: nil,
		},
		{
			name: "seq scan gordo e estimativa errada acumulam os dois",
			node: explain.PlanNode{NodeType: "Seq Scan", ActualRows: 9987, PlanRows: 10},
			want: []string{"seq scan retornando muitas linhas", "estimativa muito distante da real"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hs := EvaluateNode(&tt.node)
			if len(hs) != len(tt.want) {
				t.Fatalf("got %d highlights, want %d: %+v", len(hs), len(tt.want), hs)
			}
			for i, h := range hs {
				if h.Reason != tt.want[i] {
					t.Errorf("highlight %d: got %q, want %q", i, h.Reason, tt.want[i])
				}
			}
		})
	}
}
