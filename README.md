# mysql-align

A CLI tool for MySQL database auditing and migration readiness. It connects to a MySQL server, queries `information_schema` and `performance_schema`, and produces CSV reports or a consolidated HTML pre-migration report.

## Build

The binary is placed in the `build/` directory and named after the version constant in `main.go`.

```bash
# Bash / Git Bash
VERSION=$(grep -oP 'const version = "\K[^"]+' main.go) && go build -o "build/myalign-test-${VERSION}.exe" .

# PowerShell
$v = (Select-String 'const version = "(.+)"' main.go).Matches.Groups[1].Value; go build -o "build/myalign-test-$v.exe" .
```

Requires Go 1.23+ and a network-accessible MySQL server. No config files — all connection parameters are passed as flags at runtime.

## Usage

```
myalign <command> [flags]
```

### Connection flags

All commands except `version` accept the following flags:

| Flag | Default | Description |
|---|---|---|
| `--host` | `localhost` | Hostname or IP address of the MySQL server |
| `--port` | `3306` | TCP port |
| `--user` | `root` | MySQL user |
| `--password` | _(empty)_ | Password for the user |
| `--server-pub-key` | _(empty)_ | Path to an RSA public key PEM file (MySQL Enterprise encrypted auth) |

---

## Commands

### `version`

Print the tool version and exit.

```bash
myalign version
```

---

### `pre-migration`

Run a full set of MySQL 8.0 upgrade-compatibility checks and write a consolidated HTML report.

```bash
myalign pre-migration --host 192.168.1.10 --user admin --password secret --output ./reports
```

**Flags:**

| Flag | Default | Description |
|---|---|---|
| `--output` | _(required)_ | Output directory. Report is written as `pre_migration_report_<timestamp>.html` |
| `--max-rows` | `500` | Maximum rows per check section in the HTML report. Use `0` for unlimited. |
| `--debug` | `false` | Print debug log to stdout |

**Checks performed:**

| Check ID | What it finds | Severity |
|---|---|---|
| `CHAR_CHECK` | Schemas not using `utf8mb4` as default character set | ERROR for `utf8`; WARNING for `latin1`, `latin1_swedish_ci`, `utf8_general_ci` |
| `ENGINE_CHECK` | Tables using non-InnoDB storage engines | WARNING (MyISAM, MEMORY, FEDERATED) |
| `ROW_F_CHECK` | Tables using deprecated row formats | WARNING (Redundant, Compact, Fixed) |
| `PK_CHECK` | Base tables with no primary key | ERROR |
| `FK_CHECK` | Duplicate foreign key constraint names within a schema | WARNING |
| `VIEW_CHECK` | Views using deprecated `GROUP BY ... ASC/DESC` syntax | WARNING |
| `ROUTINE_SYNTAX_CHECK` | Stored procedures/functions using deprecated `GROUP BY ... ASC/DESC` | WARNING |
| `ROUTINE_FUNC_CHECK` | Stored procedures/functions using removed functions (`DECODE`, `ENCODE`, `COMPRESS`) | WARNING |

---

### `recon-rows`

Count rows in every base table across all non-system schemas and write a CSV.

```bash
myalign recon-rows --host 192.168.1.10 --user admin --password secret --output ./output/rows.csv
```

**Flags:**

| Flag | Default | Description |
|---|---|---|
| `--output` | _(required)_ | Output directory; file is written as `recon_rows_<timestamp>.csv` |
| `--debug` | `false` | Print debug log to stdout |

**Output columns:** `schema_name`, `table_name`, `rows`

---

### `recon-objs`

Inventory every database object (tables, views, triggers, procedures, functions, indexes) across all non-system schemas and write a CSV.

```bash
myalign recon-objs --host 192.168.1.10 --user admin --password secret --output ./output/objects.csv
```

**Flags:**

| Flag | Default | Description |
|---|---|---|
| `--output` | _(required)_ | Output directory; file is written as `recon_objects_<timestamp>.csv` |
| `--debug` | `false` | Print debug log to stdout |

**Output columns:** `object_type`, `schema_name`, `object_name`

---

### `get-size`

Report the total data + index size (MB) of every schema and write a CSV.

```bash
myalign get-size --host 192.168.1.10 --user admin --password secret --output ./output
```

**Flags:**

| Flag | Default | Description |
|---|---|---|
| `--output` | _(required)_ | Output directory; file is written as `schema_size_<timestamp>.csv` |
| `--debug` | `false` | Print debug log to stdout |

**Output columns:** `schema_name`, `size_mb`

> System schemas (`mysql`, `information_schema`, `performance_schema`, `sys`) are included in the output.

---

### `get-config`

Dump all global MySQL server variables from `performance_schema.global_variables` to a CSV.

```bash
myalign get-config --host 192.168.1.10 --user admin --password secret --output ./output/config.csv
```

**Flags:**

| Flag | Default | Description |
|---|---|---|
| `--output` | _(required)_ | Output directory; file is written as `server_config_<timestamp>.csv` |
| `--debug` | `false` | Print debug log to stdout |

**Output columns:** `variable_name`, `variable_value`

> Requires `SELECT` privilege on `performance_schema.global_variables`.

---

## Required MySQL privileges

The connecting user needs at minimum:

```sql
GRANT SELECT ON information_schema.* TO 'user'@'host';
GRANT SELECT ON performance_schema.global_variables TO 'user'@'host';
-- recon-rows also issues: SELECT COUNT(*) FROM <schema>.<table> for every base table
GRANT SELECT ON *.* TO 'user'@'host';
```

## Running integration tests

Integration tests require Docker and run against a real MySQL 5.7 container:

```bash
# Linux / macOS / Git Bash
bash scripts/test-integration.sh

# Windows PowerShell
powershell -File scripts\test-integration.ps1
```

Tests are opt-in via the `integration` build tag — `go test ./...` (no tag) passes with no database available.
