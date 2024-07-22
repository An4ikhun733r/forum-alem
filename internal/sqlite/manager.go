package sqlite

import "forum/internal/models"

func NewRepo(dsn string) (RepoI, error) {
	return OpenDB(dsn)
}

type RepoI interface {
	PostRepo
	SessionRepo
	UserRepo
}

type PostRepo interface {
	InsertSnippet(title, content, category string, user_id int) (int, error)
	GetSnippet(id int) (*models.Snippet, error)
	Latest() ([]*models.Snippet, error)
}

type SessionRepo interface {
	GetUserIDByToken(token string) (int, error)
	CreateSession(session *models.Session) error
	DeleteSessionById(userid int) error
	DeleteSessionByToken(token string) error
}

type UserRepo interface {
	GetUserByID(id int) (*models.User, error)
	Authenticate(username, password string) (int, error)
	InsertUser(username, password, email string) (int, error)
}