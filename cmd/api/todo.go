package main

import (
	"errors"
	"fmt"
	"net/http"
	"todo/internal/data"
)

func (app *application) todo(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, envelope{"server": env})

	if err != nil {
		app.logger.Printf(err.Error())
	}

}
func (app *application) getAllTodos(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParams(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"todo id": id})

	if err != nil {
		app.logger.Printf("an error is here")
	}
}

func (app *application) createTodoHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Item        string `json:"item"`
		Description string `json:"description"`
		Status      bool   `json:"status"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.logger.Fatal(err)
	}

	todo := &data.Todo{
		Item:        input.Item,
		Description: input.Description,
		Status:      input.Status,
	}

	err = app.models.Todo.InsertTodo(todo)

	if err != nil {
		fmt.Println("Failed !!!!")
		app.logger.Printf(err.Error())
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/todos/%d", todo.ID))
	err = app.writeJSON(w, http.StatusCreated, envelope{"todo": todo})

	if err != nil {
		app.logger.Printf(err.Error())
		return
	}
}

func (app *application) getTodoById(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParams(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Println(id)

	todo, err := app.models.Todo.GetTodo(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}


	err = app.writeJSON(w, http.StatusOK, envelope{"todo": todo})

	if err != nil {
		app.logger.Printf("an error is here")
	}

}

// func updateTodoHandler() {

// }
// func deleteTodoHandler() {

// }
