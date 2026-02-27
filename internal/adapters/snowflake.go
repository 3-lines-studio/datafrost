package adapters

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	"github.com/3-lines-studio/datafrost/internal/models"

	_ "github.com/snowflakedb/gosnowflake"
	gosnowflake "github.com/snowflakedb/gosnowflake"
)

type SnowflakeAdapter struct {
	conn *sql.DB
}

func NewSnowflakeAdapterRegistration() models.AdapterRegistration {
	return models.AdapterRegistration{
		Info: models.AdapterInfo{
			Type:        "snowflake",
			Name:        "Snowflake",
			Description: "Snowflake data warehouse (SSO or key pair)",
			UIConfig: models.UIConfig{
				Modes: []models.UIMode{
					{
						Key:   "browser",
						Label: "External Browser (SSO)",
						Fields: []models.FieldConfig{
							{
								Key:         "account",
								Label:       "Account",
								Type:        "text",
								Required:    true,
								Placeholder: "myorg-myaccount",
							},
							{
								Key:         "user",
								Label:       "Username",
								Type:        "text",
								Required:    true,
								Placeholder: "you@company.com",
							},
							{
								Key:         "warehouse",
								Label:       "Warehouse",
								Type:        "text",
								Required:    false,
								Placeholder: "COMPUTE_WH",
							},
							{
								Key:         "database",
								Label:       "Database",
								Type:        "text",
								Required:    false,
								Placeholder: "MY_DATABASE",
							},
							{
								Key:         "schema",
								Label:       "Schema",
								Type:        "text",
								Required:    false,
								Placeholder: "PUBLIC",
							},
							{
								Key:         "role",
								Label:       "Role",
								Type:        "text",
								Required:    false,
								Placeholder: "MY_ROLE",
							},
						},
					},
					{
						Key:   "private_key",
						Label: "Private Key",
						Fields: []models.FieldConfig{
							{
								Key:         "account",
								Label:       "Account",
								Type:        "text",
								Required:    true,
								Placeholder: "myorg-myaccount",
							},
							{
								Key:         "user",
								Label:       "Username",
								Type:        "text",
								Required:    true,
								Placeholder: "you@company.com",
							},
							{
								Key:         "private_key_pem",
								Label:       "Private Key (PEM)",
								Type:        "textarea",
								Required:    true,
								Placeholder: "-----BEGIN PRIVATE KEY-----...",
							},
							{
								Key:         "private_key_passphrase",
								Label:       "Private Key Passphrase (optional)",
								Type:        "password",
								Required:    false,
								Placeholder: "",
							},
							{
								Key:         "warehouse",
								Label:       "Warehouse",
								Type:        "text",
								Required:    false,
								Placeholder: "COMPUTE_WH",
							},
							{
								Key:         "database",
								Label:       "Database",
								Type:        "text",
								Required:    false,
								Placeholder: "MY_DATABASE",
							},
							{
								Key:         "schema",
								Label:       "Schema",
								Type:        "text",
								Required:    false,
								Placeholder: "PUBLIC",
							},
							{
								Key:         "role",
								Label:       "Role",
								Type:        "text",
								Required:    false,
								Placeholder: "MY_ROLE",
							},
						},
					},
				},
				SupportsFile: false,
			},
		},
		Factory: func() models.DatabaseAdapter {
			return &SnowflakeAdapter{}
		},
	}
}

func (a *SnowflakeAdapter) Connect(credentials map[string]any) error {
	account, ok := credentials["account"].(string)
	if !ok || account == "" {
		return fmt.Errorf("account is required")
	}

	user, ok := credentials["user"].(string)
	if !ok || user == "" {
		return fmt.Errorf("user is required")
	}

	mode, _ := credentials["mode"].(string)
	if mode == "" {
		mode = "browser"
	}

	cfg := &gosnowflake.Config{
		Account: account,
		User:    user,
	}

	switch mode {
	case "private_key":
		pemData, ok := credentials["private_key_pem"].(string)
		if !ok || strings.TrimSpace(pemData) == "" {
			return fmt.Errorf("private key PEM is required")
		}
		passphrase, _ := credentials["private_key_passphrase"].(string)

		key, err := parseSnowflakePrivateKey(pemData, passphrase)
		if err != nil {
			return fmt.Errorf("failed to parse private key: %w", err)
		}
		cfg.PrivateKey = key
		cfg.Authenticator = gosnowflake.AuthTypeJwt

	default:
		cfg.Authenticator = gosnowflake.AuthTypeExternalBrowser
	}

	if warehouse, ok := credentials["warehouse"].(string); ok && warehouse != "" {
		cfg.Warehouse = warehouse
	}
	if database, ok := credentials["database"].(string); ok && database != "" {
		cfg.Database = database
	}
	if schema, ok := credentials["schema"].(string); ok && schema != "" {
		cfg.Schema = schema
	}
	if role, ok := credentials["role"].(string); ok && role != "" {
		cfg.Role = role
	}

	dsn, err := gosnowflake.DSN(cfg)
	if err != nil {
		return fmt.Errorf("failed to build snowflake DSN: %w", err)
	}

	db, err := sql.Open("snowflake", dsn)
	if err != nil {
		return fmt.Errorf("failed to open snowflake connection: %w", err)
	}

	a.conn = db
	return nil
}

