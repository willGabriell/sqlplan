// Package explain executa EXPLAIN ANALYZE e devolve o JSON bruto formatado.
package explain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// Run roda EXPLAIN (ANALYZE, FORMAT JSON) na query e devolve o JSON indentado.
func Run(ctx context.Context, conn *pgx.Conn, query string) ([]byte, error) {
	sql := "EXPLAIN (ANALYZE, FORMAT JSON) " + query

	var raw []byte
	if err := conn.QueryRow(ctx, sql).Scan(&raw); err != nil {
		return nil, err // errors.As em cmd/root.go extrai *pgconn.PgError
	}

	var buf bytes.Buffer
	if err := json.Indent(&buf, raw, "", "  "); err != nil {
		return nil, fmt.Errorf("resposta do postgres não é JSON válido: %w", err)
	}
	return buf.Bytes(), nil
}
