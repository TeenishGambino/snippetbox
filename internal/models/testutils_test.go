package models

import (
	"database/sql"
	"os"
	"testing"
)

func newTestDB(t *testing.T) *sql.DB {
	// This is for docker : test_web:pass@(docker.for.mac.localhost:3306)/test_snippetbox?parseTime=true&multiStatements=true
	// This is for normal : test_web:pass@/test_snippetbox?parseTime=true&multiStatements=true
	// Use test_web instead of root for local
	// Maybe environment variables are better here?
	db, err := sql.Open("mysql", "root:pass@(docker.for.mac.localhost:3306)/test_snippetbox?parseTime=true&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}

	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		script, err := os.ReadFile("./testdata/teardown.sql")

		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	})
	// Return the database connection pool.
	return db
}
