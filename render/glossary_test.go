package render

import (
	"bytes"
	"strings"
	"testing"
)

func TestGlossaryFiltersToTypesGiven(t *testing.T) {
	var buf bytes.Buffer
	Glossary(&buf, []string{"Nested Loop", "Materialize", "Index Scan"})

	got := buf.String()
	for _, want := range []string{"Nested Loop:", "Materialize:", "Index Scan:"} {
		if !strings.Contains(got, want) {
			t.Fatalf("esperava %q no output, got:\n%s", want, got)
		}
	}
	for _, unwanted := range []string{"Seq Scan:", "Hash Join:", "Sort:"} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("não esperava %q no output, got:\n%s", unwanted, got)
		}
	}
}

func TestGlossaryOmitsUnknownTypesSilently(t *testing.T) {
	var buf bytes.Buffer
	Glossary(&buf, []string{"Seq Scan", "Foo Bar Scan"})

	got := buf.String()
	if strings.Contains(got, "Foo Bar Scan") {
		t.Fatalf("tipo desconhecido não deveria aparecer, got:\n%s", got)
	}
	if !strings.Contains(got, "Seq Scan:") {
		t.Fatalf("esperava Seq Scan no output, got:\n%s", got)
	}
}

func TestGlossaryEmptyWhenNoKnownTypes(t *testing.T) {
	var buf bytes.Buffer
	Glossary(&buf, []string{"Foo Bar Scan", "Another Unknown"})

	if buf.String() != "" {
		t.Fatalf("esperava output vazio, got:\n%s", buf.String())
	}
}
