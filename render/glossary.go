package render

import (
	"fmt"
	"io"
)

// glossary explica cada NodeType numa frase — usado só quando --glossary
// filtra pros tipos que de fato apareceram na execução atual.
var glossary = map[string]string{
	"Seq Scan":          "percorre a tabela inteira, linha por linha, sem usar índice",
	"Index Scan":        "busca linhas por um índice e vai até a tabela buscar o restante dos dados",
	"Index Only Scan":   "busca linhas só no índice, sem precisar acessar a tabela",
	"Bitmap Heap Scan":  "lê da tabela as linhas marcadas por um bitmap montado a partir de índice(s)",
	"Bitmap Index Scan": "monta um bitmap de páginas/linhas candidatas a partir de um índice",
	"Nested Loop":       "para cada linha do lado externo, refaz a busca no lado interno",
	"Hash Join":         "monta uma hash table de um lado e casa com as linhas do outro lado",
	"Merge Join":        "casa dois conjuntos já ordenados percorrendo os dois em paralelo",
	"Sort":              "ordena as linhas recebidas antes de repassá-las adiante",
	"Materialize":       "guarda o resultado do filho em memória pra reler sem reexecutar",
	"Unique":            "remove linhas duplicadas consecutivas de um conjunto já ordenado",
	"Aggregate":         "calcula agregações (soma, contagem, etc.) sobre as linhas recebidas",
	"Hash":              "monta a hash table usada por um Hash Join",
	"Gather":            "coleta e junta os resultados de workers paralelos",
	"WindowAgg":         "calcula funções de janela (window functions) sobre as linhas",
	"Subquery Scan":     "expõe o resultado de uma subquery como se fosse uma tabela",
	"Gather Merge":      "coleta resultados já ordenados de workers paralelos preservando a ordem",
	"Incremental Sort":  "ordena aproveitando uma ordenação parcial já existente nas linhas de entrada",
}

// Glossary escreve o bloco final explicando só os NodeType em types que têm
// entrada no map — tipo sem entrada é omitido em silêncio, sem erro ou
// placeholder. Nada é impresso se nenhum tipo sobrar.
func Glossary(w io.Writer, types []string) {
	type entry struct {
		nodeType, explanation string
	}
	var entries []entry
	for _, t := range types {
		if exp, ok := glossary[t]; ok {
			entries = append(entries, entry{t, exp})
		}
	}
	if len(entries) == 0 {
		return
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "Termos usados neste plano:")
	for _, e := range entries {
		fmt.Fprintf(w, "  %s: %s\n", e.nodeType, e.explanation)
	}
}
