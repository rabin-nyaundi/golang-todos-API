package main

import (
	"errors"
	"fmt"
	"net/http"
	"todo/internal/data"

	"github.com/julienschmidt/httprouter"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.logger.PrintFatal(err, nil)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.User.Insert(user)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrorDuplicateEmail):
			app.logger.PrintFatal(errors.New("user with email already exists"), nil)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user})

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) GetUserByEmailHandler(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	email := params.ByName("email")

	user, err := app.models.User.GetByEmail(email)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.logger.PrintError(err, nil)
			return

		default:
			app.logger.PrintError(err, nil)
			return
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user})

	if err != nil {
		app.logger.PrintError(err, nil)
	}

	fmt.Println(email, "here is the email you are looking for")
}
