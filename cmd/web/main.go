package main

import (
	"context"
	"flag"
	"github.com/Mensurui/personalBlog.git/internals/models"
	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"time"
)

type application struct {
	articleModel   *models.ArticleModel
	userModel      *models.UserDB
	templateCache  TemplateCache
	sessionManager *scs.SessionManager
	form           *form.Decoder
}

func main() {
	dns := flag.String("dns", "postgres://web:pass@localhost:5432/personalblog", "Postgres Datasource Name")

	flag.Parse()
	db, err := openDB(*dns)

	if err != nil {
		log.Println("Connection error: %s", err)
		return
	}
	articleDB := &models.ArticleModel{DB: db}
	userDB := &models.UserDB{DB: db}
	log.Println("Connection working")
	defer db.Close()
	templateCache, err := newTemplateCache()
	if err != nil {
		log.Printf("error: %v", err)
	}
	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		templateCache:  templateCache,
		articleModel:   articleDB,
		userModel:      userDB,
		sessionManager: sessionManager,
		form:           formDecoder,
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
