package entity

import "time"

type SavedQuery struct {
	ID           int64     `json:"id"`
	ConnectionID int64     `json:"connectionId"`
	Name         string    `json:"name"`
	Query        string    `json:"query"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
