package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"snippetbox.sangdennis.com/internal/models"
	"snippetbox.sangdennis.com/internal/validator"
)

// Change the signature of home() to be defined as a method against *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Because httprouter matches the "/" path exactly, we can now remove the manual
	// check below from this handler.
	// if r.URL.Path != "/" {
	// 	app.notFound(w) // Use the notFound() helper
	// 	return
	// }

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Call the newTemplateData() helper to get a templateData struct containing
	// the 'default' data (which for now is just the current year), and add the
	// snippets slice to it.
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// Use the render helper
	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// When httprouter is parsing a request, the values of any named parameters
	// will be stored in the request context.
	// You can use ParamsFromContext() function to retrieve a slice containing
	// these parameter names and values like so:
	params := httprouter.ParamsFromContext(r.Context())

	// We can then use the ByName() method to get the value of the "id" named
	// parameter from the slice and validate it as normal.
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	// Use the new render helper.
	app.render(w, http.StatusOK, "view.html", data)
}

// Add a new snippetCreate handler, which for now returns a placeholder response.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	// Initialize a new createSnippetForm instance and pass it to the template.
	// Set any default values for the form e.g snippet expiry to 365 days.
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.html", data)
}

// Define a snippetCreateForm struct to represent the form data and validation errors
// for the form fields. NOTE: All struct fields are deliberately exported (i.e start
// with capitral letter). This is because struct fields must be exported in order to
// be read by html/template package when rendering the template.
// Embed the Validator type. Embedding this means that our snippetCreateForm "inherits"
// all the fields and methods of Validator type (including FieldErrors field)
type snippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// First call r.ParseForm() which adds any data in POST request bodies to the
	// r.PostForm map. This also works in the same way for PUT and PATCH requests.
	// If there are any errors, use app.clientError() helper to send a 400 Bad Request
	// response to the user.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// The r.PostForm.Get() method always returns the form data as a *string*.
	// However, expires value should be a number to be represented as an integer.
	// Manually convert the form data to an integer using strconv.Atoi(), send
	// a 400 Bad Request response if the conversion fails.
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	// Create an instance of snippetCreateForm struct containing the values from
	// the form and an empty map for any validation errors.
	form := snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}

	// Call CheckField() directly to execute validation checks because Validation type
	// is embedded by the snippetCreateForm struct.
	// CheckField() adds the provided key and error message to the FieldErrors map if
	// the check does not evaluate to true.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank.")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long.")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank.")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365.")

	// Use the Valid() method to see if any of the checks failed.
	// If they did, re-render the template, passing in the form in the same way as before.
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	// Pass the data from snippetCreateForm instance to Insert() method
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect path to use the new clean URL format.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
