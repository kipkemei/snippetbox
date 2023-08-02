package main

import (
	"errors"
	"html/template"
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
		app.serverError(w, err)
		return
	}

	// Initialize a slice containing the paths to the two files.
	// The base template file must be the first in the slice.
	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}

	// Use template.ParseFiles() to read the template file into a
	// template set. Pass the slice of files paths as a variadic parameter
	// If an error occurs, log a detailed error msg and use  http.Error()
	// to send a generic 500 Internal Server Error response to the user.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err) // Use the serverError() helper
		return
	}

	// Create an instance of a templateData struct holding the slice of snippets.
	data := &templateData{
		Snippets: snippets,
	}

	// Use ExecuteTemplate() method on the template set to write the template
	// content as the response body. The last parameter to ExecuteTemplate()
	// Pass in the templateData struct when executing the template.
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err) // Use the serverError() helper
	}
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
			app.serverError(w, err)
		}
		return
	}

	// Initialize a slice containing the paths to view.html file,
	// plus the base layout and navigation partial.
	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/view.html",
	}

	// Parse the template files
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Create an instance of a templateData struct holding the snippet data.
	data := &templateData{
		Snippet: snippet,
	}

	// And then execute them. Pass in the snippet data a (models.Snippet struct)
	// as the final parameter.
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed) // Use the clientError() helper
		return
	}

	w.Write([]byte("Create a new snippet..."))
}