package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// the serveError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serveError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// the clientError helper sends a specific status code and corresponding description  to the user. 
//It can be used to send responses like 400 "Bad Request" when there is a problem with a user request.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound helper is a convinience wrapper around clientError which sends a 404 Not Found reponse
// to the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}