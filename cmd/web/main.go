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
}

func main() {
	app := &application{}
	err := http.ListenAndServe(":8080", app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
