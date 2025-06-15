package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

func main() {
	dbHost := flag.String("db-host", "sso-db", "Database host")
	dbPort := flag.String("db-port", "5432", "Database port")
	dbUser := flag.String("db-user", "admin", "Database user")
	dbPassword := flag.String("db-password", "", "Database password")
	dbName := flag.String("db-name", "sso_db", "Database name")
	dbSSLMode := flag.String("db-sslmode", "disable", "SSL mode (disable/require)")

	migrationsPath := flag.String("migrations-path", "", "Path to migrations folder")
	migrationsTable := flag.String("migrations-table", "migrations", "Migrations table name")

	command := flag.String("command", "up", "Migration command (up/down)")

	flag.Parse()

	if *migrationsPath == "" {
		panic("migrations-path is required")
	}
	if *dbPassword == "" {
		if pass := os.Getenv("DB_PASSWORD"); pass != "" {
			*dbPassword = pass
		} else {
			panic("db-password is required")
		}
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&x-migrations-table=%s",
		*dbUser,
		*dbPassword,
		*dbHost,
		*dbPort,
		*dbName,
		*dbSSLMode,
		*migrationsTable)

	m, err := migrate.New(
		"file://"+*migrationsPath,
		dsn,
	)
	if err != nil {
		panic(err)
	}

	switch *command {
	case "up":
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no migrations to apply")

				return
			}

			panic(err)
		}
	case "down":
		if err := m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no migrations to apply")
			}

			panic(err)
		}
	default:
		panic("unknow command")
	}

}
