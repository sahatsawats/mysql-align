//go:build integration

package features

import (
	"testing"

	"github.com/sahatsawats/mysql-align/models"
)

func TestCheckNoPK(t *testing.T) {
	results, err := CheckNoPK(testDB)
	if err != nil {
		t.Fatalf("CheckNoPK error: %v", err)
	}

	if !hasNoPKEntry(results, "fixture_nopk", "events") {
		t.Errorf("expected fixture_nopk.events to be flagged as missing PK, got %v", results)
	}
	for _, r := range results {
		if r.SchemaName == "fixture_clean" {
			t.Errorf("fixture_clean should not appear in NoPK results, got %+v", r)
		}
	}
}

func TestCheckCharSet(t *testing.T) {
	results, err := CheckCharSet(testDB)
	if err != nil {
		t.Fatalf("CheckCharSet error: %v", err)
	}

	var foundError, foundWarning bool
	for _, r := range results {
		if r.SchemaName == "fixture_utf8" && r.Severity == "ERROR" {
			foundError = true
		}
		if r.SchemaName == "fixture_latin1" && r.Severity == "WARNING" {
			foundWarning = true
		}
		if r.SchemaName == "fixture_clean" {
			t.Errorf("fixture_clean (utf8mb4) should not appear in charset results, got %+v", r)
		}
	}
	if !foundError {
		t.Errorf("expected ERROR entry for fixture_utf8 (charset=utf8), got %v", results)
	}
	if !foundWarning {
		t.Errorf("expected WARNING entry for fixture_latin1, got %v", results)
	}
}

func TestCheckEngine(t *testing.T) {
	results, err := CheckEngine(testDB)
	if err != nil {
		t.Fatalf("CheckEngine error: %v", err)
	}

	var foundMyISAM bool
	for _, r := range results {
		if r.SchemaName == "fixture_engine" && r.TableName == "t_myisam" && r.Engine == "MyISAM" {
			foundMyISAM = true
		}
		if r.SchemaName == "fixture_clean" {
			t.Errorf("fixture_clean should not appear in engine results, got %+v", r)
		}
	}
	if !foundMyISAM {
		t.Errorf("expected fixture_engine.t_myisam (MyISAM) in engine results, got %v", results)
	}
}

func TestCheckRowFormat(t *testing.T) {
	results, err := CheckRowFormat(testDB)
	if err != nil {
		t.Fatalf("CheckRowFormat error: %v", err)
	}

	var foundCompact, foundRedundant bool
	for _, r := range results {
		if r.SchemaName == "fixture_rowfmt" && r.TableName == "t_compact" {
			foundCompact = true
		}
		if r.SchemaName == "fixture_rowfmt" && r.TableName == "t_redundant" {
			foundRedundant = true
		}
	}
	if !foundCompact {
		t.Errorf("expected fixture_rowfmt.t_compact in row format results, got %v", results)
	}
	if !foundRedundant {
		t.Errorf("expected fixture_rowfmt.t_redundant in row format results, got %v", results)
	}
}

func TestCheckFKDuplication(t *testing.T) {
	results, err := CheckFKDuplication(testDB)
	if err != nil {
		t.Fatalf("CheckFKDuplication error: %v", err)
	}

	// MySQL 5.7 InnoDB enforces unique FK names per schema at the engine level,
	// so CheckFKDuplication always returns empty on 5.7 — that is correct behavior.
	// Verify no false positive is raised for fixture_fkdup (which uses distinct FK names).
	for _, r := range results {
		if r.SchemaName == "fixture_fkdup" {
			t.Errorf("unexpected duplicate FK in fixture_fkdup: %+v", r)
		}
	}
}

func TestCheckViewDeprecated(t *testing.T) {
	results, err := CheckViewDeprecated(testDB)
	if err != nil {
		t.Fatalf("CheckViewDeprecated error: %v", err)
	}

	var found bool
	for _, r := range results {
		if r.SchemaName == "fixture_views" && r.TableName == "v_bad" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected fixture_views.v_bad in deprecated view results, got %v", results)
	}
}

func TestCheckRoutineSyntaxDeprecated(t *testing.T) {
	results, err := CheckRoutineSyntaxDeprecated(testDB)
	if err != nil {
		t.Fatalf("CheckRoutineSyntaxDeprecated error: %v", err)
	}

	var found bool
	for _, r := range results {
		if r.SchemaName == "fixture_routines" && r.RoutineName == "p_bad_groupby" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected fixture_routines.p_bad_groupby in syntax deprecated results, got %v", results)
	}
}

func TestCheckRoutineFunctionDeprecated(t *testing.T) {
	results, err := CheckRoutineFunctionDeprecated(testDB)
	if err != nil {
		t.Fatalf("CheckRoutineFunctionDeprecated error: %v", err)
	}

	var foundDecode, foundCompress bool
	for _, r := range results {
		if r.SchemaName == "fixture_routines" && r.RoutineName == "f_decode" {
			foundDecode = true
		}
		if r.SchemaName == "fixture_routines" && r.RoutineName == "p_compress" {
			foundCompress = true
		}
	}
	if !foundDecode {
		t.Errorf("expected fixture_routines.f_decode in function deprecated results, got %v", results)
	}
	if !foundCompress {
		t.Errorf("expected fixture_routines.p_compress in function deprecated results, got %v", results)
	}
}

func hasNoPKEntry(results []models.InformationNoPKTable, schema, table string) bool {
	for _, r := range results {
		if r.SchemaName == schema && r.TableName == table {
			return true
		}
	}
	return false
}
