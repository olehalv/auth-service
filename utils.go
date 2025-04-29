package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
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

func decodeJsonBody(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(&v)
}

func returnHttpStatus(w http.ResponseWriter, r *http.Request, status int, message string, error error) {
	w.WriteHeader(status)
	res, _ := json.Marshal(HttpStatusResponse{Message: message})
	_, _ = w.Write(res)
	if error != nil {
		Log(fmt.Sprintf("%s: %s -> INTERNAL ERROR: %s", getIp(r), message, error))
	} else {
		Log(fmt.Sprintf("%s: %s", getIp(r), message))
	}
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

func hashString(value string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(pass), nil
}
