package config

import (
	"database/sql"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// Connect to Postgres database
func InitializeDB() (db *bun.DB) {
	localEnv, err := GetEnvMap()
	if err != nil {
		log.Panic("Failed to load .env file")
	}

	dsn := localEnv["DB_POSTGRES"]
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db = bun.NewDB(pgdb, pgdialect.New())
	return
}
