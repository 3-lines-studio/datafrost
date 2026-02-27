package database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/3-lines-studio/datafrost/internal/core/entity"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type postgresAdapter struct {
	conn *sql.DB
}

func newPostgresAdapterRegistration() entity.AdapterRegistration {
	return entity.AdapterRegistration{
		Info: entity.AdapterInfo{
			Type:        "postgres",
			Name:        "PostgreSQL",
			Description: "PostgreSQL database",
			UIConfig: entity.UIConfig{
				Modes: []entity.UIMode{
					{
						Key:   "url",
						Label: "Connection URL",
						Fields: []entity.FieldConfig{
							{
								Key:         "url",
								Label:       "Connection URL",
								Type:        "password",
								Required:    true,
								Placeholder: "postgres://user:password@localhost:5432/database?sslmode=prefer",
							},
						},
					},
					{
						Key:   "fields",
						Label: "Individual Fields",
						Fields: []entity.FieldConfig{
							{
								Key:         "host",
								Label:       "Host",
								Type:        "text",
								Required:    true,
								Placeholder: "localhost",
							},
							{
								Key:         "port",
								Label:       "Port",
								Type:        "text",
								Required:    true,
								Placeholder: "5432",
							},
							{
								Key:         "database",
								Label:       "Database",
								Type:        "text",
								Required:    true,
								Placeholder: "postgres",
							},
							{
								Key:         "username",
								Label:       "Username",
								Type:        "text",
								Required:    true,
								Placeholder: "postgres",
							},
							{
								Key:         "password",
								Label:       "Password",
								Type:        "password",
								Required:    false,
								Placeholder: "",
							},
							{
								Key:         "ssl_mode",
								Label:       "SSL Mode",
								Type:        "text",
								Required:    true,
								Placeholder: "prefer",
							},
						},
					},
				},
				SupportsFile: false,
			},
		},
		Factory: func() entity.DatabaseAdapter {
			return &postgresAdapter{}
		},
	}
}

func (a *postgresAdapter) Connect(credentials map[string]any) error {
	mode, _ := credentials["mode"].(string)

	var connStr string
	if mode == "url" {
		url, ok := credentials["url"].(string)
		if !ok || url == "" {
			return fmt.Errorf("url is required")
		}
		connStr = url
	} else {
		host, ok := credentials["host"].(string)
		if !ok || host == "" {
			return fmt.Errorf("host is required")
		}

		port := "5432"
		if p, ok := credentials["port"].(string); ok && p != "" {
			port = p
		} else if p, ok := credentials["port"].(float64); ok {
			port = strconv.Itoa(int(p))
		}

		database, ok := credentials["database"].(string)
		if !ok || database == "" {
			return fmt.Errorf("database is required")
		}

		username, ok := credentials["username"].(string)
		if !ok || username == "" {
			return fmt.Errorf("username is required")
		}

		password, _ := credentials["password"].(string)
		sslMode, _ := credentials["ssl_mode"].(string)
		if sslMode == "" {
			sslMode = "prefer"
		}

		connStr = fmt.Sprintf("host=%s port=%s dbname=%s user=%s sslmode=%s",
			host, port, database, username, sslMode)
		if password != "" {
			connStr += fmt.Sprintf(" password=%s", password)
		}
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return fmt.Errorf("failed to open postgres connection: %w", err)
	}

	a.conn = db
	return nil
}

func (a *postgresAdapter) Close() error {
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}

func (a *postgresAdapter) Ping() error {
	if a.conn == nil {
		return fmt.Errorf("not connected")
	}
	return a.conn.Ping()
}

func (a *postgresAdapter) ListTables() ([]entity.TableInfo, error) {
	rows, err := a.conn.Query(`
		SELECT table_name, 'table' as type
		FROM information_schema.tables
		WHERE table_schema = 'public'
		ORDER BY table_name
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list tables: %w", err)
	}
	defer rows.Close()

	var tables []entity.TableInfo
	for rows.Next() {
		var t entity.TableInfo
		if err := rows.Scan(&t.Name, &t.Type); err != nil {
			return nil, fmt.Errorf("failed to scan table: %w", err)
		}
		tables = append(tables, t)
	}

	return tables, rows.Err()
}

func (a *postgresAdapter) GetTableData(tableName string, limit, offset int, filters []entity.Filter) (*entity.QueryResult, error) {
	whereClause, args := buildPostgresWhereClause(filters)

	count, err := a.getFilteredTableCount(tableName, whereClause, args)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("SELECT * FROM \"%s\"", tableName)
	if whereClause != "" {
		query += " WHERE " + whereClause
	}
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	result, err := a.executeQueryWithArgs(query, args)
	if err != nil {
		return nil, err
	}

	result.Total = count
	result.Page = offset/limit + 1
	result.Limit = limit
	return result, nil
}

func (a *postgresAdapter) ExecuteQuery(query string) (*entity.QueryResult, error) {
	return a.executeQueryWithArgs(query, nil)
}

func (a *postgresAdapter) executeQueryWithArgs(query string, args []any) (*entity.QueryResult, error) {
	upperQuery := strings.ToUpper(strings.TrimSpace(query))
	isSelect := strings.HasPrefix(upperQuery, "SELECT") || strings.HasPrefix(upperQuery, "WITH")

	if !isSelect {
		return nil, fmt.Errorf("only SELECT and WITH queries are allowed")
	}

	ctx := context.Background()
	rows, err := a.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var resultRows [][]any
	for rows.Next() {
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		resultRows = append(resultRows, values)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return &entity.QueryResult{
		Columns: columns,
		Rows:    resultRows,
		Count:   len(resultRows),
		Total:   len(resultRows),
		Page:    1,
		Limit:   len(resultRows),
	}, nil
}

func (a *postgresAdapter) getFilteredTableCount(tableName, whereClause string, args []any) (int, error) {
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM \"%s\"", tableName)
	if whereClause != "" {
		countQuery += " WHERE " + whereClause
	}
	var count int
	err := a.conn.QueryRow(countQuery, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count rows: %w", err)
	}
	return count, nil
}

func (a *postgresAdapter) GetTableSchema(tableName string) (*entity.TableSchema, error) {
	schema := &entity.TableSchema{
		TableName: tableName,
	}

	columnRows, err := a.conn.Query(`
		SELECT 
			column_name,
			data_type,
			is_nullable = 'YES',
			COALESCE(column_default, '') as column_default,
			false as is_primary_key
		FROM information_schema.columns
		WHERE table_name = $1 AND table_schema = 'public'
		ORDER BY ordinal_position
	`, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}
	defer columnRows.Close()

	var columns []entity.ColumnInfo
	for columnRows.Next() {
		var col entity.ColumnInfo
		if err := columnRows.Scan(&col.Name, &col.Type, &col.Nullable, &col.DefaultValue, &col.IsPrimaryKey); err != nil {
			return nil, fmt.Errorf("failed to scan column: %w", err)
		}
		columns = append(columns, col)
	}

	pkRows, err := a.conn.Query(`
		SELECT kcu.column_name
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu 
			ON tc.constraint_name = kcu.constraint_name
		WHERE tc.table_name = $1 
			AND tc.table_schema = 'public'
			AND tc.constraint_type = 'PRIMARY KEY'
	`, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get primary keys: %w", err)
	}
	defer pkRows.Close()

	pkColumns := make(map[string]bool)
	for pkRows.Next() {
		var colName string
		if err := pkRows.Scan(&colName); err != nil {
			return nil, fmt.Errorf("failed to scan primary key: %w", err)
		}
		pkColumns[colName] = true
	}

	for i := range columns {
		if pkColumns[columns[i].Name] {
			columns[i].IsPrimaryKey = true
		}
	}
	schema.Columns = columns

	indexRows, err := a.conn.Query(`
		SELECT 
			i.relname as index_name,
			ix.indisunique as is_unique
		FROM pg_index ix
		JOIN pg_class i ON i.oid = ix.indexrelid
		JOIN pg_class t ON t.oid = ix.indrelid
		WHERE t.relname = $1
		GROUP BY i.relname, ix.indisunique
	`, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get indexes: %w", err)
	}
	defer indexRows.Close()

	var indexes []entity.IndexInfo
	for indexRows.Next() {
		var idx entity.IndexInfo
		if err := indexRows.Scan(&idx.Name, &idx.Unique); err != nil {
			return nil, fmt.Errorf("failed to scan index: %w", err)
		}

		indexColRows, err := a.conn.Query(`
			SELECT a.attname
			FROM pg_index ix
			JOIN pg_class i ON i.oid = ix.indexrelid
			JOIN pg_class t ON t.oid = ix.indrelid
			JOIN pg_attribute a ON a.attrelid = t.oid AND a.attnum = ANY(ix.indkey)
			WHERE i.relname = $1
			ORDER BY array_position(ix.indkey, a.attnum)
		`, idx.Name)
		if err != nil {
			continue
		}

		var cols []string
		for indexColRows.Next() {
			var colName string
			if err := indexColRows.Scan(&colName); err == nil {
				cols = append(cols, colName)
			}
		}
		indexColRows.Close()

		idx.Columns = cols
		indexes = append(indexes, idx)
	}
	schema.Indexes = indexes

	constraintRows, err := a.conn.Query(`
		SELECT 
			con.conname as constraint_name,
			con.contype::text as constraint_type,
			pg_get_constraintdef(con.oid) as definition
		FROM pg_constraint con
		JOIN pg_class t ON t.oid = con.conrelid
		WHERE t.relname = $1 
			AND con.contype IN ('f', 'u', 'c')
	`, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get constraints: %w", err)
	}
	defer constraintRows.Close()

	var constraints []entity.ConstraintInfo
	for constraintRows.Next() {
		var c entity.ConstraintInfo
		var typeStr string
		if err := constraintRows.Scan(&c.Name, &typeStr, &c.Definition); err != nil {
			continue
		}
		switch typeStr {
		case "f":
			c.Type = "FOREIGN KEY"
		case "u":
			c.Type = "UNIQUE"
		case "c":
			c.Type = "CHECK"
		}
		c.Column = ""
		constraints = append(constraints, c)
	}
	schema.Constraints = constraints

	return schema, nil
}

func buildPostgresWhereClause(filters []entity.Filter) (string, []any) {
	if len(filters) == 0 {
		return "", nil
	}

	var conditions []string
	var args []any
	argIndex := 1

	for _, filter := range filters {
		if filter.Column == "" {
			continue
		}

		column := fmt.Sprintf("\"%s\"", filter.Column)

		switch filter.Operator {
		case "eq":
			conditions = append(conditions, fmt.Sprintf("%s = $%d", column, argIndex))
			args = append(args, filter.Value)
			argIndex++
		case "neq":
			conditions = append(conditions, fmt.Sprintf("%s != $%d", column, argIndex))
			args = append(args, filter.Value)
			argIndex++
		case "gt":
			conditions = append(conditions, fmt.Sprintf("%s > $%d", column, argIndex))
			args = append(args, filter.Value)
			argIndex++
		case "lt":
			conditions = append(conditions, fmt.Sprintf("%s < $%d", column, argIndex))
			args = append(args, filter.Value)
			argIndex++
		case "gte":
			conditions = append(conditions, fmt.Sprintf("%s >= $%d", column, argIndex))
			args = append(args, filter.Value)
			argIndex++
		case "lte":
			conditions = append(conditions, fmt.Sprintf("%s <= $%d", column, argIndex))
			args = append(args, filter.Value)
			argIndex++
		case "like":
			conditions = append(conditions, fmt.Sprintf("%s LIKE $%d", column, argIndex))
			args = append(args, filter.Value)
			argIndex++
		case "not_like":
			conditions = append(conditions, fmt.Sprintf("%s NOT LIKE $%d", column, argIndex))
			args = append(args, filter.Value)
			argIndex++
		case "is_null":
			conditions = append(conditions, fmt.Sprintf("%s IS NULL", column))
		case "is_not_null":
			conditions = append(conditions, fmt.Sprintf("%s IS NOT NULL", column))
		}
	}

	if len(conditions) == 0 {
		return "", nil
	}

	return strings.Join(conditions, " AND "), args
}
