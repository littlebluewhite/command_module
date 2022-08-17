package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/pressly/goose/v3"
	"log"
	"new_command/app/migrate"
	"os"
)

const dialect = "mysql"

var (
	flags = flag.NewFlagSet("migration", flag.ExitOnError)
	dir   = flags.String("dir", "./migration", "directory with migration files")
)

func main() {
	flags.Usage = usage
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatalf("migration run: %v", err)
	}

	args := flags.Args()
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}

	fmt.Println(args)

	command := args[0]

	switch command {
	case "create":
		if err := goose.Run("create", nil, *dir, args[1:]...); err != nil {
			log.Fatalf("migration run: %v", err)
		}
		return
	case "fix":
		if err := goose.Run("fix", nil, *dir); err != nil {
			log.Fatalf("migration run: %v", err)
		}
		return
	}

	// initialize data sources
	appDb, err := migrate.NewMigrateDB()
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer func(appDb *sql.DB) {
		if err := appDb.Close(); err != nil {
			log.Fatalf("migration close: %v", err)
		}
	}(appDb)

	if err := goose.SetDialect(dialect); err != nil {
		log.Fatal(err)
	}

	if err := goose.Run(command, appDb, *dir, args[1:]...); err != nil {
		log.Fatalf("migration run: %v", err)
	}
}

func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

var (
	usagePrefix = `Usage: migration [OPTIONS] COMMAND
Examples:
    migration status
Options:`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
    fix                  Apply sequential ordering to migrations`
)