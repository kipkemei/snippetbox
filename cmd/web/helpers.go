package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

// the serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// the clientError helper sends a specific status code and corresponding description  to the user.
// It can be used to send responses like 400 "Bad Request" when there is a problem with a user request.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound helper is a convinience wrapper around clientError which sends a 404 Not Found reponse
// to the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	// Retrieve the appropriate template set from the cache based on the page name (like "home.html").
	// If no entry exists in the cache with the provided name, then create a new error and call the
	// serverError() helper method and return.
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// Initialize a new buffer
	buf := new(bytes.Buffer)

	// Write the template to the buffer, instead of straigh to the http.ResponseWriter.
	// If there's an error, call serverError() helper and then return.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// If the template is written to the buffer without any errors, we are safe to ahead and
	// write the HTTP status code to http.ResponseWriter.
	w.WriteHeader(status)

	// Write the contents of the buffer to the http.ResponseWriter.
	// NOTE: pass http.ResponseWriter to a function that takes io.Writer
	buf.WriteTo(w)
}
