package entity

import "time"

type Connection struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Credentials map[string]any `json:"credentials"`
	CreatedAt   time.Time      `json:"created_at"`
}

type CreateConnectionRequest struct {
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Credentials map[string]any `json:"credentials"`
}

type UpdateConnectionRequest struct {
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Credentials map[string]any `json:"credentials"`
}

type TestConnectionRequest struct {
	Type        string         `json:"type"`
	Credentials map[string]any `json:"credentials"`
}
