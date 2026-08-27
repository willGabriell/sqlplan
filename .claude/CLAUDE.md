# CLAUDE.md — sqlplan

Contexto do projeto pra Claude Code (e qualquer agente) trabalhar aqui.

## O que é

CLI em Go (`sqlplan`) que roda `EXPLAIN (ANALYZE, FORMAT JSON)` num Postgres e imprime o plano
de execução como árvore colorida, com resumo dos maiores gargalos por self time e destaque de
nós suspeitos (seq scan pesado, estimativa desatualizada). Uso: `sqlplan --dsn ... --query "..."`
ou `--file query.sql`, flag `--top N` controla quantos gargalos aparecem no resumo. Ver
[README.md](../README.md) pra exemplos de uso e saída.

## Arquitetura

```
main.go            entrypoint — trata ErrUsage (exit 2 + usage) vs erro geral (exit 1)
cmd/root.go         flags, validação, orquestração do pipeline
db/connect.go        conexão pgx única (sem pool — abre, roda, fecha)
explain/
  run.go             executa EXPLAIN ANALYZE, devolve JSON cru
  plan.go            structs ExplainResult/PlanNode + AssignIDs (numeração pré-ordem)
  parse.go           unmarshal do JSON cru
  analyze.go         self time + ranking de gargalos (FindBottlenecks)
render/
  tree.go            desenho de árvore + cor + nodeLabel (compartilhado c/ summary)
  highlight.go        regras de destaque (EvaluateNode)
  summary.go          bloco "Top N gargalos"
```

Pipeline em `cmd/root.go`: `Run` → `ParsePlan` → `AssignIDs` (uma vez, antes de tudo) →
`FindBottlenecks` → `render.Summary` → `render.Tree`. `AssignIDs` precisa rodar antes dos dois
renderers pra que o mesmo nó tenha o mesmo ID no resumo e na árvore.

## Convenções de código

- **Self time, não Actual Total Time cru**: `Actual Total Time` do Postgres é cumulativo (inclui
  filhos) e vira média por execução quando `Actual Loops > 1`. Todo ranking/cálculo de "quanto
  esse nó pesou" usa self time (`explain/analyze.go`), nunca o campo bruto — ele sempre aponta pra
  raiz e sub-estima nós em loop.
- **IDs de nó são calculados, não vêm do Postgres**: `PlanNode.ID` tem `json:"-"` e só existe
  depois de `AssignIDs` rodar. Testes que montam `PlanNode` na mão e chamam `Tree`/`Summary`
  direto (sem passar pelo pipeline) precisam chamar `explain.AssignIDs(&n, new(int))` primeiro, ou
  o ID sai `0`.
- **`nodeLabel` é compartilhado** entre `render/tree.go` e `render/summary.go` — mudança no
  formato do label (tipo + relação/índice) afeta os dois. Cada um monta o prefixo `[N]` por conta
  própria porque o resumo alinha em coluna (`%-32s`) e a árvore não.
- **Sem pool de conexão, sem framework de CLI (Cobra), sem TUI**: execução é única (abre, roda,
  imprime, sai). `flag` da standard library resolve pro escopo atual. Não introduzir essas
  dependências sem necessidade concreta.
- **Sem timeout no context da query**: `EXPLAIN ANALYZE` roda a query de verdade e pode demorar;
  timeout fica pra quando isso incomodar de fato, não antes.

## Comentários

- Comentário só quando explica um **porquê não óbvio** — invariante escondido, workaround de bug
  específico, decisão que surpreenderia quem lê. Nunca comentário descrevendo o que o código já
  diz por si (nome de função/variável bem escolhido já basta).
- **Não usar marcadores de skill/ferramenta no comentário** (tipo `ponytail:` ou qualquer prefixo
  de plugin) — o comentário deve valer por si, sem depender de quem leu qual skill.
- Ao editar um comentário existente que cita uma dessas convenções (self time, sem pool, sem
  timeout), pode reescrever sem o prefixo, mantendo só a explicação técnica.

## Testes

- `go test ./...` roda tudo. Fixtures de `render/tree_test.go` e `render/summary_test.go` montam
  `PlanNode`/`Bottleneck` na mão — ao adicionar campo novo em `PlanNode` que afete output
  (como `ID`), atualizar as fixtures e os literais `want`.
- `explain/analyze_test.go` e `explain/ids_test.go` cobrem self time e numeração de IDs
  isoladamente, sem precisar de banco.
- Não existe teste de integração contra Postgre real neste repo — `docker compose up` sobe o
  banco só pra uso manual/exploratório (ver README).

## Git / commits

- **Nunca** adicionar Claude/Anthropic como co-autor do commit (sem linha
  `Co-Authored-By: Claude ...` ou similar). Mensagem de commit não menciona a ferramenta/IA usada
  pra gerar o código.
- Mensagens seguem `tipo: descrição curta` (`feat:`, `fix:`, `docs:`, etc.), em português, como no
  histórico existente (`git log`).
