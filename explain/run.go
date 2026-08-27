// Package explain executa EXPLAIN ANALYZE e devolve o JSON bruto do plano.
package explain

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Run roda EXPLAIN (ANALYZE, FORMAT JSON) na query e devolve o JSON cru
// (array com um elemento), pra ParsePlan consumir.
func Run(ctx context.Context, conn *pgx.Conn, query string) ([]byte, error) {
	sql := "EXPLAIN (ANALYZE, FORMAT JSON) " + query

	var raw []byte
	if err := conn.QueryRow(ctx, sql).Scan(&raw); err != nil {
		return nil, err // errors.As em cmd/root.go extrai *pgconn.PgError
	}
	return raw, nil
}
