package models

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Article struct {
	Title   string
	Content string
}

type ArticleModel struct {
	DB *pgxpool.Pool
}

func (m *ArticleModel) Get(id int) (*Article, error) {
	if id != 1 {
		log.Println("No article with that id")
		return nil, nil
	}
	article := &Article{
		Title:   "Hello",
		Content: "There",
	}

	return article, nil
}
