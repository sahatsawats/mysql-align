//go:build integration

package features

import (
	"testing"
)

func TestGetSchemaSize(t *testing.T) {
	results, err := GetSchemaSize(testDB)
	if err != nil {
		t.Fatalf("GetSchemaSize error: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("expected non-empty results from GetSchemaSize")
	}

	// GetSchemaSize does NOT filter system schemas; verify at least one system schema is present.
	systemSchemas := map[string]bool{"mysql": false, "information_schema": false, "performance_schema": false}
	fixtureFound := false
	for _, r := range results {
		if _, ok := systemSchemas[r.SchemaName]; ok {
			systemSchemas[r.SchemaName] = true
		}
		if r.SchemaName == "fixture_clean" {
			fixtureFound = true
		}
	}
	for name, found := range systemSchemas {
		if !found {
			t.Errorf("expected system schema %q in GetSchemaSize results (function does not filter system schemas)", name)
		}
	}
	if !fixtureFound {
		t.Errorf("expected fixture_clean in GetSchemaSize results")
	}
}

func TestGetConfiguration(t *testing.T) {
	results, err := GetConfiguration(testDB)
	if err != nil {
		t.Fatalf("GetConfiguration error: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("expected non-empty results from GetConfiguration")
	}

	var foundVersion bool
	for _, r := range results {
		if r.VariableName == "version" {
			foundVersion = true
		}
	}
	if !foundVersion {
		t.Error("expected 'version' variable in GetConfiguration results")
	}
}
