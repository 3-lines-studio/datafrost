package adapters

import (
	"context"
	"database/sql"
	"datafrost/internal/models"
	"fmt"
	"strings"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type TursoAdapter struct {
	conn *sql.DB
}

func NewTursoAdapterRegistration() models.AdapterRegistration {
	return models.AdapterRegistration{
		Info: models.AdapterInfo{
			Type:        "turso",
			Name:        "Turso",
			Description: "Turso (libSQL) SQLite database",
			UIConfig: models.UIConfig{
				Fields: []models.FieldConfig{
					{
						Key:         "url",
						Label:       "Database URL",
						Type:        "text",
						Required:    true,
						Placeholder: "libsql://...",
					},
					{
						Key:         "token",
						Label:       "Auth Token",
						Type:        "password",
						Required:    true,
						Placeholder: "your-auth-token",
					},
				},
				SupportsFile: false,
			},
		},
		Factory: func() models.DatabaseAdapter {
			return &TursoAdapter{}
		},
	}
}

func (a *TursoAdapter) Connect(credentials map[string]interface{}) error {
	url, ok := credentials["url"].(string)
	if !ok || url == "" {
		return fmt.Errorf("url is required")
	}

	token, ok := credentials["token"].(string)
	if !ok {
		token = ""
	}

	var dbURL string
	if strings.Contains(url, "?") {
		dbURL = fmt.Sprintf("%s&authToken=%s", url, token)
	} else {
		dbURL = fmt.Sprintf("%s?authToken=%s", url, token)
	}

	database, err := sql.Open("libsql", dbURL)
	if err != nil {
		return fmt.Errorf("failed to open turso connection: %w", err)
	}

	a.conn = database
	return nil
}

func (a *TursoAdapter) Close() error {
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}

func (a *TursoAdapter) Ping() error {
	if a.conn == nil {
		return fmt.Errorf("not connected")
	}
	return a.conn.Ping()
}

func (a *TursoAdapter) ListTables() ([]models.TableInfo, error) {
	rows, err := a.conn.Query(
		"SELECT name, type FROM sqlite_master WHERE type IN ('table', 'view') AND name NOT LIKE 'sqlite_%' ORDER BY name",
	)
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

func (a *TursoAdapter) GetTableData(tableName string, limit, offset int, filters []models.Filter) (*models.QueryResult, error) {
	whereClause, args := buildTursoWhereClause(filters)

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

func (a *TursoAdapter) ExecuteQuery(query string) (*models.QueryResult, error) {
	return a.executeQueryWithArgs(query, nil)
}

func (a *TursoAdapter) executeQueryWithArgs(query string, args []interface{}) (*models.QueryResult, error) {
	upperQuery := strings.ToUpper(strings.TrimSpace(query))
	isSelect := strings.HasPrefix(upperQuery, "SELECT") || strings.HasPrefix(upperQuery, "WITH") || strings.HasPrefix(upperQuery, "PRAGMA")

	if !isSelect {
		return nil, fmt.Errorf("only SELECT, WITH, and PRAGMA queries are allowed")
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

func (a *TursoAdapter) getFilteredTableCount(tableName, whereClause string, args []interface{}) (int, error) {
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

func buildTursoWhereClause(filters []models.Filter) (string, []interface{}) {
	if len(filters) == 0 {
		return "", nil
	}

	var conditions []string
	var args []interface{}

	for _, filter := range filters {
		if filter.Column == "" {
			continue
		}

		column := fmt.Sprintf("\"%s\"", filter.Column)

		switch filter.Operator {
		case "eq":
			conditions = append(conditions, fmt.Sprintf("%s = ?", column))
			args = append(args, filter.Value)
		case "neq":
			conditions = append(conditions, fmt.Sprintf("%s != ?", column))
			args = append(args, filter.Value)
		case "gt":
			conditions = append(conditions, fmt.Sprintf("%s > ?", column))
			args = append(args, filter.Value)
		case "lt":
			conditions = append(conditions, fmt.Sprintf("%s < ?", column))
			args = append(args, filter.Value)
		case "gte":
			conditions = append(conditions, fmt.Sprintf("%s >= ?", column))
			args = append(args, filter.Value)
		case "lte":
			conditions = append(conditions, fmt.Sprintf("%s <= ?", column))
			args = append(args, filter.Value)
		case "like":
			conditions = append(conditions, fmt.Sprintf("%s LIKE ?", column))
			args = append(args, filter.Value)
		case "not_like":
			conditions = append(conditions, fmt.Sprintf("%s NOT LIKE ?", column))
			args = append(args, filter.Value)
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
