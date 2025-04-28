package main

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func setLastLoggedIn(email string) error {
	_, err := psql.Exec(
		context.Background(),
		"update users set lastLoggedIn = $1 where email = $2",
		time.Now(),
		email,
	)

	if err != nil {
		return err
	}

	return nil
}

func authUser(user LoginRequest) bool {
	var pass string

	err := psql.QueryRow(
		context.Background(),
		"select password from users where email = $1",
		user.Email,
	).Scan(&pass)

	if &pass == nil || err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(user.Pass))

	if err != nil {
		return false
	}

	return true
}

func userExists(email string) bool {
	var _email string

	err := psql.QueryRow(
		context.Background(),
		"select email from users where email = $1",
		email,
	).Scan(&_email)

	if &_email == nil || err != nil {
		return false
	}

	return true
}

func getUserDetails(email string) (UserDetailsResponse, error) {
	var _email string
	var created, lastLoggedIn time.Time

	err := psql.QueryRow(
		context.Background(),
		"select email, created, lastLoggedIn from users where email = $1",
		email,
	).Scan(&_email, &created, &lastLoggedIn)

	if err != nil || _email == "" {
		return UserDetailsResponse{}, err
	}

	return UserDetailsResponse{
		Email:        _email,
		Created:      created.Format(time.DateTime),
		LastLoggedIn: lastLoggedIn.Format(time.DateTime),
	}, nil
}
