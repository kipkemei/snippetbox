package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the flag 
	// will be stored in the 'addr' variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Use flag.Parse() to read the command-line flag value and assign it to 'addr'.
	// It should be called before the variable is used, else it'll contain the default value.
	// The application is terminated if any errors are encountered during parsing.
	flag.Parse()

	mux := http.NewServeMux()

	// Create a file server which serves files out of "./ui/static" directory.
	// The path given to http.Dir() is relative to the project root.
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})

	// Use mux.Handle() to register the file server as the handler for all URL
	// paths that start with "/static/". For matching parts, strip the prefix
	// "/static" before the request reaches the file server.
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Register the other application routes as normal.
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
    f, err := nfs.fs.Open(path)
    if err != nil {
        return nil, err
    }

    s, err := f.Stat()
	if err != nil {
        return nil, err
    }

    if s.IsDir() {
        index := filepath.Join(path, "index.html")
        if _, err := nfs.fs.Open(index); err != nil {
            closeErr := f.Close()
            if closeErr != nil {
                return nil, closeErr
            }

            return nil, err
        }
    }

    return f, nil
}    