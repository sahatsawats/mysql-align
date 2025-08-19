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