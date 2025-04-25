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
		Log(err)
		return
	}

	staticHandler()
	authRouter()

	port := os.Getenv("PORT")
	fmt.Println(fmt.Sprintf("Starting server on port %s", port))
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		Log(err)
	}
}
