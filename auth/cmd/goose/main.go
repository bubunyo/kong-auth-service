package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/bubunyo/para-services/auth/db/migrations"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", "./db/migrations", "directory with migration files")
)

func main() {
	_ = flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) < 1 {
		flags.Usage()
		return
	}

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
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	var arguments []string
	if len(args) > 0 {
		arguments = append(arguments, args[1:]...)
	}

	command := args[0]

	if err := goose.Run(command, db, *dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
