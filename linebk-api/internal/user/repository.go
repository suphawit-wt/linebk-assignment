package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

// Info aggregates user profile details and greeting messages for query results.
type Info struct {
	Name            sql.NullString `db:"name"`
	GreetingMessage sql.NullString `db:"greeting_message"`
}

// Repository defines data access methods for managing user.
type Repository interface {
	GetUserByUserID(ctx context.Context, userID string) (*User, error)
	GetUserInfo(ctx context.Context, userID string) (*Info, error)
}

type repositoryDB struct {
	db *sqlx.DB
}

// NewRepositoryDB creates a new instance of user repositoryDB with the provided database connection pool.
func NewRepositoryDB(db *sqlx.DB) Repository {
	return &repositoryDB{db: db}
}

func (r *repositoryDB) GetUserByUserID(ctx context.Context, userID string) (*User, error) {
	query := `
		SELECT 
			user_id,
			name,
			password_hash,
			passcode_hash
		FROM users
		WHERE user_id = ?
	`

	var user User
	if err := r.db.GetContext(ctx, &user, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *repositoryDB) GetUserInfo(ctx context.Context, userID string) (*Info, error) {
	query := `
		SELECT
			u.name,
			g.greeting AS greeting_message
		FROM users u
		LEFT JOIN user_greetings g ON u.user_id = g.user_id
		WHERE u.user_id = ?
	`

	var userInfo Info
	if err := r.db.GetContext(ctx, &userInfo, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &userInfo, nil
}
