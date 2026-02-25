package turso

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"datafrost/internal/models"

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

func (c *Client) GetTableData(tableName string) (*models.QueryResult, error) {
	query := fmt.Sprintf("SELECT * FROM %s LIMIT 100", tableName)
	return c.ExecuteQuery(query)
}

func (c *Client) ExecuteQuery(query string) (*models.QueryResult, error) {
	upperQuery := strings.ToUpper(strings.TrimSpace(query))
	isSelect := strings.HasPrefix(upperQuery, "SELECT") || strings.HasPrefix(upperQuery, "WITH") || strings.HasPrefix(upperQuery, "PRAGMA")

	if !isSelect {
		return nil, fmt.Errorf("only SELECT, WITH, and PRAGMA queries are allowed")
	}

	limitedQuery := query
	if !strings.Contains(upperQuery, "LIMIT") {
		limitedQuery = fmt.Sprintf("%s LIMIT 100", query)
	}

	ctx := context.Background()
	rows, err := c.conn.QueryContext(ctx, limitedQuery)
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
		Limited: len(resultRows) == 100,
	}, nil
}
