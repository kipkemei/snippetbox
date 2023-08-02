package models

import (
	"database/sql"
	"errors"
	"time"
)

// Define a Snippet type to hold data for an individual snippet.
// The fields should correspond to the fields in MySQL snippets table.
type Snippet struct {
	ID      int
	Title   string
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
	// Write the SQL statement to be executed
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	// Use the Query() method on the connection pool to execute the SQL statement
	// It returns a sql.Rows resultset containing the result of our query.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// defer rows.Close() to ensure the sql.Rows resultset is properly closed before Latest()
	// method returns. This defer should come after the error check from the Query() method.
	// Otherwise, if Query() returns an error, trying to close a nil resultset panics.
	defer rows.Close()

	// Initialize an empty slice to hold the Snippet structs.
	snippets := []*Snippet{}

	// Use rows.Next to iterate through the rows in the resultset. This prepares the first
	// (and then each subsequent) row to be acted on by rows.Scan() method. If iteration over
	// all the rows completes then the resultset automatically closes itself and frees-up the
	// underlying database connection.
	for rows.Next() {
		// Create a pointer to a new zeroed Snippet struct.
		s := &Snippet{}

		// Use rows.Scan() to copy the values from each field in the row to the new Snippet
		// object that we created. Again, the arguments to row.Scan() must be pointers to
		// the place you want to copy the data into, and the no. of arguments must be exactly
		// same as the number of columns returned by the SQL statement.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append to the slice of snippets.
		snippets = append(snippets, s)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error that
	// was encountered during the iteration. It is IMPORTANT to call this. DON'T assume that
	// a successful iteration was completed over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK then return the Snippets slice.
	return snippets, nil
}
