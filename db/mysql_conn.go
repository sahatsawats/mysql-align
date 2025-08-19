package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)


// Build DSN, prerequisite check. Return pointer of database or Error.
func InitializeDB(host *string, port *int, user *string, password *string) (*sql.DB, error) {
	var err error
	// $USER:$PWD@/tcp($HOST:$PORT)/$DATABASE)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", *user, *password, *host, *port)
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	} else {
		fmt.Println("Initialize database connection compeleted.")
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	} else {
		fmt.Println("Connected to MySQL.")
	}

	return conn, nil
}



