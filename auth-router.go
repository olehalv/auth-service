package main

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
	"time"
)

func authRouter() {
	secret := []byte(os.Getenv("JWT_SECRET"))

	http.HandleFunc("/api/auth", preventSpam(
		func(w http.ResponseWriter, r *http.Request) {
			isErr := func(err error, returnStatus int) bool {
				return Err(w, r, err, returnStatus)
			}

			if !isMethod(w, r, "POST") {
				return
			}

			var user LoginRequest
			err := json.NewDecoder(r.Body).Decode(&user)
			if isErr(err, 400) {
				return
			}

			if !userExists(user, r) {
				Err(w, r, errors.New("user doesnt exist"), 401)
				return
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"iss": os.Getenv("JWT_ISSUER"),
				"sub": user.Email,
				"exp": time.Now().Add(time.Hour).Unix(),
			})

			tokenString, err := token.SignedString(secret)
			if isErr(err, 500) {
				return
			}

			res, err := json.Marshal(AuthResponse{
				ReturnUrl: r.Header.Get("Referer"),
				Token:     tokenString,
			})
			if isErr(err, 500) {
				return
			}

			w.Header().Add("Authorization", "Bearer "+tokenString)
			_, err = w.Write(res)
			if isErr(err, 500) {
				return
			}
		}))

	http.HandleFunc("/api/user", preventSpam(
		func(w http.ResponseWriter, r *http.Request) {
			isErr := func(err error, returnStatus int) bool {
				return Err(w, r, err, returnStatus)
			}

			if !isMethod(w, r, "GET") {
				return
			}

			authHeader := r.Header.Get("Authorization")
			parts := strings.Split(authHeader, " ")

			if len(parts) != 2 || parts[0] != "Bearer" {
				redirectToIndex(w, r)
				return
			}

			authHeader = parts[1]

			token, err := jwt.Parse(
				authHeader,
				func(token *jwt.Token) (interface{}, error) {
					return secret, nil
				},
				jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
			)
			if isErr(err, 401) {
				return
			}

			sub, err := token.Claims.GetSubject()
			if isErr(err, 401) {
				return
			}
			iss, err := token.Claims.GetIssuer()
			if isErr(err, 401) {
				return
			}
			exp, err := token.Claims.GetExpirationTime()
			if isErr(err, 401) {
				return
			}

			if iss != os.Getenv("JWT_ISSUER") || exp.Before(time.Now()) {
				redirectToIndex(w, r)
				return
			}

			res, err := json.Marshal(UserResponse{Email: sub})
			if isErr(err, 401) {
				return
			}

			_, err = w.Write(res)
			if isErr(err, 500) {
				return
			}
		}))
}
