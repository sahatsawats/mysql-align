package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sahatsawats/mysql-align/db"
	"github.com/sahatsawats/mysql-align/features"
	"github.com/sahatsawats/mysql-align/utils"
)

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
	const version string = "v1.10"

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

func runPreMigration(args []string) error {
	fs := flag.NewFlagSet("pre-migration", flag.ExitOnError)
	cf := registerConnFlags(fs)
	output := fs.String("output", "", "output directory")
	debug := fs.Bool("debug", false, "enable debug log")
	fs.Parse(args)

	if *output == "" {
		return fmt.Errorf("pre-migration: --output is required")
	}

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

	fmt.Println("Checking upgrade compatability to MySQL Enterprise Edition 8.4.X")
	fmt.Println()

	charSetReport, err := features.CheckCharSet(conn)
	if err != nil {
		fmt.Println("Error:", err)
	}
	if len(charSetReport) != 0 {
		fmt.Println("[CHAR_CHECK]: NOT OK. ( CHAR_CHECK_ERR_COUNT:", len(charSetReport), ")")
		if err := utils.CharSetReportToCSV(charSetReport, *output); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("[CHAR_CHECK]: OK")
	}

	engineReport, err := features.CheckEngine(conn)
	if err != nil {
		fmt.Println("Error:", err)
	}
	if len(engineReport) != 0 {
		fmt.Println("[ENGINE_CHECK]: NOT OK. ( ENGINE_CHECK_ERR_COUNT:", len(engineReport), ")")
		if err := utils.EngineReportToCSV(engineReport, *output); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("[ENGINE_CHECK]: OK")
	}

	rowFormatReport, err := features.CheckRowFormat(conn)
	if err != nil {
		fmt.Println("Error:", err)
	}
	if len(rowFormatReport) != 0 {
		fmt.Println("[ROW_F_CHECK]: NOT OK. ( ROW_F_CHECK_ERR_COUNT:", len(rowFormatReport), ")")
		if err := utils.RowFormatReportToCSV(rowFormatReport, *output); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("[ROW_F_CHECK]: OK")
	}

	pkReport, err := features.CheckNoPK(conn)
	if err != nil {
		fmt.Println("Error:", err)
	}
	if len(pkReport) != 0 {
		fmt.Println("[PK_CHECK]: NOT OK. ( PK_ERR_COUNT:", len(pkReport), ")")
		if err := utils.PKReportToCSV(pkReport, *output); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("[PK_CHECK]: OK")
	}

	fkReport, err := features.CheckFKDuplication(conn)
	if err != nil {
		fmt.Println("Error:", err)
	}
	if fkReport != 0 {
		fmt.Println("[FK_CHECK]: NOT OK. ( FK_DUP_ERR_COUNT:", fkReport, ")")
	} else {
		fmt.Println("[FK_CHECK]: OK")
	}

	viewReport, err := features.CheckViewDeprecated(conn)
	if err != nil {
		fmt.Println("Error:", err)
	}
	if len(viewReport) != 0 {
		fmt.Println("[VIEW_CHECK]: NOT OK. ( VIEW_CHECK_ERR_COUNT:", len(viewReport), ")")
		if err := utils.ViewReportToCSV(viewReport, *output); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("[VIEW_CHECK]: OK")
	}

	syntaxRoutineReport, err := features.CheckRoutineSyntaxDeprecated(conn)
	if err != nil {
		fmt.Println("Error:", err)
	}
	if len(syntaxRoutineReport) != 0 {
		fmt.Println("[ROUTINE_SYNTAX_CHECK]: NOT OK. ( ROUTINE_SYNTAX_ERR_COUNT:", len(syntaxRoutineReport), ")")
		if err := utils.SyntaxRoutineToCSV(syntaxRoutineReport, *output); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("[ROUTINE_SYNTAX_CHECK]: OK")
	}

	functionRoutineReport, err := features.CheckRoutineFunctionDeprecated(conn)
	if err != nil {
		fmt.Println("Error:", err)
	}
	if len(functionRoutineReport) != 0 {
		fmt.Println("[ROUTINE_FUNC_CHECK]: NOT OK. ( ROUTINE_FUNC_ERR_COUNT:", len(functionRoutineReport), ")")
		if err := utils.FunctionRoutineToCSV(functionRoutineReport, *output); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("[ROUTINE_FUNC_CHECK]: OK")
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
