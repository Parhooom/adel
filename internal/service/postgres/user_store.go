package postgres

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plaintText *string
	hash       []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.plaintText = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash password  `json:"-"`
	IsAdmin      bool      `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

type UserStore interface {
	CreateUser(user *User) error
	GetUserByUsername(username string) (*User, error)
	GetUserToken(tokenPlaintext string) (*User, error)
}

var (
	ErrDuplicateUsername = errors.New("username already exists")
	ErrUserNotFound      = errors.New("user not found")
)

func (s *PostgresUserStore) CreateUser(user *User) error {
	query := `
	INSERT INTO users (username, password_hash, is_admin)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRow(query, user.Username, user.PasswordHash.hash, user.IsAdmin).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err.Error() == `ERROR: duplicate key value violates unique constraint "users_username_key" (SQLSTATE 23505)` {
			return ErrDuplicateUsername
		}

		return err
	}

	return nil
}

func (s *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}

	query := `
	SELECT u.id, u.username, u.password_hash, u.is_admin, u.created_at, u.updated_at
	FROM users u
	WHERE u.username = $1
	`

	err := s.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.PasswordHash.hash, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresUserStore) GetUserToken(tokenPlaintext string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
  	SELECT u.id, u.username, u.password_hash, u.is_admin, u.created_at, u.updated_at
  	FROM users u
  	JOIN tokens t ON t.user_id = u.id
  	WHERE t.hash = $1 AND t.expiry > $2
  	`

	user := &User{
		PasswordHash: password{},
	}

	err := s.db.QueryRow(query, tokenHash[:], time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash.hash,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}
