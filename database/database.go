package database

import (
	"database/sql"
	"fmt"
)

// CreateConnection will create a database connection by using user, password, host and database name.
func CreateConnection(user, password, host, name string) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", user, password, host, name)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
