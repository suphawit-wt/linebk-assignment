package user

import (
	"database/sql"
	"time"
)

// User represents a core user entity with authentication credentials and timestamp metadata.
type User struct {
	UserID       string         `db:"user_id"`
	Name         sql.NullString `db:"name"`
	DummyCol1    sql.NullString `db:"dummy_col_1"`
	PasswordHash sql.NullString `db:"password_hash"`
	PasscodeHash sql.NullString `db:"passcode_hash"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
}

// Greeting represents user greeting metadata associated with an user profile.
type Greeting struct {
	UserID    string         `db:"user_id"`
	Greeting  sql.NullString `db:"greeting"`
	DummyCol2 sql.NullString `db:"dummy_col_2"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}
