package adapters

import (
	"context"
	"fmt"
	"strings"

	"github.com/3-lines-studio/datafrost/internal/models"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/option"
)

type BigQueryAdapter struct {
	client    *bigquery.Client
	projectID string
	dataset   string
}

func NewBigQueryAdapterRegistration() models.AdapterRegistration {
	return models.AdapterRegistration{
		Info: models.AdapterInfo{
			Type:        "bigquery",
			Name:        "BigQuery",
			Description: "Google BigQuery database",
			UIConfig: models.UIConfig{
				Fields: []models.FieldConfig{
					{
						Key:         "project_id",
						Label:       "Project ID",
						Type:        "text",
						Required:    true,
						Placeholder: "my-project-id",
					},
					{
						Key:         "dataset",
						Label:       "Dataset",
						Type:        "text",
						Required:    true,
						Placeholder: "my_dataset",
					},
					{
						Key:         "credentials",
						Label:       "Service Account Credentials (JSON)",
						Type:        "textarea",
						Required:    true,
						Placeholder: "Paste JSON credentials here or upload file...",
					},
				},
				SupportsFile: true,
				FileTypes:    []string{".json"},
			},
		},
		Factory: func() models.DatabaseAdapter {
			return &BigQueryAdapter{}
		},
	}
}

func (a *BigQueryAdapter) Connect(credentials map[string]any) error {
	projectID, ok := credentials["project_id"].(string)
	if !ok || projectID == "" {
		return fmt.Errorf("project_id is required")
	}

	dataset, ok := credentials["dataset"].(string)
	if !ok || dataset == "" {
		return fmt.Errorf("dataset is required")
	}

	credJSON, ok := credentials["credentials"].(string)
	if !ok || credJSON == "" {
		return fmt.Errorf("credentials are required")
	}

	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID, option.WithCredentialsJSON([]byte(credJSON)))
	if err != nil {
		return fmt.Errorf("failed to create bigquery client: %w", err)
	}

	a.client = client
	a.projectID = projectID
	a.dataset = dataset
	return nil
}

func (a *BigQueryAdapter) Close() error {
	if a.client != nil {
		return a.client.Close()
	}
	return nil
}

func (a *BigQueryAdapter) Ping() error {
	if a.client == nil {
		return fmt.Errorf("not connected")
	}

	ctx := context.Background()
	_, err := a.client.Dataset(a.dataset).Metadata(ctx)
	if err != nil {
		return fmt.Errorf("dataset not found or no access: %w", err)
	}
	return nil
}

func (a *BigQueryAdapter) ListTables() ([]models.TableInfo, error) {
	if a.client == nil {
		return nil, fmt.Errorf("not connected")
	}

	ctx := context.Background()
	dataset := a.client.Dataset(a.dataset)

	it := dataset.Tables(ctx)
	var tables []models.TableInfo

	for {
		table, err := it.Next()
		if err != nil {
			break
		}
		tables = append(tables, models.TableInfo{
			Name: table.TableID,
			Type: "table",
		})
	}

	return tables, nil
}

func (a *BigQueryAdapter) GetTableData(tableName string, limit, offset int, filters []models.Filter) (*models.QueryResult, error) {
	if a.client == nil {
		return nil, fmt.Errorf("not connected")
	}

	whereClause, _ := buildBigQueryWhereClause(filters)

	countQuery := fmt.Sprintf("SELECT COUNT(*) as count FROM `%s.%s.%s`", a.projectID, a.dataset, tableName)
	if whereClause != "" {
		countQuery += " WHERE " + whereClause
	}

	ctx := context.Background()
	countResult, err := a.executeQueryWithCount(countQuery)
	if err != nil {
		return nil, err
	}

	_ = ctx

	totalCount := int64(0)
	if len(countResult.Rows) > 0 && len(countResult.Rows[0]) > 0 {
		switch v := countResult.Rows[0][0].(type) {
		case int64:
			totalCount = v
		case int:
			totalCount = int64(v)
		case float64:
			totalCount = int64(v)
		}
	}

	query := fmt.Sprintf("SELECT * FROM `%s.%s.%s`", a.projectID, a.dataset, tableName)
	if whereClause != "" {
		query += " WHERE " + whereClause
	}
	query += fmt.Sprintf(" LIMIT %d", limit)
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", offset)
	}

	result, err := a.ExecuteQuery(query)
	if err != nil {
		return nil, err
	}

	result.Total = int(totalCount)
	result.Page = offset/limit + 1
	result.Limit = limit
	return result, nil
}

func (a *BigQueryAdapter) ExecuteQuery(query string) (*models.QueryResult, error) {
	if a.client == nil {
		return nil, fmt.Errorf("not connected")
	}

	upperQuery := strings.ToUpper(strings.TrimSpace(query))
	isSelect := strings.HasPrefix(upperQuery, "SELECT") || strings.HasPrefix(upperQuery, "WITH")

	if !isSelect {
		return nil, fmt.Errorf("only SELECT and WITH queries are allowed")
	}

	return a.executeQueryWithCount(query)
}

func (a *BigQueryAdapter) executeQueryWithCount(query string) (*models.QueryResult, error) {
	ctx := context.Background()
	q := a.client.Query(query)
	it, err := q.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	schema := it.Schema
	columns := make([]string, len(schema))
	for i, field := range schema {
		columns[i] = field.Name
	}

	var resultRows [][]any
	for {
		var row []bigquery.Value
		err := it.Next(&row)
		if err != nil {
			break
		}

		convertedRow := make([]any, len(row))
		for i, val := range row {
			convertedRow[i] = convertBigQueryValue(val)
		}
		resultRows = append(resultRows, convertedRow)
	}

	if len(columns) == 0 {
		columns = a.getFallbackColumns(ctx, query, resultRows)
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

func convertBigQueryValue(val bigquery.Value) any {
	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case string:
		return v
	case int64:
		return v
	case float64:
		return v
	case bool:
		return v
	case []byte:
		return string(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func buildBigQueryWhereClause(filters []models.Filter) (string, []any) {
	if len(filters) == 0 {
		return "", nil
	}

	var conditions []string
	var args []any

	for _, filter := range filters {
		if filter.Column == "" {
			continue
		}

		switch filter.Operator {
		case "eq":
			conditions = append(conditions, fmt.Sprintf("%s = '%s'", filter.Column, filter.Value))
		case "neq":
			conditions = append(conditions, fmt.Sprintf("%s != '%s'", filter.Column, filter.Value))
		case "gt":
			conditions = append(conditions, fmt.Sprintf("%s > '%s'", filter.Column, filter.Value))
		case "lt":
			conditions = append(conditions, fmt.Sprintf("%s < '%s'", filter.Column, filter.Value))
		case "gte":
			conditions = append(conditions, fmt.Sprintf("%s >= '%s'", filter.Column, filter.Value))
		case "lte":
			conditions = append(conditions, fmt.Sprintf("%s <= '%s'", filter.Column, filter.Value))
		case "like":
			conditions = append(conditions, fmt.Sprintf("%s LIKE '%s'", filter.Column, filter.Value))
		case "not_like":
			conditions = append(conditions, fmt.Sprintf("%s NOT LIKE '%s'", filter.Column, filter.Value))
		case "is_null":
			conditions = append(conditions, fmt.Sprintf("%s IS NULL", filter.Column))
		case "is_not_null":
			conditions = append(conditions, fmt.Sprintf("%s IS NOT NULL", filter.Column))
		}
	}

	if len(conditions) == 0 {
		return "", nil
	}

	return strings.Join(conditions, " AND "), args
}

func (a *BigQueryAdapter) getFallbackColumns(ctx context.Context, query string, resultRows [][]any) []string {
	tableName := extractTableNameFromQuery(query)
	if tableName != "" {
		columns, err := a.getColumnsFromInfoSchema(ctx, tableName)
		if err == nil && len(columns) > 0 {
			return columns
		}
	}

	if len(resultRows) > 0 {
		columns := make([]string, len(resultRows[0]))
		for i := range columns {
			columns[i] = fmt.Sprintf("col_%d", i)
		}
		return columns
	}

	return []string{}
}

func extractTableNameFromQuery(query string) string {
	upperQuery := strings.ToUpper(strings.TrimSpace(query))

	if strings.HasPrefix(upperQuery, "SELECT") {
		fromIndex := strings.Index(upperQuery, " FROM ")
		if fromIndex == -1 {
			return ""
		}

		afterFrom := strings.TrimSpace(query[fromIndex+6:])

		endIndex := len(afterFrom)
		for _, keyword := range []string{" WHERE ", " GROUP ", " ORDER ", " LIMIT ", " HAVING ", " JOIN ", " UNION ", " INTERSECT ", " EXCEPT "} {
			idx := strings.Index(strings.ToUpper(afterFrom), keyword)
			if idx != -1 && idx < endIndex {
				endIndex = idx
			}
		}

		tablePart := strings.TrimSpace(afterFrom[:endIndex])
		tablePart = strings.Trim(tablePart, "`\"'")

		if strings.Contains(tablePart, ",") || strings.Contains(tablePart, "(") {
			return ""
		}

		parts := strings.Split(tablePart, ".")
		return parts[len(parts)-1]
	}

	return ""
}

func (a *BigQueryAdapter) getColumnsFromInfoSchema(ctx context.Context, tableName string) ([]string, error) {
	infoQuery := fmt.Sprintf(
		"SELECT column_name FROM `%s.%s.INFORMATION_SCHEMA.COLUMNS` WHERE table_name = '%s' ORDER BY ordinal_position",
		a.projectID, a.dataset, tableName,
	)

	q := a.client.Query(infoQuery)
	it, err := q.Read(ctx)
	if err != nil {
		return nil, err
	}

	var columns []string
	for {
		var row []bigquery.Value
		err := it.Next(&row)
		if err != nil {
			break
		}
		if len(row) > 0 && row[0] != nil {
			columns = append(columns, fmt.Sprintf("%v", row[0]))
		}
	}

	return columns, nil
}

func (a *BigQueryAdapter) GetTableSchema(tableName string) (*models.TableSchema, error) {
	if a.client == nil {
		return nil, fmt.Errorf("not connected")
	}

	ctx := context.Background()

	table := a.client.Dataset(a.dataset).Table(tableName)
	metadata, err := table.Metadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get table metadata: %w", err)
	}

	var columns []models.ColumnInfo
	for _, field := range metadata.Schema {
		columns = append(columns, models.ColumnInfo{
			Name:     field.Name,
			Type:     string(field.Type),
			Nullable: !field.Required,
		})
	}

	return &models.TableSchema{
		TableName: tableName,
		Columns:   columns,
	}, nil
}
