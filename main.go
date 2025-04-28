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

	initPsql()

	http.HandleFunc("/api/auth", preventSpam(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			auth(w, r)
			break
		default:
			returnHttpStatus(w, r, http.StatusMethodNotAllowed, "Method not allowed", nil)
			break
		}
	}))

	http.HandleFunc("/api/user", preventSpam(func(w http.ResponseWriter, r *http.Request) {
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

	err = http.ListenAndServe(host, nil)
	if err != nil {
		LogFatal(err)
	}
}
