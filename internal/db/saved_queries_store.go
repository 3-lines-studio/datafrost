package db

import (
	"database/sql"
	"time"
)

type SavedQuery struct {
	ID           int64     `json:"id"`
	ConnectionID int64     `json:"connectionId"`
	Name         string    `json:"name"`
	Query        string    `json:"query"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type SavedQueriesStore struct {
	db *sql.DB
}

func NewSavedQueriesStore(db *sql.DB) *SavedQueriesStore {
	return &SavedQueriesStore{db: db}
}

func (s *SavedQueriesStore) ListByConnection(connectionID int64) ([]SavedQuery, error) {
	rows, err := s.db.Query(
		`SELECT id, connection_id, name, query, created_at, updated_at 
		 FROM saved_queries 
		 WHERE connection_id = ? 
		 ORDER BY updated_at DESC`,
		connectionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var queries []SavedQuery
	for rows.Next() {
		var q SavedQuery
		err := rows.Scan(
			&q.ID,
			&q.ConnectionID,
			&q.Name,
			&q.Query,
			&q.CreatedAt,
			&q.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		queries = append(queries, q)
	}

	return queries, rows.Err()
}

func (s *SavedQueriesStore) GetByID(id int64) (*SavedQuery, error) {
	var q SavedQuery
	err := s.db.QueryRow(
		`SELECT id, connection_id, name, query, created_at, updated_at 
		 FROM saved_queries 
		 WHERE id = ?`,
		id,
	).Scan(
		&q.ID,
		&q.ConnectionID,
		&q.Name,
		&q.Query,
		&q.CreatedAt,
		&q.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &q, nil
}

func (s *SavedQueriesStore) Create(connectionID int64, name, query string) (*SavedQuery, error) {
	result, err := s.db.Exec(
		`INSERT INTO saved_queries (connection_id, name, query) 
		 VALUES (?, ?, ?)`,
		connectionID, name, query,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

func (s *SavedQueriesStore) Update(id int64, name, query string) (*SavedQuery, error) {
	_, err := s.db.Exec(
		`UPDATE saved_queries 
		 SET name = ?, query = ?, updated_at = CURRENT_TIMESTAMP 
		 WHERE id = ?`,
		name, query, id,
	)
	if err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

func (s *SavedQueriesStore) Delete(id int64) error {
	_, err := s.db.Exec(`DELETE FROM saved_queries WHERE id = ?`, id)
	return err
}
