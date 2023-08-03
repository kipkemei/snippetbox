package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Intialize the router
	router := httprouter.New()

	// Create a handler function which wraps notFound() helper, and then assign it as the
	// custom handler for 404 Not Found response. 
	// A custom handler for 405 Method Not Allowed can also be set by router.MethodNotAllowed
	// in the same way.
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		app.notFound(w)
	})

	// Update the pattern for the route for the static files.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath",  http.StripPrefix("/static", fileServer))

	// Create the methods using the appropriate methods, patterns and handlers.
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standard := alice.New(app.recoverPanic, app.loqRequest, secureHeaders)

	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(router)
}
