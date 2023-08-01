package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Define an application struct to hold application-wide dependencies.
type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize a new instance of application struct containing dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}

	// Initialize a new http.Server struct. Set the Addr and Handler fields so that the
	// server uses the same network address routes as before. Set the ErrorLog field
	// so that the server now uses the custom errorLog logger.
	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	// Call ListenAndServe() method on our new http.Server struct.
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}