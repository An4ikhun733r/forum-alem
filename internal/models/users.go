package models

import (
	"forum/internal/validator"
	"regexp"
	"time"
)

var EmailRX = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserSignupForm struct {
	Name string
	Password string
	Email    string
	validator.Validator
}

type UserLoginForm struct {
	Name string
	Password string
	validator.Validator
}