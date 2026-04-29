// Must use html/template (not text/template) — view/routine bodies are user SQL and
// rely on context-aware auto-escaping for safety.
package utils

import (
	_ "embed"
	"fmt"
	"html/template"
	"os"

	"github.com/sahatsawats/mysql-align/models"
)

//go:embed templates/report.tmpl.html
var reportTmplSrc string

//go:embed templates/report.css
var reportCSSSrc string

type checkSectionMeta struct {
	Name        string
	Description string
	DefaultSev  string
	CodeCols    []bool
}

var checkMeta = map[string]checkSectionMeta{
	"CHAR_CHECK": {
		Name:        "Character Set",
		Description: "Identifies schemas not using the utf8mb4 character set. Schemas using utf8 (aliased to utf8mb3 in MySQL 8.4) are flagged as ERROR; schemas using latin1 are flagged as WARNING. Both must be converted to utf8mb4 before upgrading.",
		DefaultSev:  "ERROR",
		CodeCols:    []bool{false, false, false},
	},
	"ENGINE_CHECK": {
		Name:        "Storage Engine",
		Description: "Identifies tables using a storage engine other than InnoDB (MyISAM, Memory, FEDERATED). These engines have limited or deprecated support in MySQL 8.4 Enterprise Edition.",
		DefaultSev:  "WARNING",
		CodeCols:    []bool{false, false, false, false},
	},
	"ROW_F_CHECK": {
		Name:        "Row Format",
		Description: "Identifies base tables using a legacy row format (Redundant, Compact, or Fixed). MySQL 8.4 defaults to DYNAMIC; these legacy formats may cause issues during import, export, or replication.",
		DefaultSev:  "WARNING",
		CodeCols:    []bool{false, false, false, false},
	},
	"PK_CHECK": {
		Name:        "Primary Key",
		Description: "Identifies base tables with no primary key. Tables without a primary key are incompatible with row-based replication and can severely degrade performance on Group Replication.",
		DefaultSev:  "ERROR",
		CodeCols:    []bool{false, false},
	},
	"FK_CHECK": {
		Name:        "Foreign Key Duplication",
		Description: "Identifies duplicate foreign key constraint names within a schema. Duplicate FK constraint names across tables in the same schema can cause import failures in MySQL 8.4.",
		DefaultSev:  "ERROR",
		CodeCols:    []bool{false, false, false},
	},
	"VIEW_CHECK": {
		Name:        "View Syntax",
		Description: "Identifies views that use the deprecated GROUP BY ASC/DESC syntax, which was removed in MySQL 8.0. These views must be rewritten before migration.",
		DefaultSev:  "WARNING",
		CodeCols:    []bool{false, false, true},
	},
	"ROUTINE_SYNTAX_CHECK": {
		Name:        "Routine Syntax",
		Description: "Identifies stored procedures and functions that use the deprecated GROUP BY ASC/DESC syntax, which was removed in MySQL 8.0.",
		DefaultSev:  "WARNING",
		CodeCols:    []bool{false, false, false, true},
	},
	"ROUTINE_FUNC_CHECK": {
		Name:        "Routine Functions",
		Description: "Identifies stored procedures and functions that call the deprecated built-in functions DECODE(), ENCODE(), or COMPRESS(). These functions are removed in MySQL 8.4.",
		DefaultSev:  "ERROR",
		CodeCols:    []bool{false, false, false, true},
	},
}

const maxCellLen = 2000

func truncateCell(s string) string {
	if len(s) > maxCellLen {
		return s[:maxCellLen] + " … [truncated]"
	}
	return s
}

func WritePreMigrationHTML(report models.PreMigrationReport, outPath string, maxRows int) error {
	sections := make([]models.CheckSection, len(report.Sections))
	for i, sec := range report.Sections {
		meta, ok := checkMeta[sec.ID]
		if ok {
			sec.Name = meta.Name
			sec.Description = meta.Description
			sec.DefaultSev = meta.DefaultSev
			sec.CodeCols = meta.CodeCols
			// Truncate long SQL cells for code columns.
			for j, row := range sec.Rows {
				newRow := make([]string, len(row))
				copy(newRow, row)
				for k, cell := range newRow {
					if k < len(meta.CodeCols) && meta.CodeCols[k] {
						newRow[k] = truncateCell(cell)
					}
				}
				sec.Rows[j] = newRow
			}
		}
		// Apply per-section row cap.
		if maxRows > 0 && len(sec.Rows) > maxRows {
			sec.TotalCount = sec.Count
			sec.Rows = sec.Rows[:maxRows]
			sec.Truncated = true
		}
		sections[i] = sec
	}

	overall := computeOverall(sections)

	funcMap := template.FuncMap{
		"severityClass": func(sev string) string {
			switch sev {
			case "ERROR":
				return "pill-error"
			case "WARNING":
				return "pill-warning"
			default:
				return "pill-ok"
			}
		},
		"statusClass": func(status string) string {
			switch status {
			case "issues":
				return "pill-warning"
			case "failed":
				return "pill-failed"
			default:
				return "pill-ok"
			}
		},
		"statusLabel": func(status string) string {
			switch status {
			case "issues":
				return "ISSUES"
			case "failed":
				return "FAILED"
			default:
				return "OK"
			}
		},
		"overallClass": func(s string) string {
			switch s {
			case "ALL CLEAR":
				return "all-clear"
			case "WARNINGS":
				return "warnings"
			case "ERRORS":
				return "errors"
			default:
				return "checks-failed"
			}
		},
	}

	tmpl, err := template.New("report").Funcs(funcMap).Parse(reportTmplSrc)
	if err != nil {
		return fmt.Errorf("html report: parse template: %w", err)
	}

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("html report: create file: %w", err)
	}
	defer f.Close()

	data := struct {
		Meta          models.ReportMeta
		Sections      []models.CheckSection
		CSS           template.CSS
		OverallStatus string
	}{
		Meta:          report.Meta,
		Sections:      sections,
		CSS:           template.CSS(reportCSSSrc),
		OverallStatus: overall,
	}

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("html report: render: %w", err)
	}

	fmt.Println("HTML report created successfully:", outPath)
	return nil
}

func computeOverall(sections []models.CheckSection) string {
	hasFailed := false
	hasError := false
	hasWarning := false
	for _, s := range sections {
		switch s.Status {
		case "failed":
			hasFailed = true
		case "issues":
			if s.DefaultSev == "ERROR" {
				hasError = true
			} else {
				hasWarning = true
			}
		}
	}
	if hasFailed {
		return "CHECKS FAILED"
	}
	if hasError {
		return "ERRORS"
	}
	if hasWarning {
		return "WARNINGS"
	}
	return "ALL CLEAR"
}
