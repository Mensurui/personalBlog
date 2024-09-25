package main

import (
	"bytes"
	"log"
	"net/http"
)

func (app *application) render(w http.ResponseWriter, status int, page string, data *Blog) {
	template, ok := app.templateCache[page]
	if !ok {
		log.Printf("The template %s doesn't exist", page)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Rendering template: %s", page)

	buf := new(bytes.Buffer)

	// Execute the template into a buffer
	err := template.ExecuteTemplate(buf, "base.tmpl", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write status code and content after successful template execution
	w.WriteHeader(status)
	buf.WriteTo(w)
}
