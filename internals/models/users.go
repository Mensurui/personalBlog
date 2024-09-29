package models

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type Users struct {
	email    string
	username string
	password string
}

type UserDB struct {
	DB *pgxpool.Pool
}

func (m *UserDB) Create(email string, username string, password string) error {

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		log.Printf("The error is: %s", err)
		return err
	}

	stmt := `INSERT INTO users(email, username, password)
			VALUES($1, $2, $3)
			RETURNING id`

	_, err = m.DB.Exec(context.Background(), stmt, email, username, string(hashed))

	if err != nil {
		log.Printf("The error is:  %s", err)
		return err
	}

	return nil
}

func (m *UserDB) Login(email string, password string) (int, error) {
	var id int
	var hashed []byte

	stmt := `SELECT id, password FROM users WHERE email=$1`
	err := m.DB.QueryRow(context.Background(), stmt, email).Scan(&id, &hashed)
	if err != nil {
		log.Printf("The error is %s", err)
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashed, []byte(password))

	if err != nil {
		log.Printf("The error is %s", err)
		return 0, err
	}

	return id, nil
}
