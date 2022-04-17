package data

import (
	"database/sql"
	"errors"
	"fmt"
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

	query := `
	DELETE FROM todos
	WHERE id = $1
	`

	result, err := t.DB.Exec(query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
func (t TodoModel) GetAllTodoItems() ([]*Todo, error) {
	query := `
	SELECT * FROM todos
	`

	var todo Todo

	rows, err := t.DB.Query(query)

	if err != nil {
		return nil, ErrRecordNotFound
	}

	fmt.Println(rows, "rowssss")

	todos := []*Todo{}

	for rows.Next() {
		err = rows.Scan(
			&todo.ID,
			&todo.Item,
			&todo.Description,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&todo.Status,
		)

		if err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}

	fmt.Println("hjdwtgdhgwhdqyhjs")

	return todos, nil
}
