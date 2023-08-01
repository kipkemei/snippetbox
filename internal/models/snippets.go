package models

import (
	"database/sql"
	"errors"
	"time"
)

// Define a Snippet type to hold data for an individual snippet.
// The fields should correspond to the fields in MySQL snippets table.
type Snippet struct {
	ID int
	Title string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// Write the SQL statement to be executed
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(? ? UTC_TIMESTAMP(), DATE_ADD(UTC_TMESTAMP(), INTERVAL ? DAY))`

	// Use DB.Exec() on the embedded connection pool to execute the statement.
	// The first parameter is the SQL statement, followed by fields values for
	// placeholder parameters. 
	// This method returns a sql.Result type, which contains basic information about
	// what happened when the statement was executed.
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Use the LastInsertId() method on the result to get the ID of our newly inserted
	// record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID returned has the type int64, so convert it to an int type before returning.
	return int(id), nil
}

// This will fetch a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	// Write the SQL statement to be executed
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id=?`

	// Use the QueryRow() method on the connection pool to execute the SQL statement.
	// Pass in the untrusted id variable as the value for the placeholder parameter.
	// This returns a pointer to a sql.Row object which holds the result from db.
	row := m.DB.QueryRow(stmt, id)

	// Initialize a pointer to a new zeroed Snippet struct
	s := &Snippet{}

	// Use row.Scan() to copy the values from each field in sql.Row to the corresponding
	// field in the Snippet struct. The arguments to row.Scan() are *pointers* to the place
	// you want to copy the data into, and the no. of arguments must be exactly the same as 
	// the number of columns returned by the statement.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// If the query returns no rows, then row.Scan() will return a sql.ErrNoRows error.
		// Use errors.Is() to check the specific error it is, and return our own ErrNoRecord
		// error instead.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	// If everything went OK then return the Snippet object.
	return s, nil
}

// This will return the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}