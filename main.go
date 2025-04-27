package main

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

var rootUser = LoginRequest{
	Email:    "root",
	Password: "",
}

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

	h := sha512.New()
	h.Write([]byte(uuid.NewString()))
	pass := base64.URLEncoding.EncodeToString(h.Sum(nil))

	rootUser.Password = pass
	Log(fmt.Sprintf("%s:%s", rootUser.Email, rootUser.Password))

	err = http.ListenAndServe(host, nil)
	if err != nil {
		LogFatal(err)
	}
}
