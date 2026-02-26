package adapters

import (
	"context"
	"datafrost/internal/models"
	"fmt"
	"strings"

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

func (a *BigQueryAdapter) Connect(credentials map[string]interface{}) error {
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

func buildBigQueryWhereClause(filters []models.Filter) (string, []interface{}) {
	if len(filters) == 0 {
		return "", nil
	}

	var conditions []string
	var args []interface{}

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
