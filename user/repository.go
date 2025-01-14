package user

import (
	"database/sql"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"time"
)

// Repository is interface that will manage user data.
type Repository interface {
	// CheckUserEmail will check if the user email is in use.
	CheckUserEmail(email string) (bool, error)

	// AddUser will add a new user to the database.
	AddUser(user *User) error

	// GetUserByEmail will get user by its email.
	GetUserByEmail(email string) (*User, error)

	// DeleteUser will delete user by its id.
	DeleteUser(id *uuid.UUID) error

	// AddToken will add a new token.
	AddToken(tokenID, userId *uuid.UUID, exp *time.Time) error

	// DeleteTokenById will delete a token and return true if the token was deleted.
	DeleteTokenById(id *uuid.UUID) (bool, error)

	DeleteTokenByUserId(id *uuid.UUID) error
}

type PostgresRepository struct {
	database *sql.DB
}

func NewPostgresRepository(database *sql.DB) *PostgresRepository {
	return &PostgresRepository{database: database}
}

func (p *PostgresRepository) CheckUserEmail(email string) (bool, error) {
	row := p.database.QueryRow("SELECT COUNT(id) FROM users WHERE email = $1", email)
	var count int

	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (p *PostgresRepository) AddUser(user *User) error {
	_, err := p.database.Exec("INSERT INTO users (id, email, password) VALUES ($1, $2, $3)", user.Id, user.Email, user.Password)
	return err
}

func (p *PostgresRepository) GetUserByEmail(email string) (*User, error) {
	row := p.database.QueryRow("SELECT * FROM users WHERE email = $1", email)
	var user User
	err := row.Scan(&user.Id, &user.Email, &user.Password)
	return &user, err
}

func (p *PostgresRepository) DeleteUser(id *uuid.UUID) error {
	_, err := p.database.Exec("DELETE FROM users WHERE id = $1", *id)
	return err
}

func (p *PostgresRepository) AddToken(userId, tokenId *uuid.UUID, exp *time.Time) error {
	_, err := p.database.Exec("INSERT INTO tokens(id, user_id, expire_date) VALUES ($1, $2, $3)", *tokenId, *userId, *exp)
	return err
}

func (p *PostgresRepository) DeleteTokenById(id *uuid.UUID) (bool, error) {
	result, err := p.database.Exec("DELETE FROM tokens WHERE id = $1", *id)
	if err != nil {
		return false, err
	}
	count, err := result.RowsAffected()
	return count > 0, err
}

func (p *PostgresRepository) DeleteTokenByUserId(userId *uuid.UUID) error {
	_, err := p.database.Exec("DELETE FROM tokens WHERE user_id = $1", *userId)
	return err
}
