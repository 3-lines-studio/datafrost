package models

import "time"

type Connection struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateConnectionRequest struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Token string `json:"token"`
}

type UpdateConnectionRequest struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Token string `json:"token"`
}

type TestConnectionRequest struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

type TableInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type QueryResult struct {
	Columns []string `json:"columns"`
	Rows    [][]any  `json:"rows"`
	Count   int      `json:"count"`
	Total   int      `json:"total"`
	Page    int      `json:"page"`
	Limit   int      `json:"limit"`
}

type Filter struct {
	ID       string `json:"id"`
	Column   string `json:"column"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}
