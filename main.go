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
	const version string = "v1.01"
	
	if len(os.Args) < 2 {
		fmt.Println("Usage: myalign <command> [args]")
	}

	// First args
	cmd := os.Args[1]
	switch cmd {
	case "version":
		fmt.Println(version)
		os.Exit(0)
	case "recon-row":
		var resultsReport []models.InformationSchema

		// CMD arguments
		reconcileCmd := flag.NewFlagSet("recconcile", flag.ExitOnError)
		user := reconcileCmd.String("user", "root", "User to access database")
		pwd := reconcileCmd.String("password", "", "Password for user to access database")
		host := reconcileCmd.String("host", "localhost", "Hostname or IP-Address to database server")
		port := reconcileCmd.Int("port", 3306, "Port of database server")
		serverPubPath := reconcileCmd.String("server-pub-key", "", "RSA file for transmite encryption data.")

		output := reconcileCmd.String("output", "", "Path to output csv file.")
		reconcileCmd.Parse(os.Args[2:])

		// Check output file is not empty.
		if *pwd == "" {
			fmt.Println("Please specify password for user with --password <password>")
			os.Exit(1)
		}
		if *serverPubPath == "" && *host != "localhost" {
			fmt.Println("Detect none-localhost but missing server public key. Please specify path to public key with --server-pub-key <path-to-pub-key>")
			os.Exit(1)
		}
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

	case "recon-object":
		var resultsReport []models.InformationObject
		// CMD arguments
		reconcileCmd := flag.NewFlagSet("recconcile", flag.ExitOnError)
		user := reconcileCmd.String("user", "root", "User to access database")
		pwd := reconcileCmd.String("password", "", "Password for user to access database")
		host := reconcileCmd.String("host", "localhost", "Hostname or IP-Address to database server")
		port := reconcileCmd.Int("port", 3306, "Port of database server")
		serverPubPath := reconcileCmd.String("server-pub-key", "", "RSA file for transmite encryption data.")

		output := reconcileCmd.String("output", "", "Path to output csv file.")
		reconcileCmd.Parse(os.Args[2:])

		// Check output file is not empty.
		if *pwd == "" {
			fmt.Println("Please specify password for user with --password <password>")
			os.Exit(1)
		}
		if *serverPubPath == "" && *host != "localhost" {
			fmt.Println("Detect none-localhost but missing server public key. Please specify path to public key with --server-pub-key <path-to-pub-key>")
			os.Exit(1)
		}
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
		reconcileCmd := flag.NewFlagSet("recconcile", flag.ExitOnError)
		user := reconcileCmd.String("user", "root", "User to access database")
		pwd := reconcileCmd.String("password", "", "Password for user to access database")
		host := reconcileCmd.String("host", "localhost", "Hostname or IP-Address to database server")
		port := reconcileCmd.Int("port", 3306, "Port of database server")
		serverPubPath := reconcileCmd.String("server-pub-key", "", "RSA file for transmite encryption data.")
		output := reconcileCmd.String("output", "", "Path to output csv file.")
		reconcileCmd.Parse(os.Args[2:])

		// Check output file is not empty.
		if *pwd == "" {
			fmt.Println("Please specify password for user with --password <password>")
			os.Exit(1)
		}
		if *serverPubPath == "" && *host != "localhost" {
			fmt.Println("Detect none-localhost but missing server public key. Please specify path to public key with --server-pub-key <path-to-pub-key>")
			os.Exit(1)
		}

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

	now := time.Now().Format("2006-01-02 15:04:05")
	timeEnd := fmt.Sprintf("                                              Program stopped at %s                 \n", now)
	fmt.Print(timeEnd)
	fmt.Println("________________________________________________________________________________________________________________________________________")
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