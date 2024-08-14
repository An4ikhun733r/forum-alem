package models

import (
	"forum/internal/validator"
	"time"
)

type Snippet struct {
	ID       int
	User_id  int
	Name     string
	Title    string
	Content  string
	Likes    int
	Dislikes int
	Category []string
	Created  time.Time
}

type SnippetCreateForm struct {
	Title    string
	Content  string
	Category []string
	validator.Validator
}

type UserPostReaction struct {
	UserID   int
	PostID   int
	Reaction int // 1 для лайка, -1 для дизлайка
}

type Comment struct {
	ID       int
	PostId   int
	UserId   int
	Username string
	Content  string
	Created  time.Time
}
