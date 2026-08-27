// Package render desenha o plano em árvore estilo `tree` do Unix, com cor
// por tipo de nó e destaque nos nós suspeitos (sprint 3).
package render

import (
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"

	"sqlplan/explain"
)

// baseColor mapeia tipo de nó pra cor base — mesmo tipo, mesma cor entre
// execuções diferentes.
var baseColor = map[string]*color.Color{
	"Seq Scan":          color.New(color.FgYellow),
	"Index Scan":        color.New(color.FgGreen),
	"Index Only Scan":   color.New(color.FgGreen),
	"Bitmap Heap Scan":  color.New(color.FgGreen),
	"Bitmap Index Scan": color.New(color.FgGreen),
	"Hash Join":         color.New(color.FgCyan),
	"Nested Loop":       color.New(color.FgCyan),
	"Merge Join":        color.New(color.FgCyan),
	"Sort":              color.New(color.FgBlue),
	"Aggregate":         color.New(color.FgBlue),
	"Hash":              color.New(color.FgBlue),
}

var highlightColor = color.New(color.FgRed, color.Bold)

// Tree escreve n e seus filhos em w, desenhados como árvore (`├──`/`└──`/`│`).
func Tree(w io.Writer, n *explain.PlanNode) {
	rootTime := n.ActualTotalTime
	printNode(w, n, rootTime, "", true, true)
}

func printNode(w io.Writer, n *explain.PlanNode, rootTime float64, prefix string, isLast, isRoot bool) {
	if !isRoot {
		branch := "├── "
		if isLast {
			branch = "└── "
		}
		fmt.Fprint(w, prefix+branch)
	}

	name := n.NodeType
	switch {
	case n.IndexName != "" && n.RelationName != "":
		name += fmt.Sprintf(" using %s on %s", n.IndexName, n.RelationName)
	case n.RelationName != "":
		name += " on " + n.RelationName
	}

	pct := 0.0
	if rootTime > 0 {
		pct = n.ActualTotalTime / rootTime * 100
	}

	line := fmt.Sprintf("%s (cost=%.2f..%.2f rows=%d) (actual rows=%.0f time=%.3fms, %.1f%% do total)",
		name, n.StartupCost, n.TotalCost, n.PlanRows, n.ActualRows, n.ActualTotalTime, pct)

	highlights := EvaluateNode(n)
	if len(highlights) > 0 {
		reasons := make([]string, len(highlights))
		for i, h := range highlights {
			reasons[i] = h.Reason
		}
		line += " ⚠ " + strings.Join(reasons, "; ")
		highlightColor.Fprintln(w, line)
	} else if c, ok := baseColor[n.NodeType]; ok {
		c.Fprintln(w, line)
	} else {
		fmt.Fprintln(w, line)
	}

	childPrefix := prefix
	if !isRoot {
		if isLast {
			childPrefix += "    "
		} else {
			childPrefix += "│   "
		}
	}

	for i := range n.Plans {
		printNode(w, &n.Plans[i], rootTime, childPrefix, i == len(n.Plans)-1, false)
	}
}
