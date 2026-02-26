package db

import (
	"database/sql"
	"fmt"

	"github.com/3-lines-studio/datafrost/internal/adapters"
	"github.com/3-lines-studio/datafrost/internal/models"
)

type ConnectionStore struct {
	db *sql.DB
}

func NewConnectionStore(db *sql.DB) *ConnectionStore {
	return &ConnectionStore{db: db}
}

func (s *ConnectionStore) Create(req models.CreateConnectionRequest) (*models.Connection, error) {
	credentialsJSON, err := adapters.SerializeCredentials(req.Credentials)
	if err != nil {
		return nil, err
	}

	result, err := s.db.Exec(
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

	return s.GetByID(id)
}

func (s *ConnectionStore) GetByID(id int64) (*models.Connection, error) {
	var conn models.Connection
	var credentialsJSON string
	err := s.db.QueryRow(
		"SELECT id, name, type, credentials, created_at FROM connections WHERE id = ?",
		id,
	).Scan(&conn.ID, &conn.Name, &conn.Type, &credentialsJSON, &conn.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	conn.Credentials, err = adapters.DeserializeCredentials(credentialsJSON)
	if err != nil {
		return nil, err
	}

	return &conn, nil
}

func (s *ConnectionStore) List() ([]models.Connection, error) {
	rows, err := s.db.Query(
		"SELECT id, name, type, credentials, created_at FROM connections ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list connections: %w", err)
	}
	defer rows.Close()

	var connections []models.Connection
	for rows.Next() {
		var conn models.Connection
		var credentialsJSON string
		if err := rows.Scan(&conn.ID, &conn.Name, &conn.Type, &credentialsJSON, &conn.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan connection: %w", err)
		}
		conn.Credentials, err = adapters.DeserializeCredentials(credentialsJSON)
		if err != nil {
			return nil, err
		}
		connections = append(connections, conn)
	}

	return connections, rows.Err()
}

func (s *ConnectionStore) Delete(id int64) error {
	_, err := s.db.Exec("DELETE FROM connections WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete connection: %w", err)
	}
	return nil
}

func (s *ConnectionStore) Update(id int64, req models.UpdateConnectionRequest) (*models.Connection, error) {
	credentialsJSON, err := adapters.SerializeCredentials(req.Credentials)
	if err != nil {
		return nil, err
	}

	_, err = s.db.Exec(
		"UPDATE connections SET name = ?, type = ?, credentials = ? WHERE id = ?",
		req.Name, req.Type, credentialsJSON, id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update connection: %w", err)
	}

	return s.GetByID(id)
}

func (s *ConnectionStore) SetLastConnected(id int64) error {
	_, err := s.db.Exec(
		"INSERT INTO app_state (key, value) VALUES ('last_connected_id', ?) ON CONFLICT(key) DO UPDATE SET value = ?",
		fmt.Sprintf("%d", id), fmt.Sprintf("%d", id),
	)
	return err
}

func (s *ConnectionStore) GetLastConnected() (int64, error) {
	var value string
	err := s.db.QueryRow(
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
