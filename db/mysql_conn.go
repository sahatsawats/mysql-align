package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DSN string
var DB *sql.DB

// Build DSN, prerequisite check. Return Error.
func InitializeDB(host string, port string, user string, password string) error {
	var err error
	// $USER:$PWD@/tcp($HOST:$PORT)/$DATABASE)
	DSN = fmt.Sprintf("%s:%s@/tcp(%s:%s)/", user, password, host, port)
	DB, err = sql.Open("mysql", DSN)
	if err != nil {
		return err
	} else {
		fmt.Sprintf("Initialize database connection compeleted.")
	}

	if err := DB.Ping(); err != nil {
		return err
	} else {
		fmt.Sprintf("Connected to MySQL.")
	}

	return nil
}

