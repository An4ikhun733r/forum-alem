package service

import (
	"forum/internal/models"
	"forum/internal/sqlite"
)

type service struct {
	repo sqlite.RepoI
}

type ServiceI interface {
	SnippetRepo
	UserRepo
}

func NewService(repo sqlite.RepoI) ServiceI {
	return &service{repo}
}

type SnippetRepo interface {
	InsertSnippet(cookie, title, content, category string) (int, error)
	GetSnippet(id int) (*models.Snippet, error)
	Latest() ([]*models.Snippet, error)
}

type UserRepo interface {
	InsertUser(username, password, email string) (int, error)
	Authenticate(username, password string) (*models.Session, int, error)
}

