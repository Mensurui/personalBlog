package main

import (
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	title := "Something"
	content := "Some Content"
	data := &Blog{
		Title:   title,
		Content: content,
	}
	app.render(w, http.StatusOK, "home.tmpl", data)

}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "about.tmpl", nil)
}

func (app *application) articles(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "articles.tmpl", nil)
}
