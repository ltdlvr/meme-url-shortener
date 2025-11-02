package db

import (
	"database/sql"

	log "github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

func InitDatabase() *sql.DB {
	//надо добавить cfg
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=meme-url sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open db: %s", err)
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS urls (
		id SERIAL PRIMARY KEY,
		shortcode VARCHAR(20) UNIQUE,
		long_url TEXT,
		created_at TIMESTAMP)`)
	if err != nil {
		log.Fatalf("failed to create db: %s", err)
	}

	return db
}
