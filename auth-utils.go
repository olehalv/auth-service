package main

import "os"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"pass"`
}

type AuthResponse struct {
	ReturnUrl string `json:"returnUrl"`
	Token     string `json:"token"`
}

type UserResponse struct {
	Email string `json:"email"`
}

var users []LoginRequest

var userExists = func(user LoginRequest) bool {
	users = append(users, LoginRequest{
		Email:    os.Getenv("ADMIN_EMAIL"),
		Password: os.Getenv("ADMIN_PASS")},
	)

	for _, tU := range users {
		if user.Email == tU.Email && user.Password == tU.Password {
			return true
		}
	}
	return false
}
