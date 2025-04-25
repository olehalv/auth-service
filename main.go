package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

func main() {
	clearConsole()

	err := godotenv.Load()
	if err != nil {
		LogFatal(err)
	}

	initializePostgres()
	staticHandler()
	authRouter()

	host := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))

	Log(fmt.Sprintf("Starting server with %s", host))

	err = http.ListenAndServe(host, nil)
	if err != nil {
		LogFatal(err)
	}
}
