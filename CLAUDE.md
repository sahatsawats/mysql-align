# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build — output goes to build/, filename includes the version const from main.go
# Bash / Git Bash
VERSION=$(grep -oP 'const version = "\K[^"]+' main.go) && go build -o "build/myalign-test-${VERSION}.exe" .

# PowerShell
$v = (Select-String 'const version = "(.+)"' main.go).Matches.Groups[1].Value; go build -o "build/myalign-test-$v.exe" .

# Run
./build/myalign-test-<version>.exe <command> [options]

# Tidy dependencies
go mod tidy

# Vet
go vet ./...
```

There are no tests in this project.

## Architecture

This is a CLI tool written in Go that connects to a MySQL server and produces CSV reports and a consolidated HTML pre-migration report for database auditing and migration readiness. It has no config files — all connection parameters are passed as CLI flags at runtime.

**Entry point:** `main.go` — parses the first positional argument as the subcommand, then uses `flag.FlagSet` to parse remaining args. Each subcommand opens a DB connection via `db.InitializeDB`, calls a function from `features/`, and writes output via `utils/`.

**Packages:**

- `db/` — `InitializeDB` builds a DSN and returns `*sql.DB`. Optionally registers an RSA public key (via `--server-pub-key`) for encrypted authentication against MySQL Enterprise.
- `features/` — one file per feature group; each function accepts `*sql.DB` and returns a typed slice + error. Functions query `information_schema` / `performance_schema` directly.
- `models/` — plain structs that mirror query result shapes (`models.go`) plus rendering-only aggregate types for the HTML report (`report.go`). No ORM.
- `utils/` — `csv_report.go` has one `*ToCSV` function per CSV report type; `html_report.go` renders the consolidated pre-migration HTML report from an embedded template under `utils/templates/`; `debug.go` is a package-level toggle (`SetDebug` / `Debug`).

**Subcommands and their output:**

| Command | Feature function | Output |
|---|---|---|
| `pre-migration` | Multiple checks in `migration.go` | Single timestamped HTML report `pre_migration_report_<ts>.html` in `--output` dir |
| `recon-rows` | `ReconcileRow` | Single CSV file at `--output` path |
| `recon-objs` | `ReconcileObject` | Single CSV file at `--output` path |
| `get-size` | `GetSchemaSize` | `schema_size.csv` in `--output` dir |
| `get-config` | `GetConfiguration` | Single CSV file at `--output` path |

**`pre-migration` checks** (consolidated into `<output>/pre_migration_report_<timestamp>.html`):
- `CHAR_CHECK` — schemas not using `utf8mb4` (ERROR for `utf8`, WARNING for `latin1`)
- `ENGINE_CHECK` — tables not using InnoDB (MyISAM, Memory, FEDERATED)
- `ROW_F_CHECK` — tables with Redundant/Compact/Fixed row format
- `PK_CHECK` — base tables with no primary key
- `FK_CHECK` — duplicate foreign key constraint names (schema, constraint name, occurrence count)
- `VIEW_CHECK` — views using deprecated `GROUP BY ASC/DESC` syntax
- `ROUTINE_SYNTAX_CHECK` — routines using deprecated `GROUP BY ASC/DESC`
- `ROUTINE_FUNC_CHECK` — routines using deprecated `DECODE/ENCODE/COMPRESS` functions

**`pre-migration` flags:** `--output` (required, output directory), `--max-rows` (default 500; max rows per check section in the HTML report; 0 for unlimited), `--debug`.

**Adding a new feature** follows the existing pattern:
1. Add a struct to `models/models.go`
2. Add a query function to the appropriate file in `features/` (or a new file)
3. For single-CSV output commands: add a `*ToCSV` function in `utils/csv_report.go`. For new pre-migration checks: add a new section entry in `utils/html_report.go` (`checkMeta` map) and update `report.tmpl.html` if needed.
4. Wire up a new `case` in the `switch` in `main.go`
