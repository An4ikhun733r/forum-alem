package service

import (
	"errors"
	"forum/internal/app"
	"forum/internal/models"
	"net/http"
)

func (s *service) InsertUser(username, password, email string) (int, error) {
	userId, err := s.repo.InsertUser(username, password, email)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			return 0, models.ErrDuplicateEmail
		} else if errors.Is(err, models.ErrDuplicateName) {
			return 0, models.ErrDuplicateName
		} else {
			return 0, err
		}
	}

	return int(userId), nil
}

// Authenticate checks the user's credentials.
func (s *service) Authenticate(username, password string) (*models.Session, int, error) {
	userId, err := s.repo.Authenticate(username, password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			return nil, 0, models.ErrInvalidCredentials
		} else {
			return nil, 0, err
		}
	}
	session := models.NewSession(userId)
	if err = s.repo.DeleteSessionById(userId); err != nil {
		return nil, 0, err
	}
	err = s.repo.CreateSession(session)
	if err != nil {
		return nil, 0, err
	}
	return session, int(userId), nil
}

func (s *service) GetUser(r *http.Request) (*models.User, error) {
	token := app.GetSessionCookie("session_id", r)
	userID, err := s.repo.GetUserIDByToken(token.Value)
	if err != nil {
		return nil, err
	}
	return s.repo.GetUserByID(userID)
}

func (s *service) DeleteSession(token string) error {
	if err := s.repo.DeleteSessionByToken(token); err != nil {
		return err
	}
	return nil
}