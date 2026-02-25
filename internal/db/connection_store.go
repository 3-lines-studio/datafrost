package db

import (
	"database/sql"
	"datafrost/internal/models"
	"fmt"
)

type ConnectionStore struct {
	db *sql.DB
}

func NewConnectionStore(db *sql.DB) *ConnectionStore {
	return &ConnectionStore{db: db}
}

func (s *ConnectionStore) Create(req models.CreateConnectionRequest) (*models.Connection, error) {
	result, err := s.db.Exec(
		"INSERT INTO connections (name, url, token) VALUES (?, ?, ?)",
		req.Name, req.URL, req.Token,
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
	err := s.db.QueryRow(
		"SELECT id, name, url, token, created_at FROM connections WHERE id = ?",
		id,
	).Scan(&conn.ID, &conn.Name, &conn.URL, &conn.Token, &conn.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	return &conn, nil
}

func (s *ConnectionStore) List() ([]models.Connection, error) {
	rows, err := s.db.Query(
		"SELECT id, name, url, token, created_at FROM connections ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list connections: %w", err)
	}
	defer rows.Close()

	var connections []models.Connection
	for rows.Next() {
		var conn models.Connection
		if err := rows.Scan(&conn.ID, &conn.Name, &conn.URL, &conn.Token, &conn.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan connection: %w", err)
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
	_, err := s.db.Exec(
		"UPDATE connections SET name = ?, url = ?, token = ? WHERE id = ?",
		req.Name, req.URL, req.Token, id,
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
