package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
)

var db *sql.DB

func Init(ctx context.Context, host string) {
	db, err := sql.Open("nrpostgres", host)
	if err != nil {
		panic(fmt.Sprintf("Connecting to database: %+v", err))
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: "distribute-system-schema-migrations",
	})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./db/migrations/postgres",
		"distributed-system",
		driver,
	)
	if err != nil {
		panic(fmt.Sprintf("Error connecting migrator %+v", err))
	}
	if err := m.Up(); err != nil {
		panic(fmt.Sprintf("Error making the migration %+v", err))
	}

}
