package data

import "time"

type Todo struct {
	ID          int64     `json:"id"`
	Item        string    `json:"item"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Status      bool      `json:"status"`
}
