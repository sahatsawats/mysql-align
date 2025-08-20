package features

import (
	"database/sql"

	"github.com/sahatsawats/mysql-align/models"
)

func CheckCharSet(conn *sql.DB) ([]models.CharSetObject, error) {
	const statement string = `SELECT schema_name, default_character_set_name 
	FROM information_schema.schemata WHERE schema_name NOT IN 
	('mysql', 'information_schema', 'sys', 'performance_schema') AND default_character_set_name != 'utf8mb4'`
	var listOfObject []models.CharSetObject
	var warningList = []string{"latin1", "latin1_swedish_ci", "utf8_general_ci"}
	var errorList = []string{"utf8"}

	// create empty map where key is string, and vaule is struct{} (which struct case zero memory)
	var errorSet = make(map[string]struct{})
	for _, e := range errorList {
		errorSet[e] = struct{}{}
	}
	var warnSet = make(map[string]struct{})
	for _, w := range warningList {
		warnSet[w] = struct{}{}
	}

	// query
	rows, err := conn.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// loop through each result
	for rows.Next() {
		// scan query results
		var item models.CharSetObject
		err := rows.Scan(&item.SchemaName, &item.CharSet)
		if err != nil {
			return nil, err
		}

		if _, found := errorSet[item.CharSet]; found {
			item.Severity = "ERROR"
			listOfObject = append(listOfObject, item)
		} else if _, found := warnSet[item.CharSet]; found {
			item.Severity = "WARNING"
			listOfObject = append(listOfObject, item)
		}
	}

	return listOfObject, nil
}

func CheckEngine(conn *sql.DB) ([]models.InformationTableEngine, error) {
	const statement string = `SELECT table_schema, table_name, engine 
	FROM information_schema.tables WHERE engine != 'InnoDB' AND 
	table_schema NOT IN ('mysql','performance_schema','sys','information_schema');`

	var warningTables []models.InformationTableEngine
	var warningList = []string{"MyISAM", "Memory", "Federated"}

	var warnSet = make(map[string]struct{})
	for _, w := range warningList {
		warnSet[w] = struct{}{}
	}

	// query
	rows, err := conn.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// loop through each result
	for rows.Next() {
		// scan query results
		var item models.InformationTableEngine
		err := rows.Scan(&item.SchemaName, &item.TableName, &item.Engine)
		if err != nil {
			return nil, err
		}

		if _, found := warnSet[item.Engine]; found {
			warningTables = append(warningTables, item)
		}
	}

	return warningTables, nil
}

func CheckFKDuplication(conn *sql.DB) (int, error) {
	const statement string = `SELECT COUNT(*) AS constraint_count FROM 
	INFORMATION_SCHEMA.TABLE_CONSTRAINTS WHERE CONSTRAINT_TYPE = 'FOREIGN KEY' 
	GROUP BY CONSTRAINT_SCHEMA, CONSTRAINT_NAME;`
	var FKDuplicationCounts int = 0
	// query
	rows, err := conn.Query(statement)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	// loop through each result
	for rows.Next() {
		// scan query results
		var FKCount int;
		err := rows.Scan(&FKCount)
		if err != nil {
			return 0, err
		}

		if FKCount > 1 {
			FKDuplicationCounts += 1
		}
	}

	return FKDuplicationCounts, nil
}