package postgres

import (
	"adel/internal/service/token"
	"database/sql"
	"time"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{
		db: db,
	}
}

type TokenStore interface {
	Insert(token *token.Token) error
	CreateNewToken(userID int64, ttl time.Duration) (*token.Token, error)
	DeleteAllTokensForUser(userID int64) error
}

func (t *PostgresTokenStore) CreateNewToken(userID int64, ttl time.Duration) (*token.Token, error) {
	token, err := token.GenerateToken(userID, ttl)
	if err != nil {
		return nil, err
	}

	err = t.Insert(token)
	return token, err
}

func (t *PostgresTokenStore) Insert(token *token.Token) error {
	query := `
    INSERT INTO tokens (hash, user_id, expiry)
    VALUES ($1, $2, $3)
    `

	_, err := t.db.Exec(query, token.Hash, token.UserID, token.Expiry)
	return err
}

func (t *PostgresTokenStore) DeleteAllTokensForUser(userID int64) error {
	query := `
    DELETE FROM tokens
    WHERE user_id = $1
    `

	_, err := t.db.Exec(query, userID)
	return err
}
