package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/3-lines-studio/datafrost/internal/core/entity"
)

type ConnectionRepository struct {
	db *sql.DB
}

func NewConnectionRepository(db *sql.DB) *ConnectionRepository {
	return &ConnectionRepository{db: db}
}

func (r *ConnectionRepository) Create(req entity.CreateConnectionRequest) (*entity.Connection, error) {
	credentialsJSON, err := serializeCredentials(req.Credentials)
	if err != nil {
		return nil, err
	}

	result, err := r.db.Exec(
		"INSERT INTO connections (name, type, credentials) VALUES (?, ?, ?)",
		req.Name, req.Type, credentialsJSON,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return r.GetByID(id)
}

func (r *ConnectionRepository) GetByID(id int64) (*entity.Connection, error) {
	var conn entity.Connection
	var credentialsJSON string
	err := r.db.QueryRow(
		"SELECT id, name, type, credentials, created_at FROM connections WHERE id = ?",
		id,
	).Scan(&conn.ID, &conn.Name, &conn.Type, &credentialsJSON, &conn.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	conn.Credentials, err = deserializeCredentials(credentialsJSON)
	if err != nil {
		return nil, err
	}

	return &conn, nil
}

func (r *ConnectionRepository) List() ([]entity.Connection, error) {
	rows, err := r.db.Query(
		"SELECT id, name, type, credentials, created_at FROM connections ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list connections: %w", err)
	}
	defer rows.Close()

	var connections []entity.Connection
	for rows.Next() {
		var conn entity.Connection
		var credentialsJSON string
		if err := rows.Scan(&conn.ID, &conn.Name, &conn.Type, &credentialsJSON, &conn.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan connection: %w", err)
		}
		conn.Credentials, err = deserializeCredentials(credentialsJSON)
		if err != nil {
			return nil, err
		}
		connections = append(connections, conn)
	}

	return connections, rows.Err()
}

func (r *ConnectionRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM connections WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete connection: %w", err)
	}
	return nil
}

func (r *ConnectionRepository) Update(id int64, req entity.UpdateConnectionRequest) (*entity.Connection, error) {
	credentialsJSON, err := serializeCredentials(req.Credentials)
	if err != nil {
		return nil, err
	}

	_, err = r.db.Exec(
		"UPDATE connections SET name = ?, type = ?, credentials = ? WHERE id = ?",
		req.Name, req.Type, credentialsJSON, id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update connection: %w", err)
	}

	return r.GetByID(id)
}

func (r *ConnectionRepository) SetLastConnected(id int64) error {
	_, err := r.db.Exec(
		"INSERT INTO app_state (key, value) VALUES ('last_connected_id', ?) ON CONFLICT(key) DO UPDATE SET value = ?",
		fmt.Sprintf("%d", id), fmt.Sprintf("%d", id),
	)
	return err
}

func (r *ConnectionRepository) GetLastConnected() (int64, error) {
	var value string
	err := r.db.QueryRow(
		"SELECT value FROM app_state WHERE key = 'last_connected_id'",
	).Scan(&value)

	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	var id int64
	_, err = fmt.Sscanf(value, "%d", &id)
	return id, err
}

func serializeCredentials(credentials map[string]any) (string, error) {
	data, err := json.Marshal(credentials)
	if err != nil {
		return "", fmt.Errorf("failed to serialize credentials: %w", err)
	}
	return string(data), nil
}

func deserializeCredentials(data string) (map[string]any, error) {
	var credentials map[string]any
	if err := json.Unmarshal([]byte(data), &credentials); err != nil {
		return nil, fmt.Errorf("failed to deserialize credentials: %w", err)
	}
	return credentials, nil
}
