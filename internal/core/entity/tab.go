package entity

type Tab struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Title        string `json:"title"`
	ConnectionID int    `json:"connectionId"`
	TableName    string `json:"tableName,omitempty"`
	Query        string `json:"query,omitempty"`
	Page         int    `json:"page,omitempty"`
}
