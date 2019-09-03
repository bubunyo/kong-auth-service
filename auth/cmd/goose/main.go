package main

import (
	"flag"
	"github.com/bubunyo/para-services/auth/db"
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

	conn := db.New()

	var arguments []string
	if len(args) > 0 {
		arguments = append(arguments, args[1:]...)
	}

	command := args[0]

	if err := goose.Run(command, conn.DB, *dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
