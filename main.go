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
	authRouter()

	host := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))

	Log(fmt.Sprintf("Starting server with %s", fmt.Sprintf("http://%s", host)))

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
