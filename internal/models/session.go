package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	UserID  int
	Token   string
	ExpTime time.Time
}

func NewSession(UserID int) *Session {
	return &Session{
		UserID:  UserID,
		Token:   uuid.New().String(),
		ExpTime: time.Now().Add(120 * time.Minute),
	}
}