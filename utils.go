package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func clearConsole() {
	clears := []string{"clear", "cls"}
	for _, s := range clears {
		cmd := exec.Command(s)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	}
}

var userExists = func(user LoginRequest, r *http.Request) bool {
	if user == rootUser {
		return true
	}

	var email string
	var pass string

	err := psql.QueryRow(
		context.Background(),
		fmt.Sprintf(
			"select * from users u where u.email = '%s'",
			user.Email,
		),
	).Scan(&email, &pass)

	if err != nil {
		Log(fmt.Sprintf("%s: %s", getIp(r), err))
		return false
	}

	if &email == nil || &pass == nil {
		return false
	}

	return true
}

func decodeJSON(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(&v)
}

func returnHttpStatus(w http.ResponseWriter, r *http.Request, status int, message string, error error) {
	w.WriteHeader(status)
	res, _ := json.Marshal(HttpStatusResponse{Message: message})
	_, _ = w.Write(res)
	Log(fmt.Sprintf("%s: %s -> INTERNAL ERROR: %s", getIp(r), message, error))
}

func getIp(r *http.Request) string {
	ip := r.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}

func Log(any any) {
	logger.Println(any)
}

func LogFatal(any any) {
	logger.Fatalln(any)
}
