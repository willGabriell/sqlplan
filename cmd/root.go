// Package cmd cuida de flags, validação e orquestração da CLI.
package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"

	"sqlplan/db"
	"sqlplan/explain"
	"sqlplan/render"
)

// ErrUsage indica uso incorreto das flags (dispara exit code 2 + help).
var ErrUsage = errors.New("uso inválido")

// resolveQuery decide qual query rodar a partir das flags --query/--file.
// Extraída à parte por ser a única lógica que vale testar sem banco.
func resolveQuery(query, file string) (string, error) {
	if query != "" && file != "" {
		return "", fmt.Errorf("%w: use --query ou --file, não os dois", ErrUsage)
	}
	if query == "" && file == "" {
		return "", fmt.Errorf("%w: informe --query ou --file", ErrUsage)
	}

	q := query
	if file != "" {
		b, err := os.ReadFile(file)
		if err != nil {
			return "", fmt.Errorf("falha ao ler --file: %w", err)
		}
		q = string(b)
	}

	if strings.TrimSpace(q) == "" {
		return "", errors.New("query vazia")
	}
	return q, nil
}

// Execute faz parse das flags, roda o EXPLAIN e imprime o resultado.
func Execute() error {
	dsn := flag.String("dsn", os.Getenv("PGDSN"), "DSN de conexão com o Postgres (ou variável PGDSN)")
	query := flag.String("query", "", "query SQL a explicar")
	file := flag.String("file", "", "arquivo .sql com a query a explicar")
	top := flag.Int("top", 3, "quantos gargalos mostrar no resumo")
	flag.Parse()

	if *dsn == "" {
		return fmt.Errorf("%w: informe --dsn ou defina PGDSN", ErrUsage)
	}
	if *top < 1 {
		return fmt.Errorf("%w: --top precisa ser >= 1", ErrUsage)
	}

	q, err := resolveQuery(*query, *file)
	if err != nil {
		return err
	}

	// ponytail: sem timeout — EXPLAIN ANALYZE roda a query de verdade e pode
	// pendurar. Adicionar --timeout quando isso incomodar.
	ctx := context.Background()

	conn, err := db.Connect(ctx, *dsn)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	out, err := explain.Run(ctx, conn, q)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Position > 0 {
				return fmt.Errorf("erro do postgres (posição %d): %s", pgErr.Position, pgErr.Message)
			}
			return fmt.Errorf("erro do postgres: %s", pgErr.Message)
		}
		return fmt.Errorf("falha ao executar explain: %w", err)
	}

	res, err := explain.ParsePlan(out)
	if err != nil {
		return err
	}
	bs := explain.FindBottlenecks(&res.Plan, res.ExecutionTime, *top)
	render.Summary(os.Stdout, bs, res.ExecutionTime)
	render.Tree(os.Stdout, &res.Plan)
	return nil
}
