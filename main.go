package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
// home() takes two parameters: http.ResponseWriter which provides methods
// for assembling a HTTP response and sending it to a user and the *http.Request
// parameter which is a pointer to a struct which holds information about the
// current request.
func home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/".
	// If not, use http.NotFound() to send a 404 response to client and then return.
	// If we don't return, the handler would keep executing and write the byte slice.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from Snippetbox"))
}

// Add a snippetView handler function
func snippetView(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the id parameter from the query string and convert it
	// to an integer using strconv.Atoi(). If it can't be converted to an integer, 
	// or the value is less than 1, return a 404 page not found response.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Use fmt.Fprintf() to interpolate id value with our response 
	// and write it to http.ResponseWriter
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// Add a snippetCreate handler function
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	// Use r.Method to check whether the request is using POST or not.
	if r.Method != http.MethodPost {
		// Use http.Error() to send a 405 status code "Method Not Allowed" response body.
		// Then return from the function to prevent execution of subsequent code.
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Create a new snippet..."))
}

func main() {
	// Use http.NewServeMux() to initialize a new servemux, then register
	// the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// Use http.ListenAndServe() to start a new web server.
	// Pass in two parameters: the TCP network address to listen on (":4000")
	// and the servemux. If http.ListenAndServe() returns an error,
	// use log.Fatal() to log the error message and exit.
	// Note: Errors returned by http.ListenAndServe() are always non-nil.
	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
