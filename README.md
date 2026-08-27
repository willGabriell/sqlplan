# sqlplan

CLI em Go que roda `EXPLAIN (ANALYZE, FORMAT JSON)` num Postgres e mostra o plano de execução como árvore colorida, destacando nós suspeitos (seq scan pesado, estimativa desatualizada).

## Uso

```bash
sqlplan --dsn "postgres://user:pass@localhost:5432/dbname" --query "SELECT * FROM users WHERE id = 1"
sqlplan --dsn "postgres://..." --file query.sql
```

- `--dsn`: obrigatório, ou lido de `PGDSN` se ausente
- `--query` / `--file`: exatamente um dos dois
- Sem `--dsn` ou sem `--query`/`--file`: mostra uso da CLI (exit code 2)

### Exemplo de saída

```
└── Hash Join (cost=200.00..400.00 rows=500) (actual rows=480 time=5.120ms, 100.0% do total)
    ├── Seq Scan on orders (cost=0.00..180.00 rows=8000) (actual rows=8000 time=1.900ms, 37.1% do total) ⚠  seq scan retornando muitas linhas
    └── Hash (cost=100.00..100.00 rows=500) (actual rows=500 time=0.800ms, 15.6% do total)
```

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
│   └── parse.go          # unmarshal do JSON cru
└── render/
    ├── tree.go           # desenho de árvore + cor
    └── highlight.go       # regras de destaque de nós suspeitos
```

## Dependências

- [`github.com/jackc/pgx/v5`](https://github.com/jackc/pgx) - driver Postgres
- [`github.com/fatih/color`](https://github.com/fatih/color) — cor ANSI com detecção automática de terminal

## Decisões de design

- Sem pool de conexão: é execução única (abre, roda, fecha), pool só faria sentido se existisse modo watch
- Sem framework de CLI (Cobra): `flag` da standard library resolve pro escopo atual
- Sem Bubble Tea/TUI: ferramenta imprime uma vez e sai, não tem loop de eventos nem interação pós-resultado
- Sem timeout no context da query (`ponytail` comentado em `cmd/root.go`): `EXPLAIN ANALYZE` roda a query de verdade e pode demorar; timeout fica pra quando isso incomodar de fato
