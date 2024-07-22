package models

import (
	"forum/internal/validator"
	"time"
)

type Snippet struct {
	ID       int
	User_id  int
	Title    string
	Content  string
	Likes 	 int
	Dislikes int
	Category string
	Created  time.Time
}

type SnippetCreateForm struct {
	Title   string
	Content string
	Category string
	validator.Validator
}