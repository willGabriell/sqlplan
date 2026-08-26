package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveQuery(t *testing.T) {
	fileWith := func(t *testing.T, content string) string {
		t.Helper()
		p := filepath.Join(t.TempDir(), "q.sql")
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		return p
	}

	tests := []struct {
		name      string
		query     string
		file      string
		wantQuery string
		wantErr   bool
	}{
		{name: "query direta", query: "select 1", wantQuery: "select 1"},
		{name: "via arquivo", file: fileWith(t, "select 2"), wantQuery: "select 2"},
		{name: "nenhum informado", wantErr: true},
		{name: "ambos informados", query: "select 1", file: fileWith(t, "select 2"), wantErr: true},
		{name: "arquivo inexistente", file: filepath.Join(t.TempDir(), "nope.sql"), wantErr: true},
		{name: "query só espaço em branco", query: "   ", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveQuery(tt.query, tt.file)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("esperava erro, veio nil (query=%q)", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if got != tt.wantQuery {
				t.Fatalf("got %q, want %q", got, tt.wantQuery)
			}
		})
	}
}
