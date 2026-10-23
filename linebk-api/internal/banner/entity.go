package banner

import (
	"database/sql"
	"time"
)

// Banner represents a promotional banner entity in the database
type Banner struct {
	BannerID    string         `db:"banner_id"`
	UserID      string         `db:"user_id"`
	Title       sql.NullString `db:"title"`
	Description sql.NullString `db:"description"`
	Image       sql.NullString `db:"image"`
	DummyCol11  sql.NullString `db:"dummy_col_11"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}
