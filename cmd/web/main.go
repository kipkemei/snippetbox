package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"snippetbox.sangdennis.com/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

// Define an application struct to hold application-wide dependencies.
type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *models.SnippetModel
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:Naingia12@/snippetbox?parseTime=true", "MySQL data soure name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// pass to openDB the DSN command-line flag
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initialize a new instance of application struct containing dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &models.SnippetModel{DB: db},
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
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// openDB() wraps sql.Open() and returns a sql.DB connection pool for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}