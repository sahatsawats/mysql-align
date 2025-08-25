package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sahatsawats/mysql-align/db"
	"github.com/sahatsawats/mysql-align/features"
	"github.com/sahatsawats/mysql-align/models"
	"github.com/sahatsawats/mysql-align/utils"
)

// Make docker to test
func main() {
	const version string = "v1.06"
	
	if len(os.Args) < 2 {
		fmt.Println("Usage: myalign <command> [args]")
	}

	// First args
	cmd := os.Args[1]
	switch cmd {
	case "version":
		fmt.Println(version)
		os.Exit(0)
	case "pre-migration":
		// CMD arguments
		CMD := flag.NewFlagSet("recconcile", flag.ExitOnError)
		user := CMD.String("user", "root", "User to access database")
		pwd := CMD.String("password", "", "Password for user to access database")
		host := CMD.String("host", "localhost", "Hostname or IP-Address to database server")
		port := CMD.Int("port", 3306, "Port of database server")
		serverPubPath := CMD.String("server-pub-key", "", "RSA file for transmit encryption data.")
		output := CMD.String("output-path", "", "output directory")
		CMD.Parse(os.Args[2:])

		if *output == "" {
			fmt.Println("Please specify output path for csv file with --output <file-path>")
			os.Exit(1)
		}
		
		printLogo()


		// initialize database connection. return conn object
		conn, err := db.InitializeDB(host, port, user, pwd, serverPubPath)
		defer conn.Close()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		
		fmt.Println("Checking upgrade compatability to MySQL Enterprise Edition 8.4.X")
		fmt.Println()
		
		charSetReport, err := features.CheckCharSet(conn)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		if len(charSetReport) != 0 {
			fmt.Println("[CHAR_CHECK]: NOT OK. ( CHAR_CHECK_ERR_COUNT:", len(charSetReport), ")")
			err := utils.CharSetReportToCSV(charSetReport, *output)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("[CHAR_CHECK]: OK")
		}

		engineReport, err := features.CheckEngine(conn)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		if len(engineReport) != 0 {
			fmt.Println("[ENGINE_CHECK]: NOT OK. ( ENGINE_CHECK_ERR_COUNT:", len(engineReport), ")")
			err := utils.EngineReportToCSV(engineReport, *output)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("[ENGINE_CHECK]: OK")
		}

		rowFormatReport, err := features.CheckRowFormat(conn)
		if err != nil {
			fmt.Println("Error: ", err)

		}
		if len(rowFormatReport) != 0 {
			fmt.Println("[ROW_F_CHECK]: NOT OK. ( ROW_F_CHECK_ERR_COUNT:", len(engineReport), ")")
			err := utils.RowFormatReportToCSV(rowFormatReport, *output)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("[ROW_F_CHECK]: OK")
		}


		pkReport, err := features.CheckNoPK(conn)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		if len(pkReport) != 0 {
			fmt.Println("[PK_CHECK]: NOT OK. ( PK_ERR_COUNT:", len(pkReport), ")")
			err := utils.PKReportToCSV(pkReport, *output)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("[PK_CHECK]: OK")
		}

		fkReport, err := features.CheckFKDuplication(conn)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		if fkReport != 0 {
			fmt.Println("[FK_CHECK]: NOT OK. ( FK_DUP_ERR_COUNT:", fkReport, ")")
		} else {
			fmt.Println("[FK_CHECK]: OK")
		}

		viewReport, err := features.CheckViewDeprecated(conn)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		//fmt.Println(viewReport)
		if len(viewReport) != 0 {
			fmt.Println("[VIEW_CHECK]: NOT OK. ( VIEW_CHECK_ERR_COUNT:", len(viewReport), ")")
			err := utils.ViewReportToCSV(viewReport, *output)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("[VIEW_CHECK]: OK")
		}

		syntaxRoutineReport, err := features.CheckRoutineSyntaxDeprecated(conn)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		if len(syntaxRoutineReport) != 0 {
			fmt.Println("[ROUTINE_SYNTAX_CHECK]: NOT OK. ( ROUTINE_SYNTAX_ERR_COUNT:", len(syntaxRoutineReport), ")")
			err := utils.SyntaxRoutineToCSV(syntaxRoutineReport, *output)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("[ROUTINE_SYNTAX_CHECK]: OK")
		}
		
		functionRoutineReport, err := features.CheckRoutineFunctionDeprecated(conn)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		if len(functionRoutineReport) != 0 {
			fmt.Println("[ROUTINE_FUNC_CHECK]: NOT OK. ( ROUTINE_FUNC_ERR_COUNT:", len(functionRoutineReport), ")")
			err := utils.FunctionRoutineToCSV(functionRoutineReport, *output)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("[ROUTINE_FUNC_CHECK]: OK")
		}
		
	case "recon-rows":
		var resultsReport []models.InformationSchema

		// CMD arguments
		CMD := flag.NewFlagSet("recconcile", flag.ExitOnError)
		user := CMD.String("user", "root", "User to access database")
		pwd := CMD.String("password", "", "Password for user to access database")
		host := CMD.String("host", "localhost", "Hostname or IP-Address to database server")
		port := CMD.Int("port", 3306, "Port of database server")
		serverPubPath := CMD.String("server-pub-key", "", "RSA file for transmit encryption data.")

		output := CMD.String("output", "", "Path to output csv file.")
		CMD.Parse(os.Args[2:])

		if *output == "" {
			fmt.Println("Please specify output path for csv file with --output <file-path>")
			os.Exit(1)
		}

		printLogo()

		// initialize database connection. return conn object
		conn, err := db.InitializeDB(host, port, user, pwd, serverPubPath)
		defer conn.Close()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		resultsReport, err = features.ReconcileRow(conn)
		if err != nil {
			errorMsg := fmt.Sprintf("Error on recconcile process: %s", err.Error())
			fmt.Println(errorMsg)
			os.Exit(1)
		}
		fmt.Println("total rows: ", len(resultsReport))

		// Dumping data into CSV
		err = utils.SaveInformationTablesToCSV(resultsReport, *output)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

	case "recon-objs":
		var resultsReport []models.InformationObject
		// CMD arguments
		CMD := flag.NewFlagSet("recconcile", flag.ExitOnError)
		user := CMD.String("user", "root", "User to access database")
		pwd := CMD.String("password", "", "Password for user to access database")
		host := CMD.String("host", "localhost", "Hostname or IP-Address to database server")
		port := CMD.Int("port", 3306, "Port of database server")
		serverPubPath := CMD.String("server-pub-key", "", "RSA file for transmit encryption data.")

		output := CMD.String("output", "", "Path to output csv file.")
		CMD.Parse(os.Args[2:])
		
		if *output == "" {
			fmt.Println("Please specify output path for csv file with --output <file-path>")
			os.Exit(1)
		}

		printLogo()

		// initialize database connection. return conn object
		conn, err := db.InitializeDB(host, port, user, pwd, serverPubPath)
		defer conn.Close()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		resultsReport, err = features.ReconcileObject(conn)
		if err != nil {
			errorMsg := fmt.Sprintf("Error on recconcile process: %s", err.Error())
			fmt.Println(errorMsg)
			os.Exit(1)
		}
		fmt.Println("total rows: ", len(resultsReport))

		// Dumping data into CSV
		err = utils.SaveInformationObjectToCSV(resultsReport, *output)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}


	case "get-config":
		//var configs []models.InformationConfig

		// CMD arguments
		CMD := flag.NewFlagSet("recconcile", flag.ExitOnError)
		user := CMD.String("user", "root", "User to access database")
		pwd := CMD.String("password", "", "Password for user to access database")
		host := CMD.String("host", "localhost", "Hostname or IP-Address to database server")
		port := CMD.Int("port", 3306, "Port of database server")
		serverPubPath := CMD.String("server-pub-key", "", "RSA file for transmit encryption data.")
		output := CMD.String("output", "", "Path to output csv file.")
		CMD.Parse(os.Args[2:])

		// Check output file is not empty.
		if *output == "" {
			fmt.Println("Please specify output path for csv file with --output <file-path>")
			os.Exit(1)
		}


		printLogo()

		conn, err := db.InitializeDB(host, port, user, pwd, serverPubPath)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		defer conn.Close()

		resultsReport, err := features.GetConfiguration(conn)
		if err != nil {
			errorMsg := fmt.Sprintf("Error on recconcile process: %s", err.Error())
			fmt.Println(errorMsg)
			os.Exit(1)
		}

		err = utils.SaveServerConfigurationToCSV(resultsReport, *output)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	default:
		fmt.Println("Mismatch detected. Please choose a corrective action.")
		os.Exit(1)
	}

	endProgram()
	os.Exit(0)

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