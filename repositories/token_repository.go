package repositories

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"time"
)

// TokenRepository interface managers the tokens data.
type TokenRepository interface {
	// AddToken will add a new token.
	AddToken(ctx context.Context, tokenId uuid.UUID, exp time.Time, userId int) error

	// DeleteToken will delete a token by its id.
	DeleteToken(ctx context.Context, tokenId uuid.UUID) error

	// CheckToken will search the token id the database and return its subject - The user id.
	CheckToken(ctx context.Context, tokenId uuid.UUID) (int, error)
}

type PostgresTokenRepository struct {
	db *sql.DB
}

func (r *PostgresTokenRepository) AddToken(ctx context.Context, tokenId uuid.UUID, exp time.Time, userId int) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO tokens (id, exp, user_id)
		VALUES ($1, $2, $3)`,
		tokenId,
		exp,
		userId,
	)

	return err
}

func (r *PostgresTokenRepository) DeleteToken(ctx context.Context, tokenId uuid.UUID) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM tokens
       WHERE id = $1`,
		tokenId,
	)

	return err
}

func (r *PostgresTokenRepository) CheckToken(ctx context.Context, tokenId uuid.UUID) (int, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id FROM tokens
               WHERE id = $1`,
		tokenId,
	)

	var userId int
	err := row.Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func NewPostgresTokenRepository(db *sql.DB) *PostgresTokenRepository {
	return &PostgresTokenRepository{
		db: db,
	}
}
