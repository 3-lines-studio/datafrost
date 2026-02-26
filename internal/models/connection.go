package models

import "time"

type DatabaseAdapter interface {
	Connect(credentials map[string]interface{}) error
	Close() error
	ListTables() ([]TableInfo, error)
	GetTableData(tableName string, limit, offset int, filters []Filter) (*QueryResult, error)
	ExecuteQuery(query string) (*QueryResult, error)
	Ping() error
}

type AdapterRegistration struct {
	Info    AdapterInfo
	Factory func() DatabaseAdapter
}

type Connection struct {
	ID          int64                  `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Credentials map[string]interface{} `json:"credentials"`
	CreatedAt   time.Time              `json:"created_at"`
}

type CreateConnectionRequest struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Credentials map[string]interface{} `json:"credentials"`
}

type UpdateConnectionRequest struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Credentials map[string]interface{} `json:"credentials"`
}

type TestConnectionRequest struct {
	Type        string                 `json:"type"`
	Credentials map[string]interface{} `json:"credentials"`
}

type TableInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type QueryResult struct {
	Columns []string `json:"columns"`
	Rows    [][]any  `json:"rows"`
	Count   int      `json:"count"`
	Total   int      `json:"total"`
	Page    int      `json:"page"`
	Limit   int      `json:"limit"`
}

type Filter struct {
	ID       string `json:"id"`
	Column   string `json:"column"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type AdapterInfo struct {
	Type        string   `json:"type"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	UIConfig    UIConfig `json:"ui_config"`
}

type UIConfig struct {
	Modes        []UIMode      `json:"modes,omitempty"`
	Fields       []FieldConfig `json:"fields,omitempty"`
	SupportsFile bool          `json:"supports_file"`
	FileTypes    []string      `json:"file_types,omitempty"`
}

type UIMode struct {
	Key    string        `json:"key"`
	Label  string        `json:"label"`
	Fields []FieldConfig `json:"fields"`
}

type FieldConfig struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Placeholder string `json:"placeholder,omitempty"`
}
