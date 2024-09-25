package main

import (
	"context"
	"flag"
	"github.com/jackc/pgx/v5/pgxpool"
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
	dns := flag.String("dns", "postgres://web:pass@localhost:5432/personalblog", "Postgres Datasource Name")

	flag.Parse()
	db, err := openDB(*dns)

	if err != nil {
		log.Println("Connection error: %v", err)
		return
	}
	log.Println("Connection working")
	defer db.Close()
	templateCache, err := newTemplateCache()
	if err != nil {
		log.Printf("error: %v", err)
	}

	app := &application{
		templateCache: templateCache,
	}
	err = http.ListenAndServe(":8080", app.routes())
	if err != nil {
		log.Fatal(err)
	}
}

func openDB(dns string) (*pgxpool.Pool, error) {
	dbPool, err := pgxpool.New(context.Background(), dns)
	if err != nil {
		return nil, err
	}
	if err = dbPool.Ping(context.Background()); err != nil {
		dbPool.Close()
		return nil, err
	}
	return dbPool, nil
}
