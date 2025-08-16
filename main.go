package main

import (
	"fmt"
	"github.com/sahatsawats/mysql-align/db/mysql_conn"
)



func main() {
	var err error

	err = mysql_conn.InitializeDB("140.245.109.35", "3306", "root", "TEST@flg1234")
	if err != nil {
		fmt.Sprintf("%s", err)
	}
}