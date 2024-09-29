package models

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Article struct {
	ID      int
	Title   string
	Content string
}

type ArticleModel struct {
	DB *pgxpool.Pool
}

func (m *ArticleModel) Get(id int) (*Article, error) {

	var title, content string

	stmt := `SELECT name, content
			FROM articles
			WHERE id = $1;`

	err := m.DB.QueryRow(context.Background(), stmt, id).Scan(&title, &content)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	article := &Article{
		Title:   title,
		Content: content,
	}

	return article, nil
}

func (m *ArticleModel) Add(title string, content string) (int64, error) {
	stmt := `INSERT INTO articles(name, content)
			VALUES($1, $2)
			RETURNING id;`

	var id int64

	err := m.DB.QueryRow(context.Background(), stmt, title, content).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *ArticleModel) Latest() ([]*Article, error) {

	stmt := `SELECT id, name, content
			FROM articles`

	rows, err := m.DB.Query(context.Background(), stmt)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	articles := []*Article{}

	for rows.Next() {
		article := &Article{}

		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
		)

		if err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}
