package main

import (
	"log"
	"net/http"
)

type Blog struct {
	Title   string
	Content string
}

type application struct {
	templateCache TemplateCache
}

func main() {
	templateCache, err := newTemplateCache()
	if err != nil {
		log.Println("error")
	}

	app := &application{
		templateCache: templateCache,
	}
	err = http.ListenAndServe(":8080", app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
