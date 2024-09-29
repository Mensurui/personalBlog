package main

import (
	"bytes"
	"log"
	"net/http"
)

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
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

func (app *application) newTemplateData(r *http.Request) *templateData {
	userID, ok := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	return &templateData{
		IsAuthenticated: ok && userID > 0,
	}
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	log.Println("Successfully Authenticated")
	log.Printf("User authenticated: %v", isAuthenticated)

	return isAuthenticated
}
