package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
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


func SaveInformationObjectToCSV(serverConfigs []models.InformationObject, outputFile string) error {
	fmt.Println("Starting process for dumping data into csv file...")
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error to create output file: %s", err.Error())
	}
	defer file.Close()


	writer := csv.NewWriter(file)
	defer writer.Flush()

	// CSV Header
	writer.Write([]string{"OBJECT_TYPE", "SCHEMA_NAME", "OBJECT_NAME"})
	
	// Loop through each informationSchema, wrtie each records to csv file.
	for _, item := range serverConfigs {
		row := []string{
			item.ObjectType,
			item.SchemaName,
			item.ObjectName,
		}
		writer.Write(row)
	}

	fmt.Println("CSV file created successfully.")

	return nil
}


func CharSetReportToCSV(results []models.CharSetObject, outputDir string) error {
	const fileName string = "charset_err.csv"
	var outputFile string = filepath.Join(outputDir, fileName)
	
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error to create output file: %s", err.Error())
	}
	defer file.Close()


	writer := csv.NewWriter(file)
	defer writer.Flush()

	// CSV Header
	writer.Write([]string{"SEVERITY", "SCHEMA_NAME", "CharSet"})
	
	// Loop through each informationSchema, wrtie each records to csv file.
	for _, item := range results {
		row := []string{
			item.Severity,
			item.SchemaName,
			item.CharSet,
		}
		writer.Write(row)
	}

	return nil

}

func EngineReportToCSV(results []models.InformationTableEngine, outputDir string) error {
	const fileName string = "engine_err.csv"
	var outputFile string = filepath.Join(outputDir, fileName)
	
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error to create output file: %s", err.Error())
	}
	defer file.Close()


	writer := csv.NewWriter(file)
	defer writer.Flush()

	// CSV Header
	writer.Write([]string{"SCHEMA_NAME", "TABLE_NAME", "ENGINE", "CREATE_OPTIONS"})
	
	// Loop through each informationSchema, wrtie each records to csv file.
	for _, item := range results {
		row := []string{
			item.SchemaName,
			item.TableName,
			item.Engine,
			item.CreateOptions,
		}
		writer.Write(row)
	}

	return nil

}

func RowFormatReportToCSV(results []models.InformationRowFormat, outputDir string) error {
	const fileName string = "row_format_err.csv"
	var outputFile string = filepath.Join(outputDir, fileName)
	
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error to create output file: %s", err.Error())
	}
	defer file.Close()


	writer := csv.NewWriter(file)
	defer writer.Flush()

	// CSV Header
	writer.Write([]string{"SCHEMA_NAME", "TABLE_NAME", "ENGINE", "ROW_FORMAT"})
	
	// Loop through each informationSchema, wrtie each records to csv file.
	for _, item := range results {
		row := []string{
			item.SchemaName,
			item.TableName,
			item.Engine,
			item.RowFormat,
		}
		writer.Write(row)
	}

	return nil
}


func PKReportToCSV(results []models.InformationNoPKTable, outputDir string) error {
	const fileName string = "noPK_err.csv"
	var outputFile string = filepath.Join(outputDir, fileName)
	
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error to create output file: %s", err.Error())
	}
	defer file.Close()


	writer := csv.NewWriter(file)
	defer writer.Flush()

	// CSV Header
	writer.Write([]string{"SCHEMA_NAME", "TABLE_NAME"})
	
	// Loop through each informationSchema, wrtie each records to csv file.
	for _, item := range results {
		row := []string{
			item.SchemaName,
			item.TableName,
		}
		writer.Write(row)
	}

	return nil
}

func ViewReportToCSV(results []models.InformationView, outputDir string) error {
	const fileName string = "view_err.csv"
	var outputFile string = filepath.Join(outputDir, fileName)
	
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error to create output file: %s", err.Error())
	}
	defer file.Close()


	writer := csv.NewWriter(file)
	defer writer.Flush()

	// CSV Header
	writer.Write([]string{"SCHEMA_NAME", "TABLE_NAME", "VIEW_DEF"})
	
	// Loop through each informationSchema, wrtie each records to csv file.
	for _, item := range results {
		row := []string{
			item.SchemaName,
			item.TableName,
			item.ViewDefinition,
		}
		writer.Write(row)
	}

	return nil

}

func SyntaxRoutineToCSV(results []models.InformationRoutineDeprecated, outputDir string) error {
	const fileName string = "syntax_routine_err.csv"
	var outputFile string = filepath.Join(outputDir, fileName)
	
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error to create output file: %s", err.Error())
	}
	defer file.Close()


	writer := csv.NewWriter(file)
	defer writer.Flush()

	// CSV Header
	writer.Write([]string{"SCHEMA_NAME", "ROUTINE_NAME", "ROUTINE_TYPE", "ROUTINE_DEF"})
	
	// Loop through each informationSchema, wrtie each records to csv file.
	for _, item := range results {
		row := []string{
			item.SchemaName,
			item.RoutineName,
			item.RoutineType,
			item.RoutineDefinition,
		}
		writer.Write(row)
	}

	return nil
}

func FunctionRoutineToCSV(results []models.InformationRoutineDeprecated, outputDir string) error {
	const fileName string = "func_routine_err.csv"
	var outputFile string = filepath.Join(outputDir, fileName)
	
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error to create output file: %s", err.Error())
	}
	defer file.Close()


	writer := csv.NewWriter(file)
	defer writer.Flush()

	// CSV Header
	writer.Write([]string{"SCHEMA_NAME", "ROUTINE_NAME", "ROUTINE_TYPE", "ROUTINE_DEF"})
	
	// Loop through each informationSchema, wrtie each records to csv file.
	for _, item := range results {
		row := []string{
			item.SchemaName,
			item.RoutineName,
			item.RoutineType,
			item.RoutineDefinition,
		}
		writer.Write(row)
	}

	return nil
}

func SizeToCSV(results []models.InformationSizeSchema, outputDir string) error {
	const fileName string = "schema_size.csv"
	var outputFile string = filepath.Join(outputDir, fileName)
	
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error to create output file: %s", err.Error())
	}
	defer file.Close()


	writer := csv.NewWriter(file)
	defer writer.Flush()

	// CSV Header
	writer.Write([]string{"SCHEMA_NAME", "SIZE (MB)"})
	
	// Loop through each informationSchema, wrtie each records to csv file.
	for _, item := range results {
		row := []string{
			item.SchemaName,
			strconv.FormatFloat(item.Size, 'f', 2, 32),
		}
		writer.Write(row)
	}

	return nil
}