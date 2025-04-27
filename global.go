package main

import (
	"github.com/jackc/pgx/v5"
	"log"
	"os"
)

var rootUser = LoginRequest{
	Email:    "root",
	Password: "",
}
var logger = log.New(os.Stdout, "-> ", log.Ldate|log.Ltime)
var requests []Request
var psql *pgx.Conn
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
