package db

import (
	"database/sql"
	"server/config"
)

// connection will connect to the database.
func connect(conf *config.DatabaseConfig) (*sql.DB, error) {
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
