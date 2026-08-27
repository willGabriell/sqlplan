package explain

import "testing"

// TestFindBottlenecksLoops trava o critério de aceite principal: um nó com
// Actual Loops alto (ex: lado interno de nested loop) precisa rankear pelo
// tempo total real (ActualTotalTime*Loops), não pela média por execução.
func TestFindBottlenecksLoops(t *testing.T) {
	root := PlanNode{
		NodeType:        "Nested Loop",
		ActualTotalTime: 500.5,
		ActualLoops:     1,
		Plans: []PlanNode{
			{NodeType: "Seq Scan", RelationName: "outer", ActualTotalTime: 200, ActualLoops: 1},
			{NodeType: "Index Scan", RelationName: "inner", ActualTotalTime: 0.5, ActualLoops: 1000},
		},
	}

	bs := FindBottlenecks(&root, 700.5, 3)
	if len(bs) != 3 {
		t.Fatalf("esperava 3 bottlenecks, veio %d", len(bs))
	}
	if bs[0].Node.RelationName != "inner" {
		t.Fatalf("esperava 'inner' (500ms reais em loop) no topo, veio %q", bs[0].Node.RelationName)
	}
	if bs[0].SelfTime != 500 {
		t.Fatalf("self time do nó em loop errado: got %.3f, want 500", bs[0].SelfTime)
	}
}

// TestFindBottlenecksSumMatchesExecutionTime é o teste de sanidade mais
// importante da sprint: soma de self time de todos os nós ≈ tempo total.
func TestFindBottlenecksSumMatchesExecutionTime(t *testing.T) {
	root := PlanNode{
		NodeType:        "Hash Join",
		ActualTotalTime: 10,
		ActualLoops:     1,
		Plans: []PlanNode{
			{
				NodeType:        "Seq Scan",
				RelationName:    "a",
				ActualTotalTime: 4,
				ActualLoops:     1,
			},
			{
				NodeType:        "Hash",
				ActualTotalTime: 3,
				ActualLoops:     1,
				Plans: []PlanNode{
					{NodeType: "Seq Scan", RelationName: "b", ActualTotalTime: 2, ActualLoops: 1},
				},
			},
		},
	}

	bs := FindBottlenecks(&root, 10, 999)

	var sum float64
	for _, b := range bs {
		sum += b.SelfTime
	}

	const tolerance = 0.01
	if diff := sum - root.ActualTotalTime; diff > tolerance || diff < -tolerance {
		t.Fatalf("soma dos self times = %.4f, esperava ~%.4f (tolerância %.2f)", sum, root.ActualTotalTime, tolerance)
	}
}

// TestFindBottlenecksNoNegative garante que arredondamento de float nunca
// produz self time negativo, mesmo quando filhos somam levemente mais que o pai.
func TestFindBottlenecksNoNegative(t *testing.T) {
	root := PlanNode{
		NodeType:        "Hash Join",
		ActualTotalTime: 10.0,
		ActualLoops:     1,
		Plans: []PlanNode{
			{NodeType: "Seq Scan", RelationName: "a", ActualTotalTime: 6.0000001, ActualLoops: 1},
			{NodeType: "Seq Scan", RelationName: "b", ActualTotalTime: 4.0000001, ActualLoops: 1},
		},
	}

	bs := FindBottlenecks(&root, 10, 999)
	for _, b := range bs {
		if b.SelfTime < 0 {
			t.Fatalf("self time negativo pro nó %s: %.6f", b.Node.NodeType, b.SelfTime)
		}
	}
}

// TestFindBottlenecksTopNBiggerThanTree garante que topN maior que a
// árvore devolve tudo, sem panic de slice out of range.
func TestFindBottlenecksTopNBiggerThanTree(t *testing.T) {
	root := PlanNode{
		NodeType:        "Seq Scan",
		RelationName:    "solo",
		ActualTotalTime: 1,
		ActualLoops:     1,
	}

	bs := FindBottlenecks(&root, 1, 999)
	if len(bs) != 1 {
		t.Fatalf("esperava 1 bottleneck (árvore de 1 nó), veio %d", len(bs))
	}
}
