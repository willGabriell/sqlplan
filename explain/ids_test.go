package explain

import "testing"

// TestAssignIDsPreOrder trava o critério de aceite principal: nós repetidos
// do mesmo NodeType (várias "Nested Loop") recebem IDs diferentes e
// sequenciais, na ordem raiz -> filhos da esquerda pra direita — mesma ordem
// em que a árvore é impressa.
func TestAssignIDsPreOrder(t *testing.T) {
	root := PlanNode{
		NodeType: "Nested Loop",
		Plans: []PlanNode{
			{
				NodeType: "Nested Loop",
				Plans: []PlanNode{
					{NodeType: "Seq Scan", RelationName: "a"},
					{NodeType: "Seq Scan", RelationName: "b"},
				},
			},
			{
				NodeType: "Nested Loop",
				Plans: []PlanNode{
					{NodeType: "Seq Scan", RelationName: "c"},
					{NodeType: "Seq Scan", RelationName: "d"},
				},
			},
		},
	}

	AssignIDs(&root, new(int))

	want := map[string]int{
		"root": 1,
		"a":    3,
		"b":    4,
		"c":    6,
		"d":    7,
	}

	got := map[string]int{
		"root": root.ID,
		"a":    root.Plans[0].Plans[0].ID,
		"b":    root.Plans[0].Plans[1].ID,
		"c":    root.Plans[1].Plans[0].ID,
		"d":    root.Plans[1].Plans[1].ID,
	}

	for k, wantID := range want {
		if got[k] != wantID {
			t.Fatalf("ID do nó %q: got %d, want %d", k, got[k], wantID)
		}
	}

	if root.Plans[0].ID != 2 {
		t.Fatalf("ID do primeiro Nested Loop filho: got %d, want 2", root.Plans[0].ID)
	}
	if root.Plans[1].ID != 5 {
		t.Fatalf("ID do segundo Nested Loop filho: got %d, want 5", root.Plans[1].ID)
	}

	// IDs precisam ser únicos.
	seen := map[int]bool{}
	for _, id := range got {
		if seen[id] {
			t.Fatalf("ID %d repetido entre nós diferentes", id)
		}
		seen[id] = true
	}
}

// TestAssignIDsThenFindBottlenecksSameIDs garante que o ID atribuído antes
// do pipeline é o mesmo que aparece no resumo de gargalos (FindBottlenecks)
// e na árvore — os dois consomem o mesmo ponteiro já numerado, sem
// reatribuir.
func TestAssignIDsThenFindBottlenecksSameIDs(t *testing.T) {
	root := PlanNode{
		NodeType:        "Nested Loop",
		ActualTotalTime: 500.5,
		ActualLoops:     1,
		Plans: []PlanNode{
			{NodeType: "Seq Scan", RelationName: "outer", ActualTotalTime: 200, ActualLoops: 1},
			{NodeType: "Index Scan", RelationName: "inner", ActualTotalTime: 0.5, ActualLoops: 1000},
		},
	}

	AssignIDs(&root, new(int))

	wantInnerID := root.Plans[1].ID // 3

	bs := FindBottlenecks(&root, 700.5, 1)
	if len(bs) != 1 {
		t.Fatalf("esperava 1 bottleneck, veio %d", len(bs))
	}
	if bs[0].Node.RelationName != "inner" {
		t.Fatalf("esperava 'inner' no topo, veio %q", bs[0].Node.RelationName)
	}
	if bs[0].Node.ID != wantInnerID {
		t.Fatalf("ID do bottleneck != ID atribuído antes do pipeline: got %d, want %d", bs[0].Node.ID, wantInnerID)
	}
}
