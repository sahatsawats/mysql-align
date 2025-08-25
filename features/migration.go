package features

import (
	"database/sql"
	"strings"

	"github.com/sahatsawats/mysql-align/models"
)

func CheckNoPK(conn *sql.DB) ([]models.InformationNoPKTable, error) {
	const statement string = `SELECT T.TABLE_SCHEMA, T.TABLE_NAME FROM information_schema.TABLES AS T WHERE
    T.TABLE_SCHEMA NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys')
    AND T.TABLE_TYPE = 'BASE TABLE' AND (T.TABLE_SCHEMA, T.TABLE_NAME) NOT IN (
    SELECT TABLE_SCHEMA, TABLE_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE CONSTRAINT_NAME = 'PRIMARY'
    );`
	var listOfNoPKTable []models.InformationNoPKTable
	// query
	rows, err := conn.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// loop through each result
	for rows.Next() {
		// scan query results
		var noPKTable models.InformationNoPKTable
		err := rows.Scan(&noPKTable.SchemaName, &noPKTable.TableName)
		if err != nil {
			return nil, err
		}

		listOfNoPKTable = append(listOfNoPKTable, noPKTable)
	}

	return listOfNoPKTable, nil
}

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
	const statement string = `SELECT table_schema, table_name, engine, CREATE_OPTIONS
	FROM information_schema.tables WHERE engine != 'InnoDB' AND 
	table_schema NOT IN ('mysql','performance_schema','sys','information_schema');`

	var warningTables []models.InformationTableEngine
	var warningList = []string{"MyISAM", "Memory", "FEDERATED"}

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
		err := rows.Scan(&item.SchemaName, &item.TableName, &item.Engine, &item.CreateOptions)
		if err != nil {
			return nil, err
		}

		if _, found := warnSet[item.Engine]; found {
			warningTables = append(warningTables, item)
		}
	}

	return warningTables, nil
}

func CheckRowFormat(conn *sql.DB) ([]models.InformationRowFormat, error) {
	const statement string = `SELECT TABLE_SCHEMA, TABLE_NAME, ENGINE, ROW_FORMAT  
	From information_schema.tables 
	WHERE table_type = 'BASE TABLE' AND table_schema NOT IN ('mysql','perform
	ance_schema','performance_schema', 'information_schema','sys','information_schema');`

	var warningRows []models.InformationRowFormat
	var warningList = []string{"Redundant", "Compact", "Fixed"}

	var warnSet = make(map[string]struct{})
	for _,w := range warningList {
		var lowerCase string = strings.ToLower(w) 
		warnSet[lowerCase] = struct{}{}
	}

	rows, err := conn.Query(statement)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// loop through each result
	for rows.Next() {
		// scan query results
		var item models.InformationRowFormat
		err := rows.Scan(&item.SchemaName, &item.TableName, &item.Engine, &item.RowFormat)
		if err != nil {
			return nil, err
		}
		// Use map to find a key which convert to lowercase
		if _, found := warnSet[strings.ToLower(item.RowFormat)]; found {
			warningRows = append(warningRows, item)
		}
	}

	return warningRows, nil
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
		var FKCount int
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

func CheckViewDeprecated(conn *sql.DB) ([]models.InformationView, error) {
	const statement string = `SELECT table_schema, table_name, view_definition FROM information_schema.views 
	WHERE table_schema NOT IN ('mysql', 'performance_schema', 'sys', 'information_schema') AND 
	(view_definition LIKE '%GROUP BY%ASC%' OR view_definition LIKE '%GROUP BY%DESC%');`
	var listOfDeprecateViews []models.InformationView

	rows, err := conn.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// loop through each result
	for rows.Next() {
		// scan query results
		var deprecateView models.InformationView
		err := rows.Scan(&deprecateView.SchemaName, &deprecateView.TableName, &deprecateView.ViewDefinition)
		if err != nil {
			return nil, err
		}

		listOfDeprecateViews = append(listOfDeprecateViews, deprecateView)
	}

	return listOfDeprecateViews, nil
}

func CheckRoutineSyntaxDeprecated(conn *sql.DB) ([]models.InformationRoutineDeprecated, error) {
	const statement string = `SELECT ROUTINE_SCHEMA, ROUTINE_NAME, ROUTINE_TYPE, ROUTINE_DEFINITION FROM
    information_schema.ROUTINES WHERE ROUTINE_SCHEMA NOT IN ('mysql', 'performance_schema', 'sys', 'information_schema')
    AND (ROUTINE_DEFINITION LIKE '%GROUP BY%ASC%' OR ROUTINE_DEFINITION LIKE '%GROUP BY%DESC%');`
	var listOfRoutineSyntaxDeprecated []models.InformationRoutineDeprecated

	rows, err := conn.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// loop through each result
	for rows.Next() {
		// scan query results
		var deprecateSyntax models.InformationRoutineDeprecated
		err := rows.Scan(&deprecateSyntax.SchemaName, &deprecateSyntax.RoutineName, &deprecateSyntax.RoutineType, &deprecateSyntax.RoutineDefinition)
		if err != nil {
			return nil, err
		}

		listOfRoutineSyntaxDeprecated = append(listOfRoutineSyntaxDeprecated, deprecateSyntax)
	}

	return listOfRoutineSyntaxDeprecated, nil
}

func CheckRoutineFunctionDeprecated(conn *sql.DB) ([]models.InformationRoutineDeprecated, error) {
	const statement string = `SELECT ROUTINE_SCHEMA, ROUTINE_NAME, ROUTINE_TYPE, ROUTINE_DEFINITION FROM
    information_schema.ROUTINES WHERE ROUTINE_SCHEMA NOT IN ('mysql', 'performance_schema', 'sys', 'information_schema')
    AND (ROUTINE_DEFINITION LIKE '%DECODE(%' OR ROUTINE_DEFINITION LIKE '%ENCODE(%' OR ROUTINE_DEFINITION LIKE '%COMPRESS(%')`
	var listOfRoutineFunctionDeprecated []models.InformationRoutineDeprecated

	rows, err := conn.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// loop through each result
	for rows.Next() {
		// scan query results
		var deprecateFunction models.InformationRoutineDeprecated
		err := rows.Scan(&deprecateFunction.SchemaName, &deprecateFunction.RoutineName, &deprecateFunction.RoutineType, &deprecateFunction.RoutineDefinition)
		if err != nil {
			return nil, err
		}

		listOfRoutineFunctionDeprecated = append(listOfRoutineFunctionDeprecated, deprecateFunction)
	}

	return listOfRoutineFunctionDeprecated, nil
}
