package main

import (
	"errors"
	"fmt"
	// "html/template"
	"net/http"
	"strconv"

	"snippetbox.sangdennis.com/internal/models"
)

// Change the signature of home() to be defined as a method against *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w) // Use the notFound() helper
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serveError(w, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}

	// // Initialize a slice containing the paths to the two files.
	// // The base template file must be the first in the slice.
	// files := []string{
	// 	"./ui/html/base.html",
	// 	"./ui/html/partials/nav.html",
	// 	"./ui/html/pages/home.html",
	// }

	// // Use template.ParseFiles() to read the template file into a
	// // template set. Pass the slice of files paths as a variadic parameter
	// // If an error occurs, log a detailed error msg and use  http.Error()
	// // to send a generic 500 Internal Server Error response to the user.
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serveError(w, err) // Use the serverError() helper
	// 	return
	// }

	// // Use ExecuteTemplate() method on the template set to write the template
	// // content as the response body. The last parameter to ExecuteTemplate()
	// // represents any dynamic data that we want to pass in, nil for now.
	// err = ts.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	app.serveError(w, err) // Use the serverError() helper
	// }
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper
		return
	}

	// Use the SnippetModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serveError(w, err)
		}
		return
	}

	// Write the snippet data as a plain-text HTTP response body.
	fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed) // Use the clientError() helper
		return
	}

	w.Write([]byte("Create a new snippet..."))
}