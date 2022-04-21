package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/todos", app.listTodoHandler)
	router.HandlerFunc(http.MethodGet, "/v1/todos/:id", app.getTodoByIdHandler)
	router.HandlerFunc(http.MethodPost, "/v1/todos", app.createTodoHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/todos/:id", app.updateTodoHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/todos/:id", app.deleteTodoHandler)
	// router.HandlerFunc(http.MethodGet, "/v1/todos/all", app.getAllTodoItemshandler)
	return app.recoverPanic(app.rateLimiter(router))
}
