// Package render imprime o plano em árvore de texto simples, sem cor e
// sem desenho de galhos (isso é sprint 3).
package render

import (
	"fmt"
	"io"
	"strings"

	"sqlplan/explain"
)

// Tree escreve n e seus filhos em w, um nó por linha, indentado por
// profundidade (2 espaços por nível).
func Tree(w io.Writer, n *explain.PlanNode, depth int) {
	indent := strings.Repeat("  ", depth)

	name := n.NodeType
	switch {
	case n.IndexName != "" && n.RelationName != "":
		name += fmt.Sprintf(" using %s on %s", n.IndexName, n.RelationName)
	case n.RelationName != "":
		name += " on " + n.RelationName
	}

	fmt.Fprintf(w, "%s%s (cost=%.2f..%.2f rows=%d) (actual rows=%.0f time=%.3fms)\n",
		indent, name, n.StartupCost, n.TotalCost, n.PlanRows, n.ActualRows, n.ActualTotalTime)

	for i := range n.Plans {
		Tree(w, &n.Plans[i], depth+1)
	}
}
