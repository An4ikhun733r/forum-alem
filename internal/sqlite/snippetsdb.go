package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (s *Sqlite) InsertSnippet(title, content, category string, user_id int) (int, error) {
	// SQLite uses different syntax for inserting current time and date calculations.
	stmt := `INSERT INTO snippets (user_id, title, content, category, created)
    VALUES (?, ?, ?, ?, datetime('now'))`
	fmt.Println(title)
	fmt.Println(content)
	fmt.Println(user_id)
	fmt.Println(category)
	// Execute the statement using Exec() method.
	result, err := s.DB.Exec(stmt, user_id, title, content, category)
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
	// Write the SQL statement we want to execute. Again, I've split it over two
	// lines for readability.
	stmt := `SELECT id, user_id, title, content, likes, dislikes, category, created FROM snippets
    WHERE id = ?`

	// Use the QueryRow() method on the connection pool to execute our
	// SQL statement, passing in the untrusted id variable as the value for the
	// placeholder parameter. This returns a pointer to a sql.Row object which
	// holds the result from the database.
	row := s.DB.QueryRow(stmt, id)

	// Initialize a pointer to a new zeroed Snippet struct.
	ss := &models.Snippet{}

	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct. Notice that the arguments
	// to row.Scan are *pointers* to the place you want to copy the data into,
	// and the number of arguments must be exactly the same as the number of
	// columns returned by your statement.
	err := row.Scan(&ss.ID, &ss.User_id, &ss.Title, &ss.Content, &ss.Likes, &ss.Dislikes, &ss.Category,&ss.Created)
	if err != nil {
		// If the query returns no rows, then row.Scan() will return a
		// sql.ErrNoRows error. We use the errors.Is() function check for that
		// error specifically, and return our own ErrNoRecord error
		// instead (we'll create this in a moment).
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	// If everything went OK then return the Snippet object.
	return ss, nil
}

func (s *Sqlite) Latest() ([]*models.Snippet, error) {
	// Write the SQL statement we want to execute.
	stmt := `SELECT id, user_id, title, content, likes, dislikes, category, created FROM snippets ORDER BY id DESC LIMIT 10`

	// Use the Query() method on the connection pool to execute our
	// SQL statement. This returns a sql.Rows resultset containing the result of
	// our query.
	rows, err := s.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// We defer rows.Close() to ensure the sql.Rows resultset is
	// always properly closed before the Latest() method returns. This defer
	// statement should come *after* you check for an error from the Query()
	// method. Otherwise, if Query() returns an error, you'll get a panic
	// trying to close a nil resultset.
	defer rows.Close()

	// Initialize an empty slice to hold the Snippet structs.
	snippets := []*models.Snippet{}

	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by the
	// rows.Scan() method. If iteration over all the rows completes then the
	// resultset automatically closes itself and frees-up the underlying
	// database connection.
	for rows.Next() {
		// Create a pointer to a new zeroed Snippet struct.
		s := &models.Snippet{}
		// Use rows.Scan() to copy the values from each field in the row to the
		// new Snippet object that we created. Again, the arguments to row.Scan()
		// must be pointers to the place you want to copy the data into, and the
		// number of arguments must be exactly the same as the number of
		// columns returned by your statement.
		err = rows.Scan(&s.ID, &s.User_id, &s.Title, &s.Content, &s.Likes, &s.Dislikes, &s.Category, &s.Created)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets.
		snippets = append(snippets, s)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK then return the Snippets slice.
	return snippets, nil
}