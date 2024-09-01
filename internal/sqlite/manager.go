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
	InsertSnippet(name, title, content string, category []string, user_id int) (int, error)
	GetSnippet(id int) (*models.Snippet, error)
	Latest(tags []string, filter string, userID int) ([]*models.Snippet, error)
	AddComment(postId, userId int, content string) error
	GetCommentByPostId(postId int) ([]models.Comment, error)
	GetUserReaction(userID, postID int) (int, error)
	LikePost(userID, postID int) error
	DislikePost(userID, postID int) error
	RemoveReaction(userID, postID int) error
}

type SessionRepo interface {
	GetUserNameByUserID(user_id int) (string, error)
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
