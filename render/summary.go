package render

import (
	"fmt"
	"io"

	"sqlplan/explain"
)

// Summary escreve o bloco "Top N gargalos" antes da árvore — os nós que
// mais pesaram por tempo exclusivo (self time), não tempo cumulativo.
func Summary(w io.Writer, bs []explain.Bottleneck, totalExecutionTime float64) {
	if len(bs) == 0 {
		return
	}

	fmt.Fprintf(w, "Top %d gargalos (de %.3fms totais):\n", len(bs), totalExecutionTime)

	for i, b := range bs {
		label := fmt.Sprintf("[%d] %s", b.Node.ID, nodeLabel(b.Node))
		line := fmt.Sprintf("  %d. %-32s %8.3fms  (%4.1f%%)", i+1, label, b.SelfTime, b.Percent)
		if c, ok := baseColor[b.Node.NodeType]; ok {
			c.Fprintln(w, line)
		} else {
			fmt.Fprintln(w, line)
		}
	}

	fmt.Fprintln(w)
}
