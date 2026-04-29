package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/sahatsawats/mysql-align/models"
)

func csvFilePath(outputDir, prefix string) (string, error) {
	ts := time.Now().Format("2006-01-02_15-04-05")
	fileName := fmt.Sprintf("%s_%s.csv", prefix, ts)
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return "", fmt.Errorf("create output directory: %w", err)
	}
	return filepath.Join(outputDir, fileName), nil
}

func SaveInformationTablesToCSV(informationTable []models.InformationSchema, outputDir string) error {
	outputFile, err := csvFilePath(outputDir, "recon_rows")
	if err != nil {
		return err
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error to create output file: %s", err.Error())
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Schema", "Table", "Row"})
	for _, item := range informationTable {
		writer.Write([]string{
			item.SchemaName,
			item.TableName,
			strconv.Itoa(item.Rows),
		})
	}

	fmt.Println("CSV file created:", outputFile)
	return nil
}

func SaveServerConfigurationToCSV(serverConfigs []models.InformationConfig, outputDir string) error {
	outputFile, err := csvFilePath(outputDir, "server_config")
	if err != nil {
		return err
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error to create output file: %s", err.Error())
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"SERVER_VARIABLES", "SERVER_VALUE"})
	for _, item := range serverConfigs {
		writer.Write([]string{
			item.VariableName,
			item.VariableValue,
		})
	}

	fmt.Println("CSV file created:", outputFile)
	return nil
}

func SaveInformationObjectToCSV(objects []models.InformationObject, outputDir string) error {
	outputFile, err := csvFilePath(outputDir, "recon_objects")
	if err != nil {
		return err
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error to create output file: %s", err.Error())
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"OBJECT_TYPE", "SCHEMA_NAME", "OBJECT_NAME"})
	for _, item := range objects {
		writer.Write([]string{
			item.ObjectType,
			item.SchemaName,
			item.ObjectName,
		})
	}

	fmt.Println("CSV file created:", outputFile)
	return nil
}

func SizeToCSV(results []models.InformationSizeSchema, outputDir string) error {
	outputFile, err := csvFilePath(outputDir, "schema_size")
	if err != nil {
		return err
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error to create output file: %s", err.Error())
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"SCHEMA_NAME", "SIZE (MB)"})
	for _, item := range results {
		writer.Write([]string{
			item.SchemaName,
			strconv.FormatFloat(item.Size, 'f', 2, 32),
		})
	}

	fmt.Println("CSV file created:", outputFile)
	return nil
}
