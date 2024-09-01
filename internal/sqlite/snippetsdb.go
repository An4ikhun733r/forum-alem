package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
	"strings"
)

type SnippetModel struct {
	DB *sql.DB
}

func (s *Sqlite) InsertSnippet(name, title, content string, category []string, user_id int) (int, error) {
	categoryStr := strings.Join(category, ", ")
	// SQLite uses different syntax for inserting current time and date calculations.
	stmt := `INSERT INTO snippets (user_id, name, title, content, category, created)
    VALUES (?, ?, ?, ?, ?, datetime('now'))`

	// Execute the statement using Exec() method.
	result, err := s.DB.Exec(stmt, user_id, name, title, content, categoryStr)
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

func (s *Sqlite) GetSnippet(id int) (*models.Snippet, error) {
	// Write the SQL statement we want to execute.
	stmt := `SELECT id, user_id, name, title, content, likes, dislikes, category, created FROM snippets
        WHERE id = ?`

	// Use the QueryRow() method on the connection pool to execute our
	// SQL statement, passing in the untrusted id variable as the value for the
	// placeholder parameter. This returns a pointer to a sql.Row object which
	// holds the result from the database.
	row := s.DB.QueryRow(stmt, id)

	// Initialize a pointer to a new zeroed Snippet struct.
	ss := &models.Snippet{}
	var categoryStr string

	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct.
	err := row.Scan(&ss.ID, &ss.User_id, &ss.Name, &ss.Title, &ss.Content, &ss.Likes, &ss.Dislikes, &categoryStr, &ss.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	// Convert the comma-separated string to a slice of strings.
	ss.Category = strings.Split(categoryStr, ",")

	// If everything went OK then return the Snippet object.
	return ss, nil
}

func (s *Sqlite) Latest(tags []string, filter string, userID int) ([]*models.Snippet, error) {
	// Base SQL statement
	stmt := `SELECT id, user_id, name, title, content, likes, dislikes, category, created 
             FROM snippets 
             WHERE 1=1`

	// Prepare arguments for the query
	args := []interface{}{}

	// Handle the 'liked' filter separately
	if filter == "liked" {
		likedStmt := `SELECT s.id 
                      FROM snippets s
                      JOIN user_post_reactions ul ON s.id = ul.post_id
                      WHERE ul.user_id = ? AND ul.reaction > 0`
		args = append(args, userID)

		// Modify the main statement to use the results of the liked filter
		stmt += ` AND id IN (` + likedStmt + `)`
	} else if filter == "myPosts" {
		stmt += ` AND user_id = ?`
		args = append(args, userID)
	}

	// Add tag filters to the SQL statement if there are tags
	if len(tags) > 0 {
		stmt += ` AND (`
		for i, _ := range tags {
			if i > 0 {
				stmt += ` OR `
			}
			stmt += `category LIKE ?`
			args = append(args, "%"+tags[i]+"%")
		}
		stmt += `)`
	}

	// Add ordering and limit
	stmt += ` ORDER BY id DESC LIMIT 10`

	// Execute the query
	rows, err := s.DB.Query(stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize an empty slice to hold the Snippet structs
	snippets := []*models.Snippet{}

	// Iterate through the rows
	for rows.Next() {
		s := &models.Snippet{}
		var categoryStr string

		// Scan row values
		err = rows.Scan(&s.ID, &s.User_id, &s.Name, &s.Title, &s.Content, &s.Likes, &s.Dislikes, &categoryStr, &s.Created)
		if err != nil {
			return nil, err
		}

		// Convert category string to slice of strings
		s.Category = strings.Split(categoryStr, ",")

		// Append the snippet to the slice
		snippets = append(snippets, s)
	}

	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Return the slice of snippets
	return snippets, nil
}

func (s *Sqlite) AddComment(postId, userId int, content string) error {
	op := "sqlite.AddComment"
	stmt, err := s.DB.Prepare("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}
	_, err = stmt.Exec(postId, userId, content)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sqlite) GetCommentByPostId(postId int) ([]models.Comment, error) {
	op := "sqlite.GetCommentByPostId"

	query := `
        SELECT c.id, c.post_id, c.user_id, u.name, c.content, c.created_at
        FROM comments c
        JOIN users u ON c.user_id = u.id
        WHERE c.post_id = ?
    `
	rows, err := s.DB.Query(query, postId)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	defer rows.Close()
	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.PostId, &comment.UserId, &comment.Username, &comment.Content, &comment.Created)
		if err != nil {
			return nil, fmt.Errorf("%s : %w", op, err)
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (s *Sqlite) GetUserReaction(userID, postID int) (int, error) {
	op := "sqlite.GetUserReaction"
	var reaction int
	err := s.DB.QueryRow(`SELECT reaction FROM user_post_reactions WHERE user_id = ? AND post_id = ?`, userID, postID).Scan(&reaction)
	if err != nil {
		return 0, fmt.Errorf("%s : %w", op, err)
	}
	return reaction, nil
}

func (s *Sqlite) LikePost(userID, postID int) error {
	op := "sqlite.LikePost"
	_, err := s.DB.Exec(`INSERT INTO user_post_reactions (user_id, post_id, reaction) VALUES (?, ?, 1)`, userID, postID)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}
	_, err = s.DB.Exec(`UPDATE snippets SET likes = likes + 1 WHERE id = ?`, postID)
	return err
}

func (s *Sqlite) DislikePost(userID, postID int) error {
	op := "sqlite.DislikePost"
	_, err := s.DB.Exec(`INSERT INTO user_post_reactions (user_id, post_id, reaction) VALUES (?, ?, -1)`, userID, postID)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}
	_, err = s.DB.Exec(`UPDATE snippets SET dislikes = dislikes + 1 WHERE id = ?`, postID)
	return err
}

func (s *Sqlite) RemoveReaction(userID, postID int) error {
	op := "sqlite.RemoveReaction"
	var reaction int
	err := s.DB.QueryRow(`SELECT reaction FROM user_post_reactions WHERE user_id = ? AND post_id = ?`, userID, postID).Scan(&reaction)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}
	_, err = s.DB.Exec(`DELETE FROM user_post_reactions WHERE user_id = ? AND post_id = ?`, userID, postID)
	if err != nil {
		return err
	}
	if reaction == 1 {
		_, err = s.DB.Exec(`UPDATE snippets SET likes = likes - 1 WHERE id = ?`, postID)
		if err != nil {
			return err
		}
	} else {
		_, err = s.DB.Exec(`UPDATE snippets SET dislikes = dislikes - 1 WHERE id = ?`, postID)
		if err != nil {
			return err
		}
	}
	return nil
}
