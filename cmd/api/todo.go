package main

import (
	"fmt"
	"net/http"
	"todo/internal/data"
)

func (app *application) todo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environment: %s\n", app.config.env)
	fmt.Fprintf(w, "version: %s\n", version)

	// env := envelope{
	// 	"status": "available",
	// 	"system_info": map[string]string{
	// 		"environment": app.config.env,
	// 		"version":     version,
	// 	},
	// }

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

// func createTodoHandler() {

// }

func (app *application) getTodoById(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParams(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Println(id)

	todo := data.Todo{
		ID:          1,
		Item:        "uyejhdkjldf",
		Description: "Description of the todo",
		Status:      false,
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
