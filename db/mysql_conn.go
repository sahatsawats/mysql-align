package db

import (
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

// Build DSN, prerequisite check. Return pointer of database or Error.
func InitializeDB(host *string, port *int, user *string, password *string, serverPubPath *string) (*sql.DB, error) {
	var err error
	var publicKeyName string
	var dsn string
	// if public key path is empty -> no use of public key.
	if *serverPubPath != "" {
		publicKeyName, err = registerServerPubKey(*serverPubPath)
		if err != nil {
			errMsg := fmt.Errorf("fail to register public key: %s", err.Error())
			return nil, errMsg
		}
		// add public key options
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/?tls=false&serverPubKey=%s", *user, *password, *host, *port, publicKeyName)
	} else {
		// no public key options
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/", *user, *password, *host, *port)
	}

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

func registerServerPubKey(filePath string) (string, error) {
	const publicKeyName string = "mysql-enterprise-key"
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Decode file
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PUBLIC KEY" {
		return "", fmt.Errorf("failed to decode PEM block containing public key")
	}

	// Parse data into pub structure
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// Register public key to driver
	if rsaPubKey, ok := pub.(*rsa.PublicKey); ok {
		mysql.RegisterServerPubKey(publicKeyName, rsaPubKey)
	} else {
		return "", fmt.Errorf("not a RSA public key")
	}

	return publicKeyName, nil
}



