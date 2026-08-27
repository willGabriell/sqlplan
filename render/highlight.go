package render

import "sqlplan/explain"

// Limiares ajustáveis. Constantes num lugar só, não hardcoded espalhado.
const (
	// SeqScanRowsThreshold: acima disso um Seq Scan é suspeito de faltar índice.
	SeqScanRowsThreshold = 1000.0
	// EstimateMismatchRatio: razão actual/plan rows acima disso indica
	// estatística desatualizada (falta ANALYZE).
	EstimateMismatchRatio = 10.0
)

// Highlight é um motivo de destaque que bateu contra um PlanNode.
type Highlight struct {
	Reason string
}

// EvaluateNode roda todas as regras contra n e retorna os motivos que
// bateram. Um nó pode acumular mais de um highlight ao mesmo tempo.
func EvaluateNode(n *explain.PlanNode) []Highlight {
	var hs []Highlight

	if n.NodeType == "Seq Scan" && n.ActualRows > SeqScanRowsThreshold {
		hs = append(hs, Highlight{Reason: " seq scan retornando muitas linhas"})
	}

	planRows := float64(n.PlanRows)
	hi, lo := max(n.ActualRows, planRows), min(n.ActualRows, planRows)
	// max(1, lo) evita divisão por zero quando um dos dois lados é 0 e
	// ainda conta como mismatch de verdade (ex: PlanRows=10000, ActualRows=0).
	if hi/max(1, lo) > EstimateMismatchRatio {
		hs = append(hs, Highlight{Reason: " estimativa muito distante da real"})
	}

	return hs
}
