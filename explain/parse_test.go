package explain

import "testing"

func TestParsePlan(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		wantErr bool
	}{
		{name: "array vazio", in: `[]`, wantErr: true},
		{name: "malformado", in: `{`, wantErr: true},
		{
			name: "actual rows fracionário (PG17 com loops>1)",
			in: `[{"Plan":{"Node Type":"Seq Scan","Relation Name":"users",
				"Startup Cost":0,"Total Cost":1,"Plan Rows":1,"Plan Width":1,
				"Actual Startup Time":0,"Actual Total Time":0,
				"Actual Rows":9987.00,"Actual Loops":3},
				"Planning Time":0.1,"Execution Time":0.2}]`,
		},
		{
			name: "join aninhado 3 níveis",
			in: `[{"Plan":{
				"Node Type":"Hash Join","Startup Cost":200,"Total Cost":400,
				"Plan Rows":500,"Plan Width":32,
				"Actual Startup Time":1,"Actual Total Time":5.12,
				"Actual Rows":480,"Actual Loops":1,
				"Plans":[
					{"Node Type":"Seq Scan","Relation Name":"orders",
					 "Startup Cost":0,"Total Cost":180,"Plan Rows":8000,"Plan Width":16,
					 "Actual Startup Time":0.01,"Actual Total Time":1.9,
					 "Actual Rows":8000,"Actual Loops":1},
					{"Node Type":"Hash","Startup Cost":100,"Total Cost":100,
					 "Plan Rows":500,"Plan Width":16,
					 "Actual Startup Time":0.8,"Actual Total Time":0.8,
					 "Actual Rows":500,"Actual Loops":1,
					 "Plans":[
						{"Node Type":"Seq Scan","Relation Name":"users",
						 "Startup Cost":0,"Total Cost":155,"Plan Rows":10000,"Plan Width":32,
						 "Actual Startup Time":0.012,"Actual Total Time":2.341,
						 "Actual Rows":9987,"Actual Loops":1}
					 ]}
				]},
				"Planning Time":0.123,"Execution Time":2.4}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ParsePlan([]byte(tt.in))
			if tt.wantErr {
				if err == nil {
					t.Fatalf("esperava erro, veio nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if res == nil {
				t.Fatal("resultado nil sem erro")
			}
		})
	}

	t.Run("estrutura do join aninhado", func(t *testing.T) {
		res, err := ParsePlan([]byte(tests[3].in))
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if len(res.Plan.Plans) != 2 {
			t.Fatalf("esperava 2 filhos no Hash Join, veio %d", len(res.Plan.Plans))
		}
		hash := res.Plan.Plans[1]
		if len(hash.Plans) != 1 || hash.Plans[0].RelationName != "users" {
			t.Fatalf("esperava neto Seq Scan on users, veio %+v", hash.Plans)
		}
	})
}
