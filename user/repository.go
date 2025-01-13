package user

import (
	"database/sql"
	"github.com/google/uuid"
	"server/auth"
)

// Repository is interface that will manage user data.
type Repository interface {
	// CheckUserEmail will check if the user email is in use.
	CheckUserEmail(email *string) (bool, error)

	// AddUser will add a new user to the database.
	AddUser(user *User) error

	// GetUserByEmail will get user by its email.
	GetUserByEmail(email *string) (User, error)

	// DeleteUser will delete user by its id.
	DeleteUser(id *uuid.UUID) error

	// AddToken will add a new token.
	AddToken(token *auth.RefreshTokenClaims) error

	// DeleteToken will delete a token and return true if the token was deleted.
	DeleteToken(token *auth.RefreshTokenClaims) (bool, error)
}

type PostgresRepository struct {
	database *sql.DB
}

func (p *PostgresRepository) CheckUserEmail(email *string) (bool, error) {
	row := p.database.QueryRow("SELECT COUNT(id) FROM users WHERE email = ?", *email)
	var count int

	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (p *PostgresRepository) AddUser(user *User) error {
	_, err := p.database.Exec("INSERT INTO users (id, email, password) VALUES (?, ?, ?)", user.Id, user.Email, user.Password)
	return err
}

func (p *PostgresRepository) GetUserByEmail(email *string) (User, error) {
	row := p.database.QueryRow("SELECT * FROM users WHERE email = ?", *email)
	var user User
	err := row.Scan(&user.Id, &user.Email, &user.Password)
	return user, err
}

func (p *PostgresRepository) DeleteUser(id *uuid.UUID) error {
	_, err := p.database.Exec("DELETE FROM users WHERE id = ?", *id)
	return err
}

func (p *PostgresRepository) AddToken(token *auth.RefreshTokenClaims) error {
	_, err := p.database.Exec("INSERT INTO tokens(id, expire_date) VALUES (?, ?)", token.ID, token.ExpiresAt)
	return err
}

func (p *PostgresRepository) DeleteToken(token *auth.RefreshTokenClaims) (bool, error) {
	result, err := p.database.Exec("DELETE FROM tokens WHERE id = ?", token.ID)
	if err != nil {
		return false, err
	}
	count, err := result.RowsAffected()
	return count > 0, err
}
