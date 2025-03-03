package repositories

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
)

// UserRepository interface manages the data of users.
type UserRepository interface {
	// CheckIfEmailExists will return true if the email is in use otherwise false.
	CheckIfEmailExists(ctx context.Context, email string) (bool, error)

	// CheckIfUsernameExists will return true if the email is in use otherwise false.
	CheckIfUsernameExists(ctx context.Context, username string) (bool, error)

	// AddUser will insert a new user.
	AddUser(ctx context.Context, email string, username string, password string) error
}

// PostgresUserRepository struct manages data using connection to postgres database.
type PostgresUserRepository struct {
	db *sql.DB
}

func (r *PostgresUserRepository) CheckIfEmailExists(ctx context.Context, email string) (bool, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM users 
		WHERE email = $1`,
		email,
	)

	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PostgresUserRepository) CheckIfUsernameExists(ctx context.Context, username string) (bool, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM users 
		WHERE username = $1`,
		username,
	)

	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PostgresUserRepository) AddUser(ctx context.Context, email string, username string, password string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (email, username, password) 
		VALUES ($1, $2, $3)`,
	)
	return err
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}
