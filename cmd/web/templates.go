package main

import (
	"html/template"
	"path/filepath"

	"snippetbox.sangdennis.com/internal/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// At the moment it only contains one field, but we'll add more
// to it as the build progresses.
type templateData struct {
	CurrentYear int
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as a new cache
	cache := map[string]*template.Template{}

	// Use filepath.Glob() to get a slice of all filepaths that match
	// the pattern "./ui/html/pages/*.html". It gives us a slice of all
	// the filepaths for our application 'page' templates like:
	// [ui/html/pages/home.html ui/html/pages/view.html]
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	// loop through the page filepaths one by one
	for _, page := range pages {
		// Extract the file name (e.g "home.html") from the full filepath
		// and assign it to the name variable.
		name := filepath.Base(page)

		// Parse the base template into a template set
		ts, err := template.ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() *on this template set* to add any partials.
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() *on this template set* to add the page template.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map, using the name of the page
		// (like "home.html") as the key.
		cache[name] = ts
	}

	// Return the map
	return cache, nil
}
