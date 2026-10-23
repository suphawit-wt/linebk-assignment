package transaction

import (
	"database/sql"
	"time"
)

// Transaction represents a financial transaction entity in the database
type Transaction struct {
	TransactionID string         `db:"transaction_id"`
	UserID        string         `db:"user_id"`
	Name          sql.NullString `db:"name"`
	Image         sql.NullString `db:"image"`
	IsBank        bool           `db:"isBank"`
	DummyCol6     sql.NullString `db:"dummy_col_6"`
	CreatedAt     time.Time      `db:"created_at"`
}
