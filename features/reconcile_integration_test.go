//go:build integration

package features

import (
	"testing"
)

func TestReconcileRow(t *testing.T) {
	results, err := ReconcileRow(testDB)
	if err != nil {
		t.Fatalf("ReconcileRow error: %v", err)
	}

	wantRows := map[string]int{
		"fixture_rows.t_a": 3,
		"fixture_rows.t_b": 5,
	}
	found := map[string]bool{}
	for _, r := range results {
		key := r.SchemaName + "." + r.TableName
		if want, ok := wantRows[key]; ok {
			if r.Rows != want {
				t.Errorf("%s: expected %d rows, got %d", key, want, r.Rows)
			}
			found[key] = true
		}
	}
	for key := range wantRows {
		if !found[key] {
			t.Errorf("expected %s in ReconcileRow results", key)
		}
	}
}

func TestReconcileObject(t *testing.T) {
	results, err := ReconcileObject(testDB)
	if err != nil {
		t.Fatalf("ReconcileObject error: %v", err)
	}

	if len(results) < 1 {
		t.Fatalf("expected non-empty results from ReconcileObject")
	}

	// All rows are real data — no phantom header row.
	data := results
	wantPresent := []struct{ objType, schema, namePrefix string }{
		{"Table", "fixture_rows", "t_a"},
		{"Table", "fixture_rows", "t_b"},
		{"View", "fixture_views", "v_bad"},
		{"Procedure", "fixture_routines", "p_bad_groupby"},
		{"Function", "fixture_routines", "f_decode"},
	}
	for _, w := range wantPresent {
		found := false
		for _, r := range data {
			if r.ObjectType == w.objType && r.SchemaName == w.schema && r.ObjectName == w.namePrefix {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected object {%s, %s, %s} in ReconcileObject results", w.objType, w.schema, w.namePrefix)
		}
	}
}
