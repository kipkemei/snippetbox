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
func home(w http.ResponseWriter, r *http.Request)  {
	w.Write([]byte("Hello from Snippetbox"))	
}

func main() {
	// Use http.NewServeMux() to initialize a new servemux, then register
	// the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	// Use http.ListenAndServe() to start a new web server. 
	// Pass in two parameters: the TCP network address to listen on (":4000")
	// and the servemux. If http.ListenAndServe() returns an error, 
	// use log.Fatal() to log the error message and exit.
	// Note: Errors returned by http.ListenAndServe() are always non-nil.
	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}