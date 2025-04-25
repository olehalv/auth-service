package main

import "net/http"

func staticHandler() {
	fs := http.FileServer(http.Dir("public/"))
	http.Handle("/", http.StripPrefix("/", fs))
}
