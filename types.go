package main

import "time"

type Request struct {
	Ip   string
	Time time.Time
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"pass"`
}

type AuthResponse struct {
	ReturnUrl string `json:"returnUrl"`
	Token     string `json:"token"`
}

type UserDetailsResponse struct {
	Email        string `json:"email"`
	Created      string `json:"created"`
	LastLoggedIn string `json:"lastLoggedIn"`
}

type HttpStatusResponse struct {
	Message string `json:"message"`
}
