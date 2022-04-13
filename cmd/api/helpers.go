package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type envelope map[string]interface{}

//? extract and return id param from request body
func (app *application) readIDParams(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)

	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

//? pass data to be converted to json Object
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope) error {

	jsonObject, err := json.MarshalIndent(data, "", "\t")

	if err != nil {
		return err
	}

	jsonObject = append(jsonObject, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonObject)
	return nil
}
