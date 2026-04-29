//go:build integration

package features

import (
	"database/sql"
	"os"
	"strconv"
	"testing"

	"github.com/sahatsawats/mysql-align/db"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	host := envOrDefault("MYSQL_HOST", "127.0.0.1")
	portStr := envOrDefault("MYSQL_PORT", "3307")
	user := envOrDefault("MYSQL_USER", "root")
	password := envOrDefault("MYSQL_PASSWORD", "testpw")
	serverPubKey := ""

	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic("invalid MYSQL_PORT: " + portStr)
	}

	conn, err := db.InitializeDB(&host, &port, &user, &password, &serverPubKey)
	if err != nil {
		panic("failed to connect to test MySQL: " + err.Error())
	}
	testDB = conn
	defer testDB.Close()

	os.Exit(m.Run())
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
