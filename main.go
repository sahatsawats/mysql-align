package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sahatsawats/mysql-align/db"
	"github.com/sahatsawats/mysql-align/features"
	"github.com/sahatsawats/mysql-align/models"
	"github.com/sahatsawats/mysql-align/utils"
)

const version = "v1.10"

type connFlags struct {
	user         string
	password     string
	host         string
	port         int
	serverPubKey string
}

func registerConnFlags(fs *flag.FlagSet) *connFlags {
	cf := &connFlags{}
	fs.StringVar(&cf.user,         "user",          "root",      "User to access database")
	fs.StringVar(&cf.password,     "password",      "",          "Password for user to access database")
	fs.StringVar(&cf.host,         "host",          "localhost", "Hostname or IP-Address to database server")
	fs.IntVar(   &cf.port,         "port",          3306,        "Port of database server")
	fs.StringVar(&cf.serverPubKey, "server-pub-key","",          "RSA file for transmit encryption data.")
	return cf
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: myalign <command> [args]")
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "version":
		fmt.Println(version)
		os.Exit(0)
	case "pre-migration":
		err = runPreMigration(os.Args[2:])
	case "recon-rows":
		err = runReconRows(os.Args[2:])
	case "recon-objs":
		err = runReconObjs(os.Args[2:])
	case "get-size":
		err = runGetSize(os.Args[2:])
	case "get-config":
		err = runGetConfig(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", os.Args[1])
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	endProgram()
	os.Exit(0)
}

func sectionStatus(err error, count int) string {
	if err != nil {
		return "failed"
	}
	if count > 0 {
		return "issues"
	}
	return "ok"
}

func runPreMigration(args []string) error {
	fs := flag.NewFlagSet("pre-migration", flag.ExitOnError)
	cf := registerConnFlags(fs)
	output := fs.String("output", "", "Output directory; HTML report is written as pre_migration_report_<timestamp>.html")
	maxRows := fs.Int("max-rows", 500, "Maximum rows per check section in HTML report (0 for unlimited)")
	debug := fs.Bool("debug", false, "enable debug log")
	fs.Parse(args)

	if *output == "" {
		return fmt.Errorf("pre-migration: --output is required")
	}

	runStart := time.Now()

	printLogo()

	if *debug {
		utils.SetDebug(true)
		utils.Debug("Debug mode enabled")
	}

	conn, err := db.InitializeDB(&cf.host, &cf.port, &cf.user, &cf.password, &cf.serverPubKey)
	if err != nil {
		return fmt.Errorf("pre-migration: connect: %w", err)
	}
	defer conn.Close()

	var serverVersion string
	conn.QueryRow("SELECT VERSION()").Scan(&serverVersion)

	fmt.Println("Checking upgrade compatability to MySQL Enterprise Edition 8.4.X")
	fmt.Println()

	var sections []models.CheckSection

	// CHAR_CHECK
	charSetReport, charSetErr := features.CheckCharSet(conn)
	if charSetErr != nil {
		fmt.Println("Error:", charSetErr)
	}
	if len(charSetReport) != 0 {
		fmt.Println("[CHAR_CHECK]: NOT OK. ( CHAR_CHECK_ERR_COUNT:", len(charSetReport), ")")
	} else {
		fmt.Println("[CHAR_CHECK]: OK")
	}
	{
		rows := make([][]string, len(charSetReport))
		for i, r := range charSetReport {
			rows[i] = []string{r.Severity, r.SchemaName, r.CharSet}
		}
		errStr := ""
		if charSetErr != nil {
			errStr = charSetErr.Error()
		}
		sections = append(sections, models.CheckSection{
			ID:      "CHAR_CHECK",
			Status:  sectionStatus(charSetErr, len(charSetReport)),
			Count:   len(charSetReport),
			Error:   errStr,
			Columns: []string{"SEVERITY", "SCHEMA_NAME", "CHARSET"},
			Rows:    rows,
		})
	}

	// ENGINE_CHECK
	engineReport, engineErr := features.CheckEngine(conn)
	if engineErr != nil {
		fmt.Println("Error:", engineErr)
	}
	if len(engineReport) != 0 {
		fmt.Println("[ENGINE_CHECK]: NOT OK. ( ENGINE_CHECK_ERR_COUNT:", len(engineReport), ")")
	} else {
		fmt.Println("[ENGINE_CHECK]: OK")
	}
	{
		rows := make([][]string, len(engineReport))
		for i, r := range engineReport {
			rows[i] = []string{r.SchemaName, r.TableName, r.Engine, r.CreateOptions}
		}
		errStr := ""
		if engineErr != nil {
			errStr = engineErr.Error()
		}
		sections = append(sections, models.CheckSection{
			ID:      "ENGINE_CHECK",
			Status:  sectionStatus(engineErr, len(engineReport)),
			Count:   len(engineReport),
			Error:   errStr,
			Columns: []string{"SCHEMA_NAME", "TABLE_NAME", "ENGINE", "CREATE_OPTIONS"},
			Rows:    rows,
		})
	}

	// ROW_F_CHECK
	rowFormatReport, rowFErr := features.CheckRowFormat(conn)
	if rowFErr != nil {
		fmt.Println("Error:", rowFErr)
	}
	if len(rowFormatReport) != 0 {
		fmt.Println("[ROW_F_CHECK]: NOT OK. ( ROW_F_CHECK_ERR_COUNT:", len(rowFormatReport), ")")
	} else {
		fmt.Println("[ROW_F_CHECK]: OK")
	}
	{
		rows := make([][]string, len(rowFormatReport))
		for i, r := range rowFormatReport {
			rows[i] = []string{r.SchemaName, r.TableName, r.Engine, r.RowFormat}
		}
		errStr := ""
		if rowFErr != nil {
			errStr = rowFErr.Error()
		}
		sections = append(sections, models.CheckSection{
			ID:      "ROW_F_CHECK",
			Status:  sectionStatus(rowFErr, len(rowFormatReport)),
			Count:   len(rowFormatReport),
			Error:   errStr,
			Columns: []string{"SCHEMA_NAME", "TABLE_NAME", "ENGINE", "ROW_FORMAT"},
			Rows:    rows,
		})
	}

	// PK_CHECK
	pkReport, pkErr := features.CheckNoPK(conn)
	if pkErr != nil {
		fmt.Println("Error:", pkErr)
	}
	if len(pkReport) != 0 {
		fmt.Println("[PK_CHECK]: NOT OK. ( PK_ERR_COUNT:", len(pkReport), ")")
	} else {
		fmt.Println("[PK_CHECK]: OK")
	}
	{
		rows := make([][]string, len(pkReport))
		for i, r := range pkReport {
			rows[i] = []string{r.SchemaName, r.TableName}
		}
		errStr := ""
		if pkErr != nil {
			errStr = pkErr.Error()
		}
		sections = append(sections, models.CheckSection{
			ID:      "PK_CHECK",
			Status:  sectionStatus(pkErr, len(pkReport)),
			Count:   len(pkReport),
			Error:   errStr,
			Columns: []string{"SCHEMA_NAME", "TABLE_NAME"},
			Rows:    rows,
		})
	}

	// FK_CHECK
	fkReport, fkErr := features.CheckFKDuplication(conn)
	if fkErr != nil {
		fmt.Println("Error:", fkErr)
	}
	if len(fkReport) != 0 {
		fmt.Println("[FK_CHECK]: NOT OK. ( FK_DUP_ERR_COUNT:", len(fkReport), ")")
	} else {
		fmt.Println("[FK_CHECK]: OK")
	}
	{
		rows := make([][]string, len(fkReport))
		for i, r := range fkReport {
			rows[i] = []string{r.SchemaName, r.ConstraintName, strconv.Itoa(r.Count)}
		}
		errStr := ""
		if fkErr != nil {
			errStr = fkErr.Error()
		}
		sections = append(sections, models.CheckSection{
			ID:      "FK_CHECK",
			Status:  sectionStatus(fkErr, len(fkReport)),
			Count:   len(fkReport),
			Error:   errStr,
			Columns: []string{"SCHEMA_NAME", "CONSTRAINT_NAME", "OCCURRENCES"},
			Rows:    rows,
		})
	}

	// VIEW_CHECK
	viewReport, viewErr := features.CheckViewDeprecated(conn)
	if viewErr != nil {
		fmt.Println("Error:", viewErr)
	}
	if len(viewReport) != 0 {
		fmt.Println("[VIEW_CHECK]: NOT OK. ( VIEW_CHECK_ERR_COUNT:", len(viewReport), ")")
	} else {
		fmt.Println("[VIEW_CHECK]: OK")
	}
	{
		rows := make([][]string, len(viewReport))
		for i, r := range viewReport {
			rows[i] = []string{r.SchemaName, r.TableName, r.ViewDefinition}
		}
		errStr := ""
		if viewErr != nil {
			errStr = viewErr.Error()
		}
		sections = append(sections, models.CheckSection{
			ID:      "VIEW_CHECK",
			Status:  sectionStatus(viewErr, len(viewReport)),
			Count:   len(viewReport),
			Error:   errStr,
			Columns: []string{"SCHEMA_NAME", "TABLE_NAME", "VIEW_DEFINITION"},
			Rows:    rows,
		})
	}

	// ROUTINE_SYNTAX_CHECK
	syntaxReport, syntaxErr := features.CheckRoutineSyntaxDeprecated(conn)
	if syntaxErr != nil {
		fmt.Println("Error:", syntaxErr)
	}
	if len(syntaxReport) != 0 {
		fmt.Println("[ROUTINE_SYNTAX_CHECK]: NOT OK. ( ROUTINE_SYNTAX_ERR_COUNT:", len(syntaxReport), ")")
	} else {
		fmt.Println("[ROUTINE_SYNTAX_CHECK]: OK")
	}
	{
		rows := make([][]string, len(syntaxReport))
		for i, r := range syntaxReport {
			rows[i] = []string{r.SchemaName, r.RoutineName, r.RoutineType, r.RoutineDefinition}
		}
		errStr := ""
		if syntaxErr != nil {
			errStr = syntaxErr.Error()
		}
		sections = append(sections, models.CheckSection{
			ID:      "ROUTINE_SYNTAX_CHECK",
			Status:  sectionStatus(syntaxErr, len(syntaxReport)),
			Count:   len(syntaxReport),
			Error:   errStr,
			Columns: []string{"SCHEMA_NAME", "ROUTINE_NAME", "ROUTINE_TYPE", "ROUTINE_DEFINITION"},
			Rows:    rows,
		})
	}

	// ROUTINE_FUNC_CHECK
	funcReport, funcErr := features.CheckRoutineFunctionDeprecated(conn)
	if funcErr != nil {
		fmt.Println("Error:", funcErr)
	}
	if len(funcReport) != 0 {
		fmt.Println("[ROUTINE_FUNC_CHECK]: NOT OK. ( ROUTINE_FUNC_ERR_COUNT:", len(funcReport), ")")
	} else {
		fmt.Println("[ROUTINE_FUNC_CHECK]: OK")
	}
	{
		rows := make([][]string, len(funcReport))
		for i, r := range funcReport {
			rows[i] = []string{r.SchemaName, r.RoutineName, r.RoutineType, r.RoutineDefinition}
		}
		errStr := ""
		if funcErr != nil {
			errStr = funcErr.Error()
		}
		sections = append(sections, models.CheckSection{
			ID:      "ROUTINE_FUNC_CHECK",
			Status:  sectionStatus(funcErr, len(funcReport)),
			Count:   len(funcReport),
			Error:   errStr,
			Columns: []string{"SCHEMA_NAME", "ROUTINE_NAME", "ROUTINE_TYPE", "ROUTINE_DEFINITION"},
			Rows:    rows,
		})
	}

	report := models.PreMigrationReport{
		Meta: models.ReportMeta{
			ToolVersion:          version,
			GeneratedAt:          runStart,
			GeneratedAtFormatted: runStart.Format("2006-01-02 15:04:05"),
			Host:                 cf.host,
			Port:                 cf.port,
			User:                 cf.user,
			ServerVersion:        serverVersion,
		},
		Sections: sections,
	}

	if err := os.MkdirAll(*output, 0o755); err != nil {
		return fmt.Errorf("pre-migration: create output dir: %w", err)
	}

	filename := fmt.Sprintf("pre_migration_report_%s.html", runStart.Format("2006-01-02_15-04-05"))
	outPath := filepath.Join(*output, filename)

	if err := utils.WritePreMigrationHTML(report, outPath, *maxRows); err != nil {
		return fmt.Errorf("pre-migration: write report: %w", err)
	}

	return nil
}

func runReconRows(args []string) error {
	fs := flag.NewFlagSet("recon-rows", flag.ExitOnError)
	cf := registerConnFlags(fs)
	output := fs.String("output", "", "Path to output csv file.")
	debug := fs.Bool("debug", false, "enable debug log")
	fs.Parse(args)

	if *output == "" {
		return fmt.Errorf("recon-rows: --output is required")
	}

	printLogo()

	if *debug {
		utils.SetDebug(true)
		utils.Debug("Debug mode enabled")
	}

	conn, err := db.InitializeDB(&cf.host, &cf.port, &cf.user, &cf.password, &cf.serverPubKey)
	if err != nil {
		return fmt.Errorf("recon-rows: connect: %w", err)
	}
	defer conn.Close()

	results, err := features.ReconcileRow(conn)
	if err != nil {
		return fmt.Errorf("recon-rows: reconcile rows: %w", err)
	}
	fmt.Println("total rows:", len(results))

	if err := utils.SaveInformationTablesToCSV(results, *output); err != nil {
		return fmt.Errorf("recon-rows: write csv: %w", err)
	}

	return nil
}

func runReconObjs(args []string) error {
	fs := flag.NewFlagSet("recon-objs", flag.ExitOnError)
	cf := registerConnFlags(fs)
	output := fs.String("output", "", "Path to output csv file.")
	debug := fs.Bool("debug", false, "enable debug log")
	fs.Parse(args)

	if *output == "" {
		return fmt.Errorf("recon-objs: --output is required")
	}

	printLogo()

	if *debug {
		utils.SetDebug(true)
		utils.Debug("Debug mode enabled")
	}

	conn, err := db.InitializeDB(&cf.host, &cf.port, &cf.user, &cf.password, &cf.serverPubKey)
	if err != nil {
		return fmt.Errorf("recon-objs: connect: %w", err)
	}
	defer conn.Close()

	results, err := features.ReconcileObject(conn)
	if err != nil {
		return fmt.Errorf("recon-objs: reconcile objects: %w", err)
	}
	fmt.Println("total rows:", len(results))

	if err := utils.SaveInformationObjectToCSV(results, *output); err != nil {
		return fmt.Errorf("recon-objs: write csv: %w", err)
	}

	return nil
}

func runGetSize(args []string) error {
	fs := flag.NewFlagSet("get-size", flag.ExitOnError)
	cf := registerConnFlags(fs)
	output := fs.String("output", "", "Path to output csv file.")
	debug := fs.Bool("debug", false, "enable debug log")
	fs.Parse(args)

	if *output == "" {
		return fmt.Errorf("get-size: --output is required")
	}

	printLogo()

	if *debug {
		utils.SetDebug(true)
		utils.Debug("Debug mode enabled")
	}

	conn, err := db.InitializeDB(&cf.host, &cf.port, &cf.user, &cf.password, &cf.serverPubKey)
	if err != nil {
		return fmt.Errorf("get-size: connect: %w", err)
	}
	defer conn.Close()

	results, err := features.GetSchemaSize(conn)
	if err != nil {
		return fmt.Errorf("get-size: query size: %w", err)
	}
	fmt.Println("total rows:", len(results))

	if err := utils.SizeToCSV(results, *output); err != nil {
		return fmt.Errorf("get-size: write csv: %w", err)
	}

	return nil
}

func runGetConfig(args []string) error {
	fs := flag.NewFlagSet("get-config", flag.ExitOnError)
	cf := registerConnFlags(fs)
	output := fs.String("output", "", "Path to output csv file.")
	debug := fs.Bool("debug", false, "enable debug log")
	fs.Parse(args)

	if *output == "" {
		return fmt.Errorf("get-config: --output is required")
	}

	printLogo()

	if *debug {
		utils.SetDebug(true)
		utils.Debug("Debug mode enabled")
	}

	conn, err := db.InitializeDB(&cf.host, &cf.port, &cf.user, &cf.password, &cf.serverPubKey)
	if err != nil {
		return fmt.Errorf("get-config: connect: %w", err)
	}
	defer conn.Close()

	results, err := features.GetConfiguration(conn)
	if err != nil {
		return fmt.Errorf("get-config: query configuration: %w", err)
	}

	if err := utils.SaveServerConfigurationToCSV(results, *output); err != nil {
		return fmt.Errorf("get-config: write csv: %w", err)
	}

	return nil
}

func printLogo() {
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Print(`
=======================================================================================================================================
   _____ _          _     _                _         ____                                          _     _           _ _           _
  |  ___(_)_ __ ___| |_  | |    ___   __ _(_) ___   / ___|___  _ __ ___  _ __   __ _ _ __  _   _  | |   (_)_ __ ___ (_) |_ ___  __| |
  | |_  | | '__/ __| __| | |   / _ \ / _  | |/ __| | |   / _ \| '_   _ \| '_ \ / _  | '_ \| | | | | |   | | '_   _ \| | __/ _ \/ _  |
  |  _| | | |  \__ \ |_  | |__| (_) | (_| | | (__  | |__| (_) | | | | | | |_) | (_| | | | | |_| | | |___| | | | | | | | ||  __/ (_| |
  |_|   |_|_|  |___/\__| |_____\___/ \__, |_|\___|  \____\___/|_| |_| |_| .__/ \__,_|_| |_|\__, | |_____|_|_| |_| |_|_|\__\___|\__,_|
                                     |___/                              |_|                |___/
=======================================================================================================================================`)
	fmt.Println()
	fmt.Println(" © 2025 First Logic Company")
	fmt.Println(" All rights reserved. Proprietary software.")
	fmt.Println(" This script is for migration purposes ONLY.")
	fmt.Println(" Unauthorized use, copying, or distribution is strictly prohibited.")
	fmt.Println("________________________________________________________________________________________________________________________________________")
	timestart := fmt.Sprintf("                                              Program starting at %s                 \n", now)
	fmt.Print(timestart)
}

func endProgram() {
	now := time.Now().Format("2006-01-02 15:04:05")
	timeEnd := fmt.Sprintf("                                              Program stopped at %s                 \n", now)
	fmt.Print(timeEnd)
	fmt.Println("________________________________________________________________________________________________________________________________________")
}
