package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"os"
)

func initPsql() {
	var err error
	psql, err = pgx.Connect(context.Background(), os.Getenv("PSQL_URL"))

	if err != nil {
		LogFatal(err)
	}

	Log("Connected to database")
}
