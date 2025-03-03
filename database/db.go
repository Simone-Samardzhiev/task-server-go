package database

import (
	"database/sql"
	"server/config"
)

// Connect will connect to the database.
func Connect(conf *config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", conf.Url)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(conf.MaxOpenConnections)
	db.SetMaxIdleConns(conf.MaxIdleConnections)

	return db, nil
}
