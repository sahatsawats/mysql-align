package models

import "time"

type PreMigrationReport struct {
	Meta     ReportMeta
	Sections []CheckSection
}

type ReportMeta struct {
	ToolVersion          string
	GeneratedAt          time.Time
	GeneratedAtFormatted string
	Host                 string
	Port                 int
	User                 string
	ServerVersion        string
}

type CheckSection struct {
	ID          string
	Name        string
	Description string
	DefaultSev  string
	Status      string
	Error       string
	Count       int
	Columns     []string
	Rows        [][]string
	CodeCols    []bool
	Truncated   bool
	TotalCount  int
	Note        string
}
