package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/3-lines-studio/datafrost/internal/models"

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

func (a *TursoAdapter) Connect(credentials map[string]any) error {
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

func (a *TursoAdapter) executeQueryWithArgs(query string, args []any) (*models.QueryResult, error) {
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

func (a *TursoAdapter) getFilteredTableCount(tableName, whereClause string, args []any) (int, error) {
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

func (a *TursoAdapter) GetTableSchema(tableName string) (*models.TableSchema, error) {
	schema := &models.TableSchema{
		TableName: tableName,
	}

	escapedTableName := strings.ReplaceAll(tableName, "'", "''")

	columnRows, err := a.conn.Query(fmt.Sprintf("PRAGMA table_info('%s')", escapedTableName))
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}
	defer columnRows.Close()

	var columns []models.ColumnInfo
	for columnRows.Next() {
		var col models.ColumnInfo
		var cid int
		var notNull int
		var pk int
		var dfltValue sql.NullString

		if err := columnRows.Scan(&cid, &col.Name, &col.Type, &notNull, &dfltValue, &pk); err != nil {
			return nil, fmt.Errorf("failed to scan column: %w", err)
		}

		col.Nullable = notNull == 0
		col.IsPrimaryKey = pk == 1
		if dfltValue.Valid {
			col.DefaultValue = dfltValue.String
		}

		columns = append(columns, col)
	}
	schema.Columns = columns

	indexListRows, err := a.conn.Query(fmt.Sprintf("PRAGMA index_list('%s')", escapedTableName))
	if err != nil {
		return nil, fmt.Errorf("failed to get index list: %w", err)
	}
	defer indexListRows.Close()

	var indexes []models.IndexInfo
	for indexListRows.Next() {
		var seq int
		var indexName string
		var unique int
		var origin string
		var partial int

		if err := indexListRows.Scan(&seq, &indexName, &unique, &origin, &partial); err != nil {
			return nil, fmt.Errorf("failed to scan index list: %w", err)
		}

		idx := models.IndexInfo{
			Name:   indexName,
			Unique: unique == 1,
		}

		escapedIndexName := strings.ReplaceAll(indexName, "'", "''")
		indexInfoRows, err := a.conn.Query(fmt.Sprintf("PRAGMA index_info('%s')", escapedIndexName))
		if err != nil {
			continue
		}

		var cols []string
		for indexInfoRows.Next() {
			var seqno int
			var cid int
			var colName string
			if err := indexInfoRows.Scan(&seqno, &cid, &colName); err == nil && colName != "" {
				cols = append(cols, colName)
			}
		}
		indexInfoRows.Close()

		idx.Columns = cols
		indexes = append(indexes, idx)
	}
	schema.Indexes = indexes

	return schema, nil
}

func buildTursoWhereClause(filters []models.Filter) (string, []any) {
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
