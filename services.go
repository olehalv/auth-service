package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
)

func auth(w http.ResponseWriter, r *http.Request) {
	var oneHourFromNow = time.Now().Add(time.Hour)
	var user LoginRequest

	err := decodeJsonBody(r.Body, &user)

	if err != nil {
		returnHttpStatus(w, r, http.StatusBadRequest, "Bad JSON", err)
		return
	}

	if !authUser(user) {
		returnHttpStatus(w, r, http.StatusUnauthorized, "Wrong username or password", err)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "auth-service",
		"sub": user.Email,
		"exp": oneHourFromNow.Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		returnHttpStatus(w, r, http.StatusInternalServerError, "Could not sign token", err)
		return
	}

	err = setLastLoggedIn(user.Email)

	if err != nil {
		returnHttpStatus(w, r, http.StatusInternalServerError, "Internal server error", err)
		return
	}

	cookie := http.Cookie{
		Name:     cookieName,
		Value:    tokenString,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Expires:  oneHourFromNow,
	}

	http.SetCookie(w, &cookie)

	res, _ := json.Marshal(AuthResponse{
		Token: tokenString,
	})

	_, _ = w.Write(res)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	var reqBodyToken string
	err := decodeJsonBody(r.Body, &reqBodyToken)
	cookieToken, err1 := r.Cookie(cookieName)

	if err != nil && err1 != nil {
		returnHttpStatus(w, r, http.StatusBadRequest, "No token provided (cookie or body)", err)
		return
	}

	var token *jwt.Token

	if &cookieToken != nil {
		token, err = jwt.Parse(
			cookieToken.Value,
			func(token *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			},
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		)

		if err != nil {
			returnHttpStatus(w, r, http.StatusBadRequest, "Bad authentication token", err)
			return
		}
	} else {
		token, err = jwt.Parse(
			reqBodyToken,
			func(token *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			},
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		)

		if err != nil {
			returnHttpStatus(w, r, http.StatusBadRequest, "Bad authentication token", err)
			return
		}
	}

	email, err := token.Claims.GetSubject()

	if err != nil {
		returnHttpStatus(w, r, http.StatusBadRequest, "Bad authentication token", err)
		return
	}

	iss, err := token.Claims.GetIssuer()

	if err != nil {
		returnHttpStatus(w, r, http.StatusBadRequest, "Bad authentication token", err)
		return
	}

	if iss != "auth-service" {
		returnHttpStatus(w, r, http.StatusBadRequest, "Bad authentication token", nil)
		return
	}

	user, err := getUserDetails(email)

	if err != nil {
		returnHttpStatus(w, r, http.StatusBadRequest, "Bad authentication token", err)
		return
	}

	res, err := json.Marshal(user)

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
	var user RegisterRequest

	err := decodeJsonBody(r.Body, &user)

	if err != nil {
		returnHttpStatus(w, r, http.StatusBadRequest, "Bad JSON", err)
		return
	}

	if user.Email == "" || user.Pass == "" || user.InvCode == "" {
		returnHttpStatus(w, r, http.StatusBadRequest, "Please provide email, pass and invitation code", nil)
		return
	}

	if user.InvCode != os.Getenv("INV_CODE") {
		returnHttpStatus(w, r, http.StatusBadRequest, "Wrong invitation code", nil)
		return
	}

	if userExists(user.Email) {
		returnHttpStatus(w, r, http.StatusConflict, "User with email already exists", nil)
		return
	}

	pass, err := hashString(user.Pass)

	if err != nil {
		returnHttpStatus(w, r, http.StatusInternalServerError, "Internal server error", err)
		return
	}

	row, err := psql.Query(
		context.Background(),
		"insert into users (email, password, created, lastLoggedIn) values ($1, $2, $3, $3)",
		user.Email,
		pass,
		time.Now(),
	)

	if err != nil {
		returnHttpStatus(w, r, http.StatusInternalServerError, "Could not create user, try again later", err)
		return
	}

	row.Close()

	returnHttpStatus(w, r, http.StatusCreated, fmt.Sprintf("User created: %s", user.Email), nil)
}
