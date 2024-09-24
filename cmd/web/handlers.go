package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	files := []string{
		filepath.Join("ui", "html", "base.tmpl"),
		filepath.Join("ui", "html", "pages", "home.tmpl"),
		filepath.Join("ui", "html", "partials", "nav.tmpl"),
	}

	ts, err := template.ParseFiles(files...)

	if err != nil {
		log.Fatal(err)
		return
	}
	title := "First Blog"
	content := "My first content"
	data := &Blog{
		Title:   title,
		Content: content,
	}

	err = ts.Execute(w, data)

	if err != nil {
		log.Fatal(err)
		return
	}

}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	files := []string{
		filepath.Join("ui", "html", "base.tmpl"),
		filepath.Join("ui", "html", "pages", "about.tmpl"),
		filepath.Join("ui", "html", "partials", "nav.tmpl"),
	}

	ts, err := template.ParseFiles(files...)

	if err != nil {
		log.Fatal(err)
		return
	}

	err = ts.Execute(w, nil)
}
