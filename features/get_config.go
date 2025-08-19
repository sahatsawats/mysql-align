package features

import (
	"database/sql"
	"github.com/sahatsawats/mysql-align/models"
)

func GetConfiguration(conn *sql.DB) ([]models.InformationConfig, error) {
	// return 2 feilds, variable_name and variable_value.
	const statement string = "SELECT VARIABLE_NAME, VARIABLE_VALUE FROM performance_schema.global_variables"
	var serverConfigs []models.InformationConfig
	// query statement
	rows, err := conn.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// loop through each result
	for rows.Next() {
		var config models.InformationConfig

		err := rows.Scan(&config.VariableName, &config.VariableVaule)
		if err != nil {
			return nil, err
		}

		serverConfigs = append(serverConfigs, config)
	}

	return serverConfigs, nil
}
