package main

import (
	"context"
	"fmt"
	"net/http"
)

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

var userExists = func(user LoginRequest, r *http.Request) bool {
	if user == rootUser {
		return true
	}

	var email string
	var pass string

	err := psql.QueryRow(
		context.Background(),
		fmt.Sprintf(
			"select * from users u where u.email = '%s' and u.password = '%s'",
			user.Email,
			user.Password,
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
