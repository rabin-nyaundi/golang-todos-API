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
		app.logger.PrintFatal(err, nil)
	}

}

func (app *application) listTodoHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Item string
		data.Filters
	}

	qs := r.URL.Query()

	input.Item = app.readString(qs, "item", "")
	input.Filters.Page = app.readInt(qs, "page", 1)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "item", "description", "status", "-id", "-item", "-status"}

	todos, metadata, err := app.models.Todo.GetAllTodoItems(input.Item, input.Filters)

	if err != nil {
		app.logger.PrintFatal(err, nil)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "todos": todos})

	if err != nil {
		app.logger.PrintError(err, nil)
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
		app.logger.PrintFatal(err, nil)
	}

	todo := &data.Todo{
		Item:        input.Item,
		Description: input.Description,
		Status:      input.Status,
	}

	err = app.models.Todo.InsertTodo(todo)

	if err != nil {
		app.logger.PrintError(err, nil)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/todos/%d", todo.ID))
	err = app.writeJSON(w, http.StatusCreated, envelope{"todo": todo})

	if err != nil {
		app.logger.PrintError(err, nil)
		return
	}
}

func (app *application) getTodoByIdHandler(w http.ResponseWriter, r *http.Request) {

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
		app.logger.PrintError(err, nil)
	}

}

func (app *application) updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParams(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

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

	var input struct {
		Item        *string `json:"item"`
		Description *string `json:"description"`
		Status      *bool   `json:"status"`
	}

	err = app.readJSON(w, r, &input)

	if err != nil {
		fmt.Println("Hey it failed why fail here", err)
		app.notFoundResponse(w, r)
		return
	}

	if input.Item != nil {
		todo.Item = *input.Item
	}

	if input.Description != nil {
		todo.Description = *input.Description
	}

	if input.Status != nil {
		todo.Status = *input.Status
	}

	err = app.models.Todo.UpdateTodo(todo)

	if err != nil {
		fmt.Println("Hyey it failed here")
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"todo": todo})

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParams(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Todo.DeleteTodo(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "todo deleted successfully"})

	if err != nil {
		app.logger.PrintError(err, nil)
		app.serverErrorResponse(w, r, err)
	}
}
