package service

import (
	"database/sql"
	"errors"
	"forum/internal/models"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (s *service) InsertSnippet(cookie, title, content string, category []string) (int, error) {
	userID, err := s.repo.GetUserIDByToken(cookie)
	if err != nil {
		return 0, err
	}

	name, err := s.repo.GetUserNameByUserID(userID)
	if err != nil {
		return 0, err
	}

	postID, err := s.repo.InsertSnippet(name, title, content, category, userID)
	if err != nil {
		return 0, err
	}

	return int(postID), nil
}

func (s *service) GetSnippet(id int) (*models.Snippet, error) {
	return s.repo.GetSnippet(id)
}

func (s *service) Latest(tags []string, filter string, userID int) ([]*models.Snippet, error) {
	snippets, err := s.repo.Latest(tags, filter, userID)
	if err != nil {
		return nil, err
	}
	return snippets, nil
}

func (s *service) DislikePost(userID, postID int) error {
	reaction, err := s.repo.GetUserReaction(userID, postID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if reaction == 1 || reaction == -1 {
		if err := s.repo.RemoveReaction(userID, postID); err != nil {
			return err
		}
	}

	if reaction != -1 {
		if err := s.repo.DislikePost(userID, postID); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) LikePost(userID, postID int) error {
	reaction, err := s.repo.GetUserReaction(userID, postID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if reaction == 1 || reaction == -1 {
		if err := s.repo.RemoveReaction(userID, postID); err != nil {
			return err
		}
	}

	if reaction != 1 {
		if err := s.repo.LikePost(userID, postID); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) AddComment(postId, userId int, content string) error {
	err := s.repo.AddComment(postId, userId, content)
	return err
}

func (s *service) GetCommentByPostId(postId int) ([]models.Comment, error) {
	comment, err := s.repo.GetCommentByPostId(postId)
	if err != nil {
		return nil, err
	}
	return comment, nil
}