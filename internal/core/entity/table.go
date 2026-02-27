package entity

type TableInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type TableSchema struct {
	TableName   string           `json:"table_name"`
	Columns     []ColumnInfo     `json:"columns"`
	Indexes     []IndexInfo      `json:"indexes"`
	Constraints []ConstraintInfo `json:"constraints"`
}

type ColumnInfo struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Nullable     bool   `json:"nullable"`
	DefaultValue string `json:"default_value"`
	IsPrimaryKey bool   `json:"is_primary_key"`
}

type IndexInfo struct {
	Name    string   `json:"name"`
	Unique  bool     `json:"unique"`
	Columns []string `json:"columns"`
}

type ConstraintInfo struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Column     string `json:"column"`
	Definition string `json:"definition"`
}
