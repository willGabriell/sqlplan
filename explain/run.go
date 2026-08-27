// Package explain executa EXPLAIN ANALYZE e devolve o JSON bruto do plano.
package explain

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

// Run roda EXPLAIN (ANALYZE, FORMAT JSON) na query e devolve o JSON cru
// (array com um elemento), pra ParsePlan consumir.
//
// EXPLAIN ANALYZE executa a query de verdade — se for UPDATE/DELETE/INSERT,
// altera dados. Por isso roda sempre dentro de transação com rollback
// explícito no final, nunca commit; isso resolve pra qualquer tipo de
// comando sem precisar fazer parsing de SQL pra distinguir SELECT de escrita.
func Run(ctx context.Context, conn *pgx.Conn, query string) ([]byte, error) {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("iniciar transação: %w", err)
	}

	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && !errors.Is(rbErr, pgx.ErrTxClosed) {
			fmt.Fprintf(os.Stderr, "aviso: falha ao reverter transação: %v\n", rbErr)
		}
	}()

	sql := "EXPLAIN (ANALYZE, FORMAT JSON) " + query

	var raw []byte
	if err := tx.QueryRow(ctx, sql).Scan(&raw); err != nil {
		return nil, err // errors.As em cmd/root.go extrai *pgconn.PgError
	}
	return raw, nil
}
