# sqlplan

CLI em Go que roda `EXPLAIN (ANALYZE, FORMAT JSON)` num Postgres e mostra o plano de execução como árvore colorida, destacando nós suspeitos (seq scan pesado, estimativa desatualizada).

## Uso

```bash
sqlplan --dsn "postgres://user:pass@localhost:5432/dbname" --query "SELECT * FROM users WHERE id = 1"
sqlplan --dsn "postgres://..." --file query.sql
sqlplan --dsn "postgres://..." --query "..." --top 5
```

- `--dsn`: obrigatório, ou lido de `PGDSN` se ausente
- `--query` / `--file`: exatamente um dos dois
- `--top`: quantos nós aparecem no resumo de gargalos (padrão 3, mínimo 1)
- Sem `--dsn` ou sem `--query`/`--file`, ou `--top` < 1: mostra uso da CLI (exit code 2)

### Exemplo de saída

```
Top 3 gargalos (de 2.400ms totais):
  1. Hash Join                           1.203ms  (50.1%)
  2. Seq Scan on orders                  0.891ms  (37.1%)
  3. Seq Scan on users                   0.306ms  (12.8%)

└── Hash Join (cost=200.00..400.00 rows=500) (actual rows=480 time=5.120ms, 100.0% do total)
    ├── Seq Scan on orders (cost=0.00..180.00 rows=8000) (actual rows=8000 time=1.900ms, 37.1% do total) ⚠ seq scan retornando muitas linhas
    └── Hash (cost=100.00..100.00 rows=500) (actual rows=500 time=0.800ms, 15.6% do total)
```

O resumo ranqueia por **tempo exclusivo** (self time) de cada nó, não pelo `Actual Total Time` cru do Postgres — esse campo é cumulativo (inclui os filhos) e é média por execução quando `Actual Loops > 1`. Ordenar pelo campo cru sempre aponta pra raiz da árvore e sub-estima nós dentro de nested loop; self time é a técnica usada por explain.depesz.com e pgMustard.

Cor desliga sozinha quando a saída não é terminal (pipe/redirect pra arquivo).

## Setup local

Banco de teste via Docker:

```bash
docker compose up
```

Sobe Postgres 17 em `localhost:5432` (user/pass/db = `sqlplan`).

```bash
go build -o sqlplan.exe .
go test ./...
```

## Estrutura

```
sqlplan/
├── main.go            # entrypoint, trata erro de uso vs erro geral
├── cmd/root.go         # flags, validação, orquestração
├── db/connect.go        # conexão pgx (sem pool, abre/roda/fecha)
├── explain/
│   ├── run.go           # executa EXPLAIN ANALYZE, retorna JSON cru
│   ├── plan.go           # structs ExplainResult / PlanNode
│   ├── parse.go          # unmarshal do JSON cru
│   └── analyze.go        # self time + ranking de gargalos (sprint 4)
└── render/
    ├── tree.go           # desenho de árvore + cor
    ├── highlight.go       # regras de destaque de nós suspeitos
    └── summary.go         # bloco "Top N gargalos" (sprint 4)
```

## Dependências

- [`github.com/jackc/pgx/v5`](https://github.com/jackc/pgx) - driver Postgres
- [`github.com/fatih/color`](https://github.com/fatih/color) — cor ANSI com detecção automática de terminal

## Decisões de design

- Sem pool de conexão: é execução única (abre, roda, fecha), pool só faria sentido se existisse modo watch
- Sem framework de CLI (Cobra): `flag` da standard library resolve pro escopo atual
- Sem Bubble Tea/TUI: ferramenta imprime uma vez e sai, não tem loop de eventos nem interação pós-resultado
- Sem timeout no context da query (`ponytail` comentado em `cmd/root.go`): `EXPLAIN ANALYZE` roda a query de verdade e pode demorar; timeout fica pra quando isso incomodar de fato
- Ranking de gargalos por self time, não por `Actual Total Time` bruto: o campo do Postgres é cumulativo (soma os filhos) e vira média por execução quando o nó roda em loop — ordenar por ele sempre aponta a raiz da árvore em vez do nó que pesou de verdade
