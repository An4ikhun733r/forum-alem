package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func (s *Sqlite) GetUserByID(id int) (*models.User, error) {
	op := "sqlite.GetUserByID"
	var u models.User
	stmt := `SELECT id, name, email, created FROM users WHERE id=?`
	err := s.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Name, &u.Email, &u.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &u, nil
}

func (s *Sqlite) Authenticate(username, password string) (int, error) {
	var id int
	var hashedPassword string
	// Prepare the SQL statement to select the hashed password
	stmt := `SELECT id, hashed_password FROM users WHERE name = ?`

	// Scan the result into variables
	err := s.DB.QueryRow(stmt, username).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
		 return 0, models.ErrInvalidCredentials
		} else {
		 return 0, err
		}
	   }

	// Compare the hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		 return 0, models.ErrInvalidCredentials
		} else {
		 return 0, err
		}
	   }

	return id, nil // Authentication successful
}

func (s *Sqlite) InsertUser(username, password, email string) (int, error) {
	var exists bool

    checkStmt := `SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)`
    err := s.DB.QueryRow(checkStmt, email).Scan(&exists)

    if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			return 0, err
		} else {
			return 0, err
		}
	}
    if exists {
        return 0, models.ErrDuplicateEmail
    }
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	// SQLite uses different syntax for inserting current time and date calculations.
	stmt := `INSERT INTO users (name, hashed_password, email, created)
    VALUES (?, ?, ?, datetime('now'))`

	// Execute the statement using Exec() method.
    result, err := s.DB.Exec(stmt, username, hashedPassword, email)
    if err != nil {
        return 0, err
    }

	// Retrieve the last inserted ID using LastInsertId() method.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}