package models

import (
	"errors"
)

var ErrNoRecord = errors.New("models: no matching record found")
var ErrInvalidCredentials = errors.New("models: invalid credentials")
var ErrDuplicateEmail = errors.New("models: duplicate email")
var ErrNotValidPostForm = errors.New("models: no valid post form")
var ErrDuplicateName    = errors.New("models: duplicate name")
