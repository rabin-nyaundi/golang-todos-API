package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Todo TodoModel
	User UsersModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Todo: TodoModel{DB: db},
		User: UsersModel{DB: db},
	}
}
