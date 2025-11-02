package repository

import (
	"database/sql"
	"fmt"
	"meme-url/internal/models"

	"math/rand"

	log "github.com/sirupsen/logrus"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var memeWords = []string{
	"sigma", "based", "yeet", "big_chungus", "amogus",
	"cringe", "chad", "omegalul", "bruh", "sus",
}

func GetShortLink(db *sql.DB, url string) (*models.Response, error) {
	var shortcode string

	err := db.QueryRow("SELECT shortcode FROM urls WHERE long_url = $1", url).Scan(&shortcode)
	if err == sql.ErrNoRows {
		shortcode, err = generateShortcode(db)
		if err != nil {
			log.Errorf("error while generating a shortcode: %v", err)
			return nil, err
		}

		_, err = db.Exec("INSERT INTO urls (long_url, shortcode) VALUES ($1, $2)", url, shortcode)
		if err != nil {
			log.Errorf("failed to insert url: %v", err)
			return nil, err
		}

	} else if err != nil {
		log.Errorf("database error: %v", err)
		return nil, err
	}

	response := &models.Response{
		UrlOriginal:  url,
		Shortcode:    shortcode,
		UrlShortened: "http://localhost:3030/" + shortcode,
	}
	return response, nil

}

func GetLongURL(db *sql.DB, shortcode string) (string, error) {
	var longUrl string
	err := db.QueryRow("SELECT long_url FROM urls WHERE shortcode = $1", shortcode).Scan(&longUrl)
	return longUrl, err
}

func generateShortcode(db *sql.DB) (string, error) {
	maxAttempts := 10

	for i := 0; i < maxAttempts; i++ {
		randomWord := memeWords[rand.Intn(len(memeWords))]

		randomStr := make([]byte, 6)
		for j := 0; j < 6; j++ {
			randomStr[j] = charset[rand.Intn(len(charset))]
		}

		shortcode := randomWord + "_" + string(randomStr)

		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM urls WHERE shortcode = $1)", shortcode).Scan(&exists)
		if err != nil {
			log.Errorf("database error: %v", err)
			return "", err
		}

		if !exists {
			return shortcode, err
		}
	}
	return "", fmt.Errorf("failed to generate unique shortcode after %d attempts", maxAttempts)
}

func DeleteShortcode(db *sql.DB, shortcode string) error {
	result, err := db.Exec("DELETE FROM urls WHERE shortcode = $1", shortcode)
	if err != nil {
		log.Errorf("database error: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Errorf("failed to get rows affected for shortcode %s: %v", shortcode, err)
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return err
}
