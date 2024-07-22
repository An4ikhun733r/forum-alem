package sqlite

import (
	"database/sql"
	"fmt"
)

type Sqlite struct {
	DB *sql.DB
}

func OpenDB(dsn string) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	queries := []string{
		`CREATE TABLE IF NOT EXISTS snippets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		likes INTEGER DEFAULT 0,
    	dislikes INTEGER DEFAULT 0,
		category TEXT NOT NULL,
		created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			token TEXT NOT NULL,
			exp_time TIMESTAMP NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
  		name TEXT NOT NULL,
    	email TEXT NOT NULL,
    	hashed_password CHAR(60) NOT NULL,
    	created DATETIME NOT NULL,
    	UNIQUE(email)
	);`,
	}
	for _, query := range queries {
		stmt, err := db.Prepare(query)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", "sqlite.NewDB", err)
		}
		_, err = stmt.Exec()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", "sqlite.NewDB", err)
		}
		stmt.Close()
	}
	return &Sqlite{DB: db}, nil
}
