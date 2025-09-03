package features

import (
	"database/sql"
	"fmt"
	"github.com/sahatsawats/mysql-align/utils"
	"github.com/sahatsawats/mysql-align/models"
)
// Maybe some schema or table contain some special character which required to place in "" -> adjust again
func ReconcileRow(conn *sql.DB) ([]models.InformationSchema, error) {
	var tableInformations []models.InformationSchema
	var schemas []string
	var err error

	// get all schemas in database.
	schemas, err = getAllSchemas(conn)
	if err != nil {
		return nil, err
	}

	// logging total number of schema.
	numberOfSchemas := len(schemas)
	fmt.Printf("Found %d schemas in database.\n", numberOfSchemas)

	// prepare-statement to query all tables within given database.
	getTableStatement, err := conn.Prepare("SELECT table_schema, table_name FROM information_schema.tables WHERE table_schema = ?")
	if err != nil {
		return nil, err
	}
	defer getTableStatement.Close()

	// loop through each schema name.
	for _, schema := range schemas {
		// get all tables within given database.
		tables, err := getTablesInSchema(getTableStatement, schema)
		if err != nil {
			return nil, err
		}

		if len(tables) == 0 {
			fmt.Printf("schema: %s have zero table. skip process...", schema)
			continue
		}

		// loop query an sum of rows in given tables.
		for i := range tables {
			var rowCounts int
			// Get schema and table name
			var table models.InformationSchema = tables[i]
			// detect empty string.
			if table.SchemaName == "" || table.TableName == "" {
				fmt.Println("detect empty schema or table name. skip process...")
				continue
			}
			// Count row statement that received schema and table name
			countRowStatement := fmt.Sprintf("SELECT COUNT(*) FROM `%s`.`%s`", table.SchemaName, table.TableName)
			utils.Debug(fmt.Sprintf("Executing row count query: %s",countRowStatement))
			// Query count of row
			row := conn.QueryRow(countRowStatement)
			err := row.Scan(&rowCounts)
			if err != nil {
				return nil, err
			}

			// change field row in tableInformation.
			tables[i].Rows = rowCounts
			tableInformations = append(tableInformations, tables[i])
		}
	}

	return tableInformations, nil
}

func getAllSchemas(conn *sql.DB) ([]string, error) {
	// variables to hold schema name
	var listOfSchemas []string
	// statement to find all of schemas
	stmt := "SELECT schema_name FROM information_schema.schemata WHERE schema_name NOT IN ('mysql', 'performance_schema', 'information_schema', 'sys')"

	utils.Debug("Starting process: getAllSchemas...")

	// execute statement query
	rows, err := conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		// schema name from each row query
		var schemaName string
		// bind results to schemaName variable
		err := rows.Scan(&schemaName)
		if err != nil {
			return nil, err
		}
		utils.Debug(fmt.Sprintf("Schema found: %s", schemaName))
		// append results to list
		listOfSchemas = append(listOfSchemas, schemaName)
	}
	utils.Debug(fmt.Sprintf("Total schemas found: %d -> %v", len(listOfSchemas), listOfSchemas))
	return listOfSchemas, nil
}

func getTablesInSchema(prepareStmt *sql.Stmt, schemaName string) ([]models.InformationSchema, error) {
	var listOfInformationTable []models.InformationSchema

	utils.Debug(fmt.Sprintf("Starting table retrieval for schema: %s", schemaName))

	// query all of tables name in given schema
	tables, err := prepareStmt.Query(schemaName)
	if err != nil {
		return nil, err
	}

	for tables.Next() {
		// scan results and bind it to variables
		var tableSchema string
		var tableName string
		if err := tables.Scan(&tableSchema, &tableName); err != nil {
			return nil, err
		}

		// bind a query results to datatype.
		informationTable := models.InformationSchema{
			SchemaName: tableSchema,
			TableName:  tableName,
		}
		utils.Debug(fmt.Sprintf("Table discovered: %s.%s", tableSchema, tableName))

		// append datatype to list
		listOfInformationTable = append(listOfInformationTable, informationTable)
	}

	return listOfInformationTable, nil
}

func ReconcileObject(conn *sql.DB) ([]models.InformationObject, error) {
	const statement string = `SELECT 'ObjectType','DatabaseName','ObjectName'
	UNION ALL
	SELECT 'Table', TABLE_SCHEMA, TABLE_NAME 
	FROM information_schema.TABLES 
	WHERE TABLE_SCHEMA NOT IN ('mysql', 'performance_schema', 'information_schema', 'sys')

	UNION ALL
	SELECT 'View', TABLE_SCHEMA, TABLE_NAME 
	FROM information_schema.VIEWS 
	WHERE TABLE_SCHEMA NOT IN ('mysql', 'performance_schema', 'information_schema', 'sys')

	UNION ALL
	SELECT 'Trigger', TRIGGER_SCHEMA, TRIGGER_NAME 
	FROM information_schema.TRIGGERS 
	WHERE TRIGGER_SCHEMA NOT IN ('mysql', 'performance_schema', 'information_schema', 'sys')

	UNION ALL
	SELECT 'Procedure', ROUTINE_SCHEMA, ROUTINE_NAME 
	FROM information_schema.ROUTINES 
	WHERE ROUTINE_TYPE='PROCEDURE' 
	AND ROUTINE_SCHEMA NOT IN ('mysql', 'performance_schema', 'information_schema', 'sys')

	UNION ALL
	SELECT 'Function', ROUTINE_SCHEMA, ROUTINE_NAME 
	FROM information_schema.ROUTINES 
	WHERE ROUTINE_TYPE='FUNCTION' 
	AND ROUTINE_SCHEMA NOT IN ('mysql', 'performance_schema', 'information_schema', 'sys')

	UNION ALL
	SELECT 'Index', TABLE_SCHEMA, CONCAT(TABLE_NAME, ' (', INDEX_NAME, ')')
	FROM information_schema.STATISTICS
	WHERE TABLE_SCHEMA NOT IN ('mysql', 'performance_schema', 'information_schema', 'sys')

	ORDER BY DatabaseName, ObjectType, ObjectName;`

	var informationObjects []models.InformationObject
	// query statement
	rows, err := conn.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// loop through each result
	for rows.Next() {
		var object models.InformationObject

		err := rows.Scan(&object.ObjectType, &object.SchemaName, &object.ObjectName)
		if err != nil {
			return nil, err
		}

		utils.Debug(fmt.Sprintf("Object discovered: {object_type: %s, schema_name: %s, object_name: %s", object.ObjectType, object.SchemaName, object.ObjectName))

		informationObjects = append(informationObjects, object)
	}

	return informationObjects, nil
}