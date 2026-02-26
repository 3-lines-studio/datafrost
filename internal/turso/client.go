package turso

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/3-lines-studio/datafrost/internal/models"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type Client struct {
	conn *sql.DB
}

func NewClient(url, token string) (*Client, error) {
	var dbURL string
	if strings.Contains(url, "?") {
		dbURL = fmt.Sprintf("%s&authToken=%s", url, token)
	} else {
		dbURL = fmt.Sprintf("%s?authToken=%s", url, token)
	}

	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open turso connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping turso: %w", err)
	}

	return &Client{conn: db}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) ListTables() ([]models.TableInfo, error) {
	rows, err := c.conn.Query(
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

func (c *Client) GetTableCount(tableName string) (int, error) {
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM \"%s\"", tableName)
	var count int
	err := c.conn.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count rows: %w", err)
	}
	return count, nil
}

func (c *Client) getFilteredTableCount(tableName, whereClause string, args []any) (int, error) {
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM \"%s\"", tableName)
	if whereClause != "" {
		countQuery += " WHERE " + whereClause
	}
	var count int
	err := c.conn.QueryRow(countQuery, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count rows: %w", err)
	}
	return count, nil
}

func buildWhereClause(filters []models.Filter) (string, []any) {
	if len(filters) == 0 {
		return "", nil
	}

	var conditions []string
	var args []any

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

func (c *Client) executeQueryWithArgs(query string, args []any) (*models.QueryResult, error) {
	upperQuery := strings.ToUpper(strings.TrimSpace(query))
	isSelect := strings.HasPrefix(upperQuery, "SELECT") || strings.HasPrefix(upperQuery, "WITH") || strings.HasPrefix(upperQuery, "PRAGMA")

	if !isSelect {
		return nil, fmt.Errorf("only SELECT, WITH, and PRAGMA queries are allowed")
	}

	ctx := context.Background()
	rows, err := c.conn.QueryContext(ctx, query, args...)
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

func (c *Client) GetTableData(tableName string, limit, offset int, filters []models.Filter) (*models.QueryResult, error) {
	whereClause, args := buildWhereClause(filters)

	count, err := c.getFilteredTableCount(tableName, whereClause, args)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("SELECT * FROM \"%s\"", tableName)
	if whereClause != "" {
		query += " WHERE " + whereClause
	}
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	result, err := c.executeQueryWithArgs(query, args)
	if err != nil {
		return nil, err
	}

	result.Total = count
	result.Page = offset/limit + 1
	result.Limit = limit
	return result, nil
}

func (c *Client) ExecuteQuery(query string) (*models.QueryResult, error) {
	upperQuery := strings.ToUpper(strings.TrimSpace(query))
	isSelect := strings.HasPrefix(upperQuery, "SELECT") || strings.HasPrefix(upperQuery, "WITH") || strings.HasPrefix(upperQuery, "PRAGMA")

	if !isSelect {
		return nil, fmt.Errorf("only SELECT, WITH, and PRAGMA queries are allowed")
	}

	ctx := context.Background()
	rows, err := c.conn.QueryContext(ctx, query)
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
