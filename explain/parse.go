package explain

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ParsePlan faz unmarshal do JSON cru devolvido por Run (array com um
// elemento) e retorna esse elemento. Erro explícito se o JSON vier
// malformado ou vazio, em vez de panic de index out of range.
func ParsePlan(raw []byte) (*ExplainResult, error) {
	var results []ExplainResult
	if err := json.Unmarshal(raw, &results); err != nil {
		return nil, fmt.Errorf("json do explain inválido: %w", err)
	}
	if len(results) == 0 {
		return nil, errors.New("explain retornou array vazio")
	}
	return &results[0], nil
}
