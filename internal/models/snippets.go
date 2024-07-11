package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// This just a reference to the database instance
// Encapsulating the database connection pool within a higher-level of abstraction.
// sql.DB is just 0 or more connections to the database.
// This is kind of like the middle man that separates the "ui" layer with the data layer
type SnippetModel struct {
	DB *sql.DB
}

// Remember this is just like a function in a class
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	//Double quotes also works if its in a single line, else backmarks can be used for multi lines
	stmt := `INSERT INTO snippets (title, content, created, expires) 
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// This function is supported my MySql but not by PostgreSQL
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Need to convert because the type that lastInsertId returns is int64.
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id=?`

	row := m.DB.QueryRow(stmt, id)

	s := &Snippet{}

	//After the query is completed by QueryRow, we get the data, and
	// .Scan just copies the values to s.
	// Scan function automatically converts raw output from SQL to native Go types.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// If query returns no rows, then .Scan will return an ErrNoRows.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}

	}
	return s, nil

	return nil, nil
}

// Returns the 	10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
