package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

type Connection struct{ DB *sql.DB }

func New() *Connection {
	var (
		host     = os.Getenv("DB_HOST")
		port     = 5432
		dbname   = os.Getenv("DB_DATABASE")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("auth-app: failed to open DB: %v\n", err)
	}
	return &Connection{db}
}