func (a *SnowflakeAdapter) Close() error {
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}

func (a *SnowflakeAdapter) Ping() error {
	if a.conn == nil {
		return fmt.Errorf("not connected")
	}
	return a.conn.Ping()
}

func (a *SnowflakeAdapter) ListTables() ([]models.TableInfo, error) {
	rows, err := a.conn.Query(`
		SELECT TABLE_NAME, TABLE_TYPE
		FROM INFORMATION_SCHEMA.TABLES
		WHERE TABLE_SCHEMA = CURRENT_SCHEMA()
		AND TABLE_TYPE IN ('BASE TABLE', 'VIEW')
		ORDER BY TABLE_NAME
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list tables: %w", err)
	}
	defer rows.Close()

	var tables []models.TableInfo
	for rows.Next() {
		var t models.TableInfo
		var tableType string
		if err := rows.Scan(&t.Name, &tableType); err != nil {
			return nil, fmt.Errorf("failed to scan table: %w", err)
		}
		if tableType == "BASE TABLE" {
			t.Type = "table"
		} else {
			t.Type = "view"
		}
		t.FullName = t.Name
		tables = append(tables, t)
	}

	return tables, rows.Err()
}

// ListTree returns the database → schema → table hierarchy. If the
// current session already has a database and schema selected, it
// returns just that branch to avoid redundant discovery.
func (a *SnowflakeAdapter) ListTree() ([]models.TreeNode, error) {
	currentDB, currentSchema, err := a.currentContext()
	if err != nil {
		return nil, err
	}

	// If both database and schema are set, just list that schema's tables.
	if currentDB != "" && currentSchema != "" {
		tables, err := a.listTablesForSchema(currentDB, currentSchema)
		if err != nil {
			return nil, err
		}
		return []models.TreeNode{
			{
				Name:     currentDB,
				Type:     "database",
				FullName: currentDB,
				Children: []models.TreeNode{
					{
						Name:     currentSchema,
						Type:     "schema",
						FullName: fmt.Sprintf("%s.%s", currentDB, currentSchema),
						Children: tables,
					},
				},
			},
		}, nil
	}

	// If database is set but schema is not, list schemas within it.
	if currentDB != "" {
		schemas, err := a.listSchemas(currentDB)
		if err != nil {
			return nil, err
		}
		return []models.TreeNode{
			{
				Name:     currentDB,
				Type:     "database",
				FullName: currentDB,
				Children: schemas,
			},
		}, nil
	}

	// No database set: list all accessible databases and their schemas/tables.
	dbs, err := a.listDatabases()
	if err != nil {
		return nil, err
	}
	return dbs, nil
}

func (a *SnowflakeAdapter) GetTableData(tableName string, limit, offset int, filters []models.Filter) (*models.QueryResult, error) {
	whereClause, args := buildSnowflakeWhereClause(filters)

	count, err := a.getFilteredTableCount(tableName, whereClause, args)
	if err != nil {
		return nil, err
	}

	qualifiedName := quoteSnowflakeIdentifierPath(tableName)
	query := fmt.Sprintf(`SELECT * FROM %s`, qualifiedName)
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

func (a *SnowflakeAdapter) ExecuteQuery(query string) (*models.QueryResult, error) {
	return a.executeQueryWithArgs(query, nil)
}

func (a *SnowflakeAdapter) executeQueryWithArgs(query string, args []any) (*models.QueryResult, error) {
	upper := strings.ToUpper(strings.TrimSpace(query))
	if !strings.HasPrefix(upper, "SELECT") && !strings.HasPrefix(upper, "WITH") {
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

func (a *SnowflakeAdapter) currentContext() (string, string, error) {
	var dbName, schemaName sql.NullString
	if err := a.conn.QueryRow(`SELECT CURRENT_DATABASE(), CURRENT_SCHEMA()`).Scan(&dbName, &schemaName); err != nil {
		return "", "", fmt.Errorf("failed to get current context: %w", err)
	}
	return strings.TrimSpace(dbName.String), strings.TrimSpace(schemaName.String), nil
}

func (a *SnowflakeAdapter) listDatabases() ([]models.TreeNode, error) {
	rows, err := a.conn.Query(`SELECT DATABASE_NAME FROM SNOWFLAKE.INFORMATION_SCHEMA.DATABASES ORDER BY DATABASE_NAME`)
	if err != nil {
		return nil, fmt.Errorf("failed to list databases: %w", err)
	}
	defer rows.Close()

	var databases []models.TreeNode
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, fmt.Errorf("failed to scan database: %w", err)
		}
		schemaNodes, err := a.listSchemas(dbName)
		if err != nil {
			return nil, err
		}
		databases = append(databases, models.TreeNode{
			Name:     dbName,
			Type:     "database",
			FullName: dbName,
			Children: schemaNodes,
		})
	}

	return databases, rows.Err()
}

func (a *SnowflakeAdapter) listSchemas(database string) ([]models.TreeNode, error) {
	query := fmt.Sprintf(`SELECT SCHEMA_NAME FROM "%s".INFORMATION_SCHEMA.SCHEMATA ORDER BY SCHEMA_NAME`, database)
	rows, err := a.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list schemas for %s: %w", database, err)
	}
	defer rows.Close()

	var schemas []models.TreeNode
	for rows.Next() {
		var schemaName string
		if err := rows.Scan(&schemaName); err != nil {
			return nil, fmt.Errorf("failed to scan schema: %w", err)
		}
		tables, err := a.listTablesForSchema(database, schemaName)
		if err != nil {
			return nil, err
		}
		schemas = append(schemas, models.TreeNode{
			Name:     schemaName,
			Type:     "schema",
			FullName: fmt.Sprintf("%s.%s", database, schemaName),
			Children: tables,
		})
	}

	return schemas, rows.Err()
}

func (a *SnowflakeAdapter) listTablesForSchema(database, schema string) ([]models.TreeNode, error) {
	query := fmt.Sprintf(`SELECT TABLE_NAME, TABLE_TYPE FROM "%s".INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_TYPE IN ('BASE TABLE', 'VIEW') ORDER BY TABLE_NAME`, database)
	rows, err := a.conn.Query(query, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to list tables for %s.%s: %w", database, schema, err)
	}
	defer rows.Close()

	var tables []models.TreeNode
	for rows.Next() {
		var name, tableType string
		if err := rows.Scan(&name, &tableType); err != nil {
			return nil, fmt.Errorf("failed to scan table: %w", err)
		}
		typeVal := "table"
		if tableType == "VIEW" {
			typeVal = "view"
		}
		tables = append(tables, models.TreeNode{
			Name:     name,
			Type:     typeVal,
			FullName: fmt.Sprintf("%s.%s.%s", database, schema, name),
		})
	}

	return tables, rows.Err()
}

func parseSnowflakePrivateKey(pemString, passphrase string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemString))
	if block == nil {
		return nil, errors.New("invalid PEM data")
	}

	var keyBytes []byte
	var err error
	if x509.IsEncryptedPEMBlock(block) {
		if passphrase == "" {
			return nil, errors.New("private key is encrypted but no passphrase was provided")
		}
		keyBytes, err = x509.DecryptPEMBlock(block, []byte(passphrase))
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt private key: %w", err)
		}
	} else {
		keyBytes = block.Bytes
	}

	// Try PKCS8 first
	if key, err := x509.ParsePKCS8PrivateKey(keyBytes); err == nil {
		if rsaKey, ok := key.(*rsa.PrivateKey); ok {
			return rsaKey, nil
		}
		return nil, errors.New("private key is not RSA")
	}

	// Fallback to PKCS1
	if key, err := x509.ParsePKCS1PrivateKey(keyBytes); err == nil {
		return key, nil
	}

	return nil, errors.New("failed to parse private key")
}

func quoteSnowflakeIdentifierPath(name string) string {
	parts := strings.Split(name, ".")
	for i, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if strings.HasPrefix(p, "\"") && strings.HasSuffix(p, "\"") {
			parts[i] = p
			continue
		}
		parts[i] = fmt.Sprintf("\"%s\"", p)
	}
	return strings.Join(parts, ".")
}

func (a *SnowflakeAdapter) getFilteredTableCount(tableName, whereClause string, args []any) (int, error) {
	qualifiedName := quoteSnowflakeIdentifierPath(tableName)
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, qualifiedName)
	if whereClause != "" {
		countQuery += " WHERE " + whereClause
	}
	var count int
	if err := a.conn.QueryRow(countQuery, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count rows: %w", err)
	}
	return count, nil
}

func buildSnowflakeWhereClause(filters []models.Filter) (string, []any) {
	if len(filters) == 0 {
		return "", nil
	}

	var conditions []string
	var args []any

	for _, filter := range filters {
		if filter.Column == "" {
			continue
		}

		column := fmt.Sprintf(`"%s"`, filter.Column)

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
