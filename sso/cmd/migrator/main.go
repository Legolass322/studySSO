package main

import (
	"errors"
	"flag"
	"fmt"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

func main() {
	var connstr, migrationPath, migrationTable string

	flag.StringVar(&connstr, "conn", "", "")
	flag.StringVar(&migrationPath, "mpath", "", "")
	flag.StringVar(&migrationTable, "mtable", "migrations", "")
	flag.Parse()

	if connstr == "" {
		panic("No connection string")
	}

	if migrationPath == "" {
		panic("No migration path folder")
	}

	m, err := migrate.New(
		"file://" + migrationPath,
		fmt.Sprintf("%s?x-migrations-table=%s&sslmode=disable", connstr, migrationTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No changes to apply")
			return
		}

		panic(err)
	}
}