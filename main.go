package main

import (
	"log"
	"net/http"
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
	w.Write([]byte("Display a specific snippet..."))
}

// Add a snippetCreate handler function
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	// Use r.Method to check whether the request is using POST or not.
	if r.Method != "POST" {
		// If it's not, use w.WriteHeader() method to send a 405 status code 
		// and w.Write() method to write a "Method Not Allowed" response body.
		// Then return from the function to prevent execution of subsequent code.
		w.Header().Set("Allow", "POST")
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
		return
	}

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