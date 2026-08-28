package explain

// ExplainResult é o elemento único do array que EXPLAIN (FORMAT JSON) devolve.
type ExplainResult struct {
	Plan          PlanNode `json:"Plan"`
	PlanningTime  float64  `json:"Planning Time"`
	ExecutionTime float64  `json:"Execution Time"`
}

// PlanNode é um nó da árvore de execução. Campos como RelationName e
// IndexName só existem em certos tipos de nó (scan/index scan) — ficam
// zero value quando ausentes.
type PlanNode struct {
	NodeType          string  `json:"Node Type"`
	RelationName      string  `json:"Relation Name,omitempty"`
	Alias             string  `json:"Alias,omitempty"`
	IndexName         string  `json:"Index Name,omitempty"`
	StartupCost       float64 `json:"Startup Cost"`
	TotalCost         float64 `json:"Total Cost"`
	PlanRows          int64   `json:"Plan Rows"`
	PlanWidth         int64   `json:"Plan Width"`
	ActualStartupTime float64 `json:"Actual Startup Time"`
	ActualTotalTime   float64 `json:"Actual Total Time"`
	// ActualRows é float64, não int64: Postgres 17 emite valor fracionário
	// (ex: 9987.00) quando Actual Loops > 1, e int64 quebraria o unmarshal.
	ActualRows  float64    `json:"Actual Rows"`
	ActualLoops int64      `json:"Actual Loops"`
	Plans       []PlanNode `json:"Plans,omitempty"`

	// ID não vem do Postgres — atribuído por AssignIDs pra cruzar o resumo
	// de gargalos com o nó certo na árvore.
	ID int `json:"-"`
}

// AssignIDs numera os nós em pré-ordem (raiz, depois filhos da esquerda pra
// direita) começando em 1 — mesma ordem em que render.Tree imprime, pra que
// o ID do resumo de gargalos aponte pro nó certo da árvore.
func AssignIDs(n *PlanNode, counter *int) {
	*counter++
	n.ID = *counter
	for i := range n.Plans {
		AssignIDs(&n.Plans[i], counter)
	}
}
