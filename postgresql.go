package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"os"
)

var psql *pgx.Conn

func initializePostgres() {
	var err error
	psql, err = pgx.Connect(context.Background(), os.Getenv("PSQL_URL"))

	if err != nil {
		LogFatal(err)
	}

	Log("Connected to database")
}
