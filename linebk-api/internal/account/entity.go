package account

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

// Account represents an account entity in the database.
type Account struct {
	AccountID     string         `db:"account_id"`
	UserID        string         `db:"user_id"`
	Type          sql.NullString `db:"type"`
	Currency      sql.NullString `db:"currency"`
	AccountNumber sql.NullString `db:"account_number"`
	Issuer        sql.NullString `db:"issuer"`
	DummyCol3     sql.NullString `db:"dummy_col_3"`
	CreatedAt     time.Time      `db:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at"`
}

// Balance represents the financial balance details associated with an account.
type Balance struct {
	AccountID string          `db:"account_id"`
	UserID    string          `db:"user_id"`
	Amount    decimal.Decimal `db:"amount"`
	DummyCol4 sql.NullString  `db:"dummy_col_4"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}

// Details represents UI and functional attributes for an account.
type Details struct {
	AccountID     string         `db:"account_id"`
	UserID        string         `db:"user_id"`
	Color         sql.NullString `db:"color"`
	IsMainAccount bool           `db:"is_main_account"`
	Progress      uint8          `db:"progress"`
	DummyCol5     sql.NullString `db:"dummy_col_5"`
	CreatedAt     time.Time      `db:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at"`
}

// Flags represents boolean or metadata tags associated with an account.
type Flags struct {
	FlagID    int       `db:"flag_id"`
	AccountID string    `db:"account_id"`
	UserID    string    `db:"user_id"`
	FlagType  string    `db:"flag_type"`
	FlagValue string    `db:"flag_value"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
