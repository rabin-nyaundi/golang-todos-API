package data

import (
	"context"
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	args := []interface{}{todo.Item, todo.Description, todo.Status}

	return t.DB.QueryRowContext(ctx, query, args...).Scan(&todo.ID, &todo.CreatedAt, &todo.Status)
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := t.DB.QueryRowContext(ctx, query, id).Scan(
		// &[]byte{},
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
			fmt.Println("Error with timeout here")
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	args := []interface{}{
		todo.Item,
		todo.Description,
		todo.Status,
		todo.ID,
	}
	return t.DB.QueryRowContext(ctx, query, args...).Scan(&todo.Item)
}

func (t TodoModel) DeleteTodo(id int64) error {

	query := `
	DELETE FROM todos
	WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	result, err := t.DB.ExecContext(ctx, query, id)

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
func (t TodoModel) GetAllTodoItems(item string, filters Filters) ([]*Todo, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), * FROM todos
		WHERE (STRPOS(LOWER(item), LOWER($1)) > 0 OR $1 = '')
		ORDER BY %s %s, id ASC LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	fmt.Println(filters.sortColumn())

	//  (LOWER(item) = LOWER($1) OR $1 = '')
	// (to_tsvector('simple',item) @@ plainto_tsquery('simple',$1) OR $1='')

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	args := []interface{}{item, filters.limit(), filters.offset()}

	rows, err := t.DB.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, Metadata{}, ErrRecordNotFound
	}

	defer rows.Close()

	totalRecords := 0
	todos := []*Todo{}

	for rows.Next() {

		var todo Todo

		err = rows.Scan(
			&totalRecords,
			&todo.ID,
			&todo.Item,
			&todo.Description,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&todo.Status,
		)

		if err != nil {
			return nil, Metadata{}, err
		}
		todos = append(todos, &todo)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(filters.Page, filters.PageSize, totalRecords)

	fmt.Println(metadata)

	return todos, metadata, nil
}
