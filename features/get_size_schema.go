package features

import (
	"database/sql"
	"fmt"
	"github.com/sahatsawats/mysql-align/utils"
	"github.com/sahatsawats/mysql-align/models"
)

func GetSchemaSize(conn *sql.DB) ([]models.InformationSizeSchema, error) {
	const statement string = `SELECT 
		table_schema AS database_name,
		ROUND(SUM(data_length + index_length) / 1024 / 1024, 2) AS size_mb
	FROM 
		information_schema.tables
	GROUP BY 
		table_schema
	ORDER BY 
		size_mb DESC;`

	var informationObjects []models.InformationSizeSchema
	// query statement
	rows, err := conn.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// loop through each result
	for rows.Next() {
		var object models.InformationSizeSchema

		err := rows.Scan(&object.SchemaName, &object.Size)
		if err != nil {
			return nil, err
		}

		utils.Debug(fmt.Sprintf("Object discovered: {schema_name: %s, size: %s", object.SchemaName, object.Size))

		informationObjects = append(informationObjects, object)
	}

	return informationObjects, nil
}