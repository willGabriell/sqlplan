// Package db cuida da conexão com o Postgres.
package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// Connect abre uma conexão única (sem pool: essa CLI abre, roda, fecha).
// Pool fica pra quando existir modo watch.
func Connect(ctx context.Context, dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar: %w", err)
	}
	return conn, nil
}
