package entity

type DatabaseAdapter interface {
	Connect(credentials map[string]any) error
	Close() error
	ListTables() ([]TableInfo, error)
	GetTableData(tableName string, limit, offset int, filters []Filter) (*QueryResult, error)
	ExecuteQuery(query string) (*QueryResult, error)
	Ping() error
	GetTableSchema(tableName string) (*TableSchema, error)
}

type AdapterRegistration struct {
	Info    AdapterInfo
	Factory func() DatabaseAdapter
}
