package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/todos", app.todo)
	router.HandlerFunc(http.MethodGet, "/v1/todos/:id", app.getTodoById)

	return router
}
