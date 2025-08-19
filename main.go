package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sahatsawats/mysql-align/db"
	"github.com/sahatsawats/mysql-align/models"
	"github.com/sahatsawats/mysql-align/features"
	"github.com/sahatsawats/mysql-align/utils"


)

// Make docker to test
func main() {
	const version string = "v1.00"

	if len(os.Args) <2 {
		fmt.Println("Usage: myalign <command> [args]")
	}

	// First args
	cmd := os.Args[1]
	switch cmd {
	case "version":
		fmt.Println(version)
	case "recon-row":
		var resultsReport []models.InformationSchema

		// CMD arguments
		recconcileCmd := flag.NewFlagSet("recconcile", flag.ExitOnError)
		user := recconcileCmd.String("user", "root", "User to access database")
		pwd := recconcileCmd.String("password", "", "Password for user to access database")
		host := recconcileCmd.String("host", "localhost", "Hostname or IP-Address to database server")
		port := recconcileCmd.Int("port", 3306, "Port of database server")
		output := recconcileCmd.String("output", "", "Path to output csv file.")
		recconcileCmd.Parse(os.Args[2:])

		// Check output file is not empty.
		if *pwd == "" {
			fmt.Println("Please specify password for user with --password <password>")
			os.Exit(1)
		}
		if *output == "" {
			fmt.Println("Please specify output path for csv file with --output <file-path>")
			os.Exit(1)
		}

		// initialize database connection. return conn object
		conn, err := db.InitializeDB(host, port, user, pwd)
		defer conn.Close()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		
		resultsReport, err = features.Reconcile(conn)
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

		os.Exit(0)
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

		conn, err := db.InitializeDB(host, port, user, pwd)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		defer conn.Close()


	default:
		fmt.Println("Mismatch detected. Please choose a corrective action.")
	}


	
}