package repository

import (
	"database/sql"
	"time"

	"github.com/3-lines-studio/datafrost/internal/core/entity"
)

type SavedQueryRepository struct {
	db *sql.DB
}

func NewSavedQueryRepository(db *sql.DB) *SavedQueryRepository {
	return &SavedQueryRepository{db: db}
}

func (r *SavedQueryRepository) ListByConnection(connectionID int64) ([]entity.SavedQuery, error) {
	rows, err := r.db.Query(
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

	var queries []entity.SavedQuery
	for rows.Next() {
		var q entity.SavedQuery
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

func (r *SavedQueryRepository) GetByID(id int64) (*entity.SavedQuery, error) {
	var q entity.SavedQuery
	err := r.db.QueryRow(
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

func (r *SavedQueryRepository) Create(connectionID int64, name, query string) (*entity.SavedQuery, error) {
	result, err := r.db.Exec(
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

	return r.GetByID(id)
}

func (r *SavedQueryRepository) Update(id int64, name, query string) (*entity.SavedQuery, error) {
	_, err := r.db.Exec(
		`UPDATE saved_queries 
		 SET name = ?, query = ?, updated_at = ? 
		 WHERE id = ?`,
		name, query, time.Now(), id,
	)
	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *SavedQueryRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM saved_queries WHERE id = ?`, id)
	return err
}
