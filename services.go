package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
	"time"
)

func auth(w http.ResponseWriter, r *http.Request) {
	var user LoginRequest

	err := decodeJSON(r.Body, &user)

	if err != nil {
		returnHttpStatus(w, r, http.StatusBadRequest, "Bad JSON", err)
		return
	}

	if !userExists(user, r) {
		returnHttpStatus(w, r, http.StatusUnauthorized, "Wrong username or password", nil)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": os.Getenv("JWT_ISSUER"),
		"sub": user.Email,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		returnHttpStatus(w, r, http.StatusInternalServerError, "Could not sign token", err)
		return
	}

	res, err := json.Marshal(AuthResponse{
		ReturnUrl: r.Header.Get("Referer"),
		Token:     tokenString,
	})

	if err != nil {
		returnHttpStatus(w, r, http.StatusInternalServerError, "Could not return token", err)
		return
	}

	_, err = w.Write(res)

	if err != nil {
		returnHttpStatus(w, r, http.StatusInternalServerError, "Could not return token", err)
		return
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		returnHttpStatus(w, r, http.StatusBadRequest, "Bad authentication token", nil)
		return
	}

	authHeader = parts[1]

	token, err := jwt.Parse(
		authHeader,
		func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	if err != nil {
		returnHttpStatus(w, r, http.StatusBadRequest, "Bad authentication token", err)
		return
	}

	sub, err := token.Claims.GetSubject()

	if err != nil {
		returnHttpStatus(w, r, http.StatusBadRequest, "Bad authentication token", err)
		return
	}

	iss, err := token.Claims.GetIssuer()

	if err != nil {
		returnHttpStatus(w, r, http.StatusBadRequest, "Bad authentication token", err)
		return
	}

	if iss != os.Getenv("JWT_ISSUER") {
		returnHttpStatus(w, r, http.StatusBadRequest, "Bad authentication token", nil)
		return
	}

	res, err := json.Marshal(UserResponse{Email: sub})

	if err != nil {
		returnHttpStatus(w, r, http.StatusInternalServerError, "Internal server error", err)
		return
	}

	_, err = w.Write(res)

	if err != nil {
		returnHttpStatus(w, r, http.StatusInternalServerError, "Internal server error", err)
		return
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user LoginRequest

	err := decodeJSON(r.Body, &user)

	if err != nil {
		returnHttpStatus(w, r, http.StatusBadRequest, "Bad JSON", err)
		return
	}

	if user.Email == "" || user.Password == "" {
		returnHttpStatus(w, r, http.StatusBadRequest, "Please provide email:pass", nil)
		return
	}

	if userExists(user, r) {
		returnHttpStatus(w, r, http.StatusConflict, "User with email already exists", nil)
		return
	}

	row, err := psql.Query(
		context.Background(),
		fmt.Sprintf(
			"insert into users (email, password) values ('%s', '%s')",
			user.Email,
			user.Password,
		),
	)

	if err != nil {
		returnHttpStatus(w, r, http.StatusInternalServerError, "Could not create user, try again later", err)
		return
	}

	err = row.Err()

	if err != nil {
		returnHttpStatus(w, r, http.StatusInternalServerError, "Internal server error", err)
		return
	}

	row.Close()

	returnHttpStatus(w, r, http.StatusCreated, fmt.Sprintf("User created: %s", user.Email), nil)
}
