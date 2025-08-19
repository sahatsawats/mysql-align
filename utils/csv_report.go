package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/sahatsawats/mysql-align/models"
)


func SaveInformationTablesToCSV(informationTable []models.InformationSchema, outputFile string) error {
	fmt.Println("Starting process for dumping data into csv file...")
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error to create output file: %s", err.Error())
	}
	defer file.Close()


	writer := csv.NewWriter(file)
	defer writer.Flush()

	// CSV Header
	writer.Write([]string{"Schema", "Table", "Row"})
	
	// Loop through each informationSchema, wrtie each records to csv file.
	for _, item := range informationTable {
		row := []string{
			item.SchemaName,
			item.TableName,
			strconv.Itoa(item.Rows),
		}
		writer.Write(row)
	}

	fmt.Println("CSV file created successfully.")

	return nil
}

func SaveServerConfigurationToCSV(serverConfigs []models.InformationConfig, outputFile string) error {
	fmt.Println("Starting process for dumping data into csv file...")
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error to create output file: %s", err.Error())
	}
	defer file.Close()


	writer := csv.NewWriter(file)
	defer writer.Flush()

	// CSV Header
	writer.Write([]string{"SERVER_VARIABLES", "SERVER_VAULE"})
	
	// Loop through each informationSchema, wrtie each records to csv file.
	for _, item := range serverConfigs {
		row := []string{
			item.VariableName,
			item.VariableVaule,
		}
		writer.Write(row)
	}

	fmt.Println("CSV file created successfully.")

	return nil
}