package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/3-lines-studio/datafrost/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresAdapter struct {
	conn *sql.DB
}

func NewPostgresAdapterRegistration() models.AdapterRegistration {
	return models.AdapterRegistration{
		Info: models.AdapterInfo{
			Type:        "postgres",
			Name:        "PostgreSQL",
			Description: "PostgreSQL database",
			UIConfig: models.UIConfig{
				Modes: []models.UIMode{
					{
						Key:   "url",
						Label: "Connection URL",
						Fields: []models.FieldConfig{
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
						Fields: []models.FieldConfig{
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
		Factory: func() models.DatabaseAdapter {
			return &PostgresAdapter{}
		},
	}
}

func (a *PostgresAdapter) Connect(credentials map[string]any) error {
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

func (a *PostgresAdapter) Close() error {
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}

func (a *PostgresAdapter) Ping() error {
	if a.conn == nil {
		return fmt.Errorf("not connected")
	}
	return a.conn.Ping()
}

func (a *PostgresAdapter) ListTables() ([]models.TableInfo, error) {
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

	var tables []models.TableInfo
	for rows.Next() {
		var t models.TableInfo
		if err := rows.Scan(&t.Name, &t.Type); err != nil {
			return nil, fmt.Errorf("failed to scan table: %w", err)
		}
		tables = append(tables, t)
	}

	return tables, rows.Err()
}

func (a *PostgresAdapter) GetTableData(tableName string, limit, offset int, filters []models.Filter) (*models.QueryResult, error) {
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

func (a *PostgresAdapter) ExecuteQuery(query string) (*models.QueryResult, error) {
	return a.executeQueryWithArgs(applyRowLimit(query), nil)
}


func (a *PostgresAdapter) executeQueryWithArgs(query string, args []any) (*models.QueryResult, error) {
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

	return &models.QueryResult{
		Columns: columns,
		Rows:    resultRows,
		Count:   len(resultRows),
		Total:   len(resultRows),
		Page:    1,
		Limit:   len(resultRows),
	}, nil
}

func (a *PostgresAdapter) getFilteredTableCount(tableName, whereClause string, args []any) (int, error) {
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

func buildPostgresWhereClause(filters []models.Filter) (string, []any) {
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
