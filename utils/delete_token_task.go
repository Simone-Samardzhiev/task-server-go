package utils

import (
	"database/sql"
	"log"
	"time"
)

// StartDeleteTokenTask will start 24 ticker that will be used to call [deleteTokens]
func StartDeleteTokenTask(db *sql.DB) {
	err := deleteTokens(db) // Run immediately
	if err != nil {
		log.Printf("Error deleting tokens: %v", err)
	}

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		<-ticker.C
		err := deleteTokens(db)
		if err != nil {
			log.Printf("Error deleting tokens: %v", err)
		}
	}
}

// deleteTokens will delete all tokens that are expired.
func deleteTokens(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM tokens WHERE exp < NOW()")
	return err
}
