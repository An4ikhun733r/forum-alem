package service

import (
	"database/sql"
	"fmt"
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

func (s *service) InsertSnippet(cookie, title, content, category string) (int, error) {
	userID, err := s.repo.GetUserIDByToken(cookie)
	if err != nil {
		return 0, err
	}
	fmt.Println(cookie)
	fmt.Println(title)
	fmt.Println(content)
	fmt.Println(category)
	postID, err := s.repo.InsertSnippet(title, content, category, userID)
	if err != nil {
		return 0, err
	}

	return int(postID), nil
}

func (s *service) GetSnippet(id int) (*models.Snippet, error) {
	return s.repo.GetSnippet(id)
}

func (s *service) Latest() ([]*models.Snippet, error) {
	snippets, err := s.repo.Latest()
	if err != nil {
		return nil, err
	}
	return snippets, nil
}
