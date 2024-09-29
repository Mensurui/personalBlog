package main

import (
	"fmt"
	"github.com/Mensurui/personalBlog.git/internals/models"
	"github.com/Mensurui/personalBlog.git/ui"
	"html/template"
	"io/fs"
	"log"
	"path/filepath"
)

type TemplateCache map[string]*template.Template

type templateData struct {
	Article         *models.Article
	Articles        []*models.Article
	User            *models.Users
	IsAuthenticated bool
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")

	if err != nil {
		log.Print(err)
	}

	if len(pages) == 0 {
		log.Print("No templates found in the specified path")
		return nil, fmt.Errorf("no templates found")
	}

	log.Printf("Templates found: %v", pages)

	for _, page := range pages {

		// Get the name of the template
		name := filepath.Base(page)

		// Create a slice of patterns to include the base and partials
		patterns := []string{
			"html/base.tmpl",       // Base template
			"html/partials/*.tmpl", // Partials (e.g., headers, footers)
			page,                   // Specific page template
		}

		// Use ParseFS() to load templates from the embedded filesystem
		ts, err := template.New(name).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// Add the parsed template set to the cache
		cache[name] = ts

	}

	return cache, nil
}
