package data

import (
	"database/sql"
	"errors"
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

func (t TodoModel) GetTodo(id int64) (*Todo, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT * FROM todos
	WHERE id = $1
	`

	var todo Todo

	err := t.DB.QueryRow(query, id).Scan(
		&todo.ID,
		&todo.Item,
		&todo.Description,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.Status,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound

		default:
			return nil, err
		}
	}
	return &todo, nil
}

func (t TodoModel) UpdateTodo(todo *Todo) error {
	query := `
	UPDATE todos
	SET item = $1, description = $2, status = $3
	WHERE id = $4
	RETURNING item
	`

	args := []interface{}{
		todo.Item,
		todo.Description,
		todo.Status,
		todo.ID,
	}
	return t.DB.QueryRow(query, args...).Scan(&todo.Item)
}

func (t TodoModel) DeleteTodo(id int64) error {
	return nil
}

// curl -i -d "$BODY" localhost:4010/v1/todos
// BODY='{"item":"description":Buy greens","I go to the market and buy some greens","status":"false"}'
