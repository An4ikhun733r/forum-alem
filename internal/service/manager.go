package service

import (
	"forum/internal/models"
	"forum/internal/sqlite"
	"net/http"
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
	InsertSnippet(cookie, title, content string, category []string) (int, error)
	GetSnippet(id int) (*models.Snippet, error)
	Latest(tags []string, filter string, userID int) ([]*models.Snippet, error)
	DislikePost(userID, postID int) error
	LikePost(userID, postID int) error
	AddComment(postId, userId int, content string) error
	GetCommentByPostId(postId int) ([]models.Comment, error)
}

type UserRepo interface {
	InsertUser(username, password, email string) (int, error)
	Authenticate(username, password string) (*models.Session, int, error)
	DeleteSession(token string) error
	GetUser(r *http.Request) (*models.User, error)
}
