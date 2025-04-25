package main

import (
	"fmt"
	"net/http"
)

func redirectToIndex(w http.ResponseWriter, r *http.Request) {
	Log(fmt.Sprintf("%s: %s", getIp(r), "Redirecting to login"))
	http.Redirect(w, r, "/index.html", http.StatusSeeOther)
}

func isMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Log(fmt.Sprintf("%s: %s", getIp(r), "Method not allowed"))
		return false
	}
	return true
}

func Err(w http.ResponseWriter, r *http.Request, err error, status int) bool {
	if err != nil {
		w.WriteHeader(status)
		Log(fmt.Sprintf("%s: %s", getIp(r), err))
		return true
	}
	return false
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
