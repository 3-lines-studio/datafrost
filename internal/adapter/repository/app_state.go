package repository

import "database/sql"

type AppStateRepository struct {
	db *sql.DB
}

func NewAppStateRepository(db *sql.DB) *AppStateRepository {
	return &AppStateRepository{db: db}
}

func (r *AppStateRepository) Get(key string) (string, error) {
	var value string
	err := r.db.QueryRow("SELECT value FROM app_state WHERE key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

func (r *AppStateRepository) Set(key, value string) error {
	_, err := r.db.Exec(
		"INSERT INTO app_state (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = ?",
		key, value, value,
	)
	return err
}
