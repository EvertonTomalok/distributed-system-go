package postgres

import (
	"context"
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
	log "github.com/sirupsen/logrus"
)

var (
	Db      *sql.DB
	Started bool = false
)

func Init(ctx context.Context, host string) *sql.DB {
	db, err := sql.Open("nrpostgres", host)

	if err != nil {
		log.Panicf("Connecting to database: %+v", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: "distribute-system-schema-migrations",
	})
	if err != nil {
		log.Panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./db/migrations/postgres",
		"distributed-system",
		driver,
	)
	if err != nil {
		log.Panicf("Error connecting migrator %+v", err)
	}
	if err := m.Up(); err != nil {
		if string(err.Error()) != "no change" {
			log.Panicf("Error making the migration -> %+v", err)
		}
	}
	return db
}

func Check(database *sql.DB) error {
	if err := database.Ping(); err != nil {
		return err
	}
	Started = true
	return nil
}
