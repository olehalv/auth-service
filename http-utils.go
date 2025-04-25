package main

import (
	"fmt"
	"net/http"
)

func redirectToIndex(w http.ResponseWriter, r *http.Request) {
	Log(fmt.Sprintf("%s: %s", r.Host, "Redirecting to login"))
	http.Redirect(w, r, "/index.html", http.StatusSeeOther)
}

func isMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Log(fmt.Sprintf("%s: %s", r.Host, "Method not allowed"))
		return false
	}
	return true
}

func Err(w http.ResponseWriter, r *http.Request, err error, status int) bool {
	if err != nil {
		w.WriteHeader(status)
		Log(fmt.Sprintf("%s: %s", r.Host, err))
		return true
	}
	return false
}
