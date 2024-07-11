package models

import (
	"database/sql"
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

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Need to convert because the type that lastInsertId returns is int64.
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	return nil, nil
}

// Returns the 	10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
