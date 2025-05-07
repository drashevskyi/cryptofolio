package store

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

type TestableDB struct {
	DB *sql.DB
}

func NewTestDB(t *testing.T) *TestableDB {
	dsn := "postgres://postgres:postgres@localhost:5432/cryptofolio_test?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

	// Clean + migrate
	db.Exec("DROP TABLE IF EXISTS assets")
	Migrate(db)

	return &TestableDB{DB: db}
}

func (tdb *TestableDB) Close() {
	tdb.DB.Close()
}
