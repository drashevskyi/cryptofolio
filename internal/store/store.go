package store

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	dsn := os.Getenv("POSTGRES_DSN")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}
	return db
}

func Migrate(db *sql.DB) {
	query := `
    CREATE TABLE IF NOT EXISTS assets (
        id SERIAL PRIMARY KEY,
        "user" TEXT NOT NULL,
        label TEXT NOT NULL,
        currency TEXT NOT NULL CHECK (currency IN ('BTC', 'ETH', 'LTC')),
        amount NUMERIC NOT NULL CHECK (amount >= 0)
    );
    `
	if _, err := db.Exec(query); err != nil {
		log.Fatal("Failed to run migration:", err)
	}
}
