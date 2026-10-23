package card

import (
	"database/sql"
	"time"
)

// DebitCard represents a debit card entity in the database
type DebitCard struct {
	CardID    string         `db:"card_id"`
	UserID    string         `db:"user_id"`
	Name      sql.NullString `db:"name"`
	DummyCol7 sql.NullString `db:"dummy_col_7"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

// DebitCardStatus represents the operational status of a debit card
type DebitCardStatus struct {
	CardID    string         `db:"card_id"`
	UserID    string         `db:"user_id"`
	Status    sql.NullString `db:"status"`
	DummyCol8 sql.NullString `db:"dummy_col_8"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

// DebitCardDesign represents the visual styling properties of a debit card
type DebitCardDesign struct {
	CardID      string         `db:"card_id"`
	UserID      string         `db:"user_id"`
	Color       sql.NullString `db:"color"`
	BorderColor sql.NullString `db:"border_color"`
	DummyCol9   sql.NullString `db:"dummy_col_9"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}

// DebitCardDetails represents technical issuer and numbering details for a debit card
type DebitCardDetails struct {
	CardID     string         `db:"card_id"`
	UserID     string         `db:"user_id"`
	Issuer     sql.NullString `db:"issuer"`
	Number     sql.NullString `db:"number"`
	DummyCol10 sql.NullString `db:"dummy_col_10"`
	CreatedAt  time.Time      `db:"created_at"`
	UpdatedAt  time.Time      `db:"updated_at"`
}
