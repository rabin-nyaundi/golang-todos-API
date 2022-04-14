package data

import (
	"database/sql"
	"time"
)

type Todo struct {
	ID          int64     `json:"id"`
	Item        string    `json:"item"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Status      bool      `json:"status"`
}

type TodoModel struct {
	DB *sql.DB
}

func (t TodoModel) InsertTodo(todo *Todo) error {
	query := `
		INSERT INTO todos (item, description, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, status
	`

	args := []interface{}{todo.Item, todo.Description, todo.Status}
	return t.DB.QueryRow(query, args...).Scan(&todo.ID, &todo.CreatedAt, &todo.Status)
}

func (t TodoModel) UpdateTodo(todo *Todo) error {
	return nil
}

func (t TodoModel) GetTodo(todo *Todo) error {
	return nil
}

func (t TodoModel) DeleteTodo(id int64) error {
	return nil
}

// curl -i -d "$BODY" localhost:4010/v1/todos
// BODY='{"item":"description":Buy greens","I go to the market and buy some greens","status":"false"}'
