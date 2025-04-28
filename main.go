package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"net/http"
	"os"
)

func main() {
	clearConsole()

	err := godotenv.Load()
	if err != nil {
		LogFatal(err)
	}

	initPsql()

	mux := http.NewServeMux()

	mux.HandleFunc("/api/auth", preventSpam(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			auth(w, r)
			break
		default:
			returnHttpStatus(w, r, http.StatusMethodNotAllowed, "Method not allowed", nil)
			break
		}
	}))

	mux.HandleFunc("/api/user", preventSpam(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getUser(w, r)
			break
		case "POST":
			createUser(w, r)
			break
		default:
			returnHttpStatus(w, r, http.StatusMethodNotAllowed, "Method not allowed", nil)
			break
		}

	}))

	host := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))

	Log(fmt.Sprintf("URL: %s", fmt.Sprintf("http://%s", host)))

	err = http.ListenAndServe(host, cors.AllowAll().Handler(mux))
	if err != nil {
		LogFatal(err)
	}
}
