package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/3-lines-studio/datafrost/internal/core/entity"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteAdapter struct {
	conn *sql.DB
}

func newSQLiteAdapterRegistration() entity.AdapterRegistration {
	return entity.AdapterRegistration{
		Info: entity.AdapterInfo{
			Type:        "sqlite",
			Name:        "SQLite File",
			Description: "Local SQLite database file",
			UIConfig: entity.UIConfig{
				Fields: []entity.FieldConfig{
					{
						Key:         "path",
						Label:       "Database File Path",
						Type:        "text",
						Required:    true,
						Placeholder: "/path/to/database.db or ./relative/path.db",
					},
				},
			},
		},
		Factory: func() entity.DatabaseAdapter {
			return &sqliteAdapter{}
		},
	}
}

func (a *sqliteAdapter) Connect(credentials map[string]any) error {
	path, ok := credentials["path"].(string)
	if !ok || path == "" {
		return fmt.Errorf("path is required")
	}

	database, err := sql.Open("sqlite3", path)
	if err != nil {
		return fmt.Errorf("failed to open sqlite connection: %w", err)
	}

	a.conn = database
	return nil
}

func (a *sqliteAdapter) Close() error {
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}

func (a *sqliteAdapter) Ping() error {
	if a.conn == nil {
		return fmt.Errorf("not connected")
	}
	return a.conn.Ping()
}

func (a *sqliteAdapter) ListTables() ([]entity.TableInfo, error) {
	rows, err := a.conn.Query(
		"SELECT name, type FROM sqlite_master WHERE type IN ('table', 'view') AND name NOT LIKE 'sqlite_%' ORDER BY name",
	)
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

func (a *sqliteAdapter) GetTableData(tableName string, limit, offset int, filters []entity.Filter) (*entity.QueryResult, error) {
	whereClause, args := buildSQLiteWhereClause(filters)

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

func (a *sqliteAdapter) ExecuteQuery(query string) (*entity.QueryResult, error) {
	return a.executeQueryWithArgs(query, nil)
}

func (a *sqliteAdapter) executeQueryWithArgs(query string, args []any) (*entity.QueryResult, error) {
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

	return &entity.QueryResult{
		Columns: columns,
		Rows:    resultRows,
		Count:   len(resultRows),
		Total:   len(resultRows),
		Page:    1,
		Limit:   len(resultRows),
	}, nil
}

func (a *sqliteAdapter) getFilteredTableCount(tableName, whereClause string, args []any) (int, error) {
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

func (a *sqliteAdapter) GetTableSchema(tableName string) (*entity.TableSchema, error) {
	schema := &entity.TableSchema{
		TableName: tableName,
	}

	escapedTableName := strings.ReplaceAll(tableName, "'", "''")

	columnRows, err := a.conn.Query(fmt.Sprintf("PRAGMA table_info('%s')", escapedTableName))
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}
	defer columnRows.Close()

	var columns []entity.ColumnInfo
	for columnRows.Next() {
		var col entity.ColumnInfo
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

	var indexes []entity.IndexInfo
	for indexListRows.Next() {
		var seq int
		var indexName string
		var unique int
		var origin string
		var partial int

		if err := indexListRows.Scan(&seq, &indexName, &unique, &origin, &partial); err != nil {
			return nil, fmt.Errorf("failed to scan index list: %w", err)
		}

		idx := entity.IndexInfo{
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

func buildSQLiteWhereClause(filters []entity.Filter) (string, []any) {
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
