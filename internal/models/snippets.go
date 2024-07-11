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
	return 0, nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	return nil, nil
}

// Returns the 	10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
