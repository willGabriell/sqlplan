package explain

import "sort"

// Bottleneck é um nó da árvore com seu tempo exclusivo (self time) já
// calculado — quanto do tempo total da query foi gasto no próprio nó, sem
// contar filhos.
type Bottleneck struct {
	Node     *PlanNode
	SelfTime float64 // ms, exclusivo do nó, já multiplicado por loops
	Percent  float64 // % do Execution Time da query
}

// FindBottlenecks achata a árvore, calcula self time de cada nó (ver
// flatten) e devolve os topN maiores, ordenados decrescente.
//
// Actual Total Time do Postgres é cumulativo (inclui filhos) e é média por
// execução quando Actual Loops > 1 — por isso não dá pra rankear pelo campo
// cru direto, ele sempre aponta pra raiz. Self time = tempo do nó menos soma
// dos filhos, todos já multiplicados por loops.
func FindBottlenecks(root *PlanNode, totalExecutionTime float64, topN int) []Bottleneck {
	flat := flatten(root)

	sort.SliceStable(flat, func(i, j int) bool {
		return flat[i].SelfTime > flat[j].SelfTime
	})

	if topN < len(flat) {
		flat = flat[:topN]
	}

	for i := range flat {
		if totalExecutionTime > 0 {
			flat[i].Percent = flat[i].SelfTime / totalExecutionTime * 100
		}
	}

	return flat
}

func flatten(n *PlanNode) []Bottleneck {
	loops := max(1, n.ActualLoops)
	total := n.ActualTotalTime * float64(loops)

	var childrenTime float64
	var flat []Bottleneck

	for i := range n.Plans {
		child := &n.Plans[i]
		childLoops := max(1, child.ActualLoops)
		childrenTime += child.ActualTotalTime * float64(childLoops)
		flat = append(flat, flatten(child)...)
	}

	selfTime := total - childrenTime
	if selfTime < 0 {
		selfTime = 0 // guarda contra arredondamento de ponto flutuante
	}

	flat = append(flat, Bottleneck{Node: n, SelfTime: selfTime})
	return flat
}
