package models


type InformationSchema struct {
	SchemaName	string;
	TableName	string;
	Rows 		int;
}

type InformationConfig struct {
	VariableName	string;
	VariableVaule	string;
}

type InformationObject struct {
	ObjectType	string;
	SchemaName	string;
	ObjectName	string;
}

type CharSetObject struct {
	Severity	string;
	SchemaName	string;
	CharSet		string;
}

type InformationTableEngine struct {
	SchemaName 	string;
	TableName	string;
	Engine		string;
	CreateOptions string;
}

type InformationRowFormat struct {
	SchemaName 	string;
	TableName	string;
	Engine		string;
	RowFormat string;
}

type InformationView struct {
	SchemaName 		string;
	TableName		string;
	ViewDefinition	string;
}

type InformationRoutineDeprecated struct {
	SchemaName 		string;
	RoutineName		string;
	RoutineType		string;
	RoutineDefinition	string;
}

type InformationNoPKTable struct {
	SchemaName	string
	TableName	string
}