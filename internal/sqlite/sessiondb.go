package sqlite

import (
	"fmt"
	"forum/internal/models"
)

func (s *Sqlite) GetUserIDByToken(token string) (int, error) {
	stmt := `SELECT user_id FROM sessions WHERE token = ?`
	var userID int
	err := s.DB.QueryRow(stmt, token).Scan(&userID)
	if err != nil {
		return -1, err
	}
	return userID, nil
}

func (s *Sqlite) GetUserNameByUserID(user_id int) (string, error) {
	stmt := `SELECT name FROM users WHERE id = ?`
	var name string
	err := s.DB.QueryRow(stmt, user_id).Scan(&name)
	if err != nil {
		return "This user don't have name", err
	}
	return name, nil
}

func (s *Sqlite) CreateSession(session *models.Session) error {
	op := "sqlite.CreateSession"
	stmt := `INSERT INTO sessions(user_id, token, exp_time) VALUES(?, ?, ?)`
	_, err := s.DB.Exec(stmt, session.UserID, session.Token, session.ExpTime)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Sqlite) DeleteSessionById(userid int) error {
	op := "sqlite.DeleteSessionById"
	stmt := `DELETE FROM sessions WHERE user_id = ?`
	if _, err := s.DB.Exec(stmt, userid); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Sqlite) DeleteSessionByToken(token string) error {
	op := "sqlite.DeleteSessionByToken"
	stmt := `DELETE FROM sessions WHERE token = ?`
	if _, err := s.DB.Exec(stmt, token); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
