package card

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

// DebitCardInfo aggregates card metadata, status, design, and details for query results.
type DebitCardInfo struct {
	CardID      string         `db:"card_id"`
	UserID      string         `db:"user_id"`
	Name        sql.NullString `db:"name"`
	Status      sql.NullString `db:"status"`
	Issuer      sql.NullString `db:"issuer"`
	Number      sql.NullString `db:"number"`
	Color       sql.NullString `db:"color"`
	BorderColor sql.NullString `db:"border_color"`
}

// Repository defines data access methods for managing card
type Repository interface {
	GetDebitCards(ctx context.Context, userID string) ([]DebitCardInfo, error)
}

type repositoryDB struct {
	db *sqlx.DB
}

// NewRepositoryDB creates a new instance of card repositoryDB with the provided database connection pool.
func NewRepositoryDB(db *sqlx.DB) Repository {
	return &repositoryDB{db: db}
}

func (r *repositoryDB) GetDebitCards(
	ctx context.Context,
	userID string,
) ([]DebitCardInfo, error) {
	query := `
		SELECT 
			dc.card_id,
			dc.name,
			ds.status,
			dd.issuer,
			dd.number,
			dgn.color,
			dgn.border_color
		FROM debit_cards dc
		LEFT JOIN debit_card_status ds ON dc.card_id = ds.card_id
		LEFT JOIN debit_card_details dd ON dc.card_id = dd.card_id
		LEFT JOIN debit_card_design dgn ON dc.card_id = dgn.card_id
		WHERE dc.user_id = ?;
	`

	var debitCards []DebitCardInfo
	if err := r.db.SelectContext(ctx, &debitCards, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []DebitCardInfo{}, nil
		}

		return nil, err
	}

	return debitCards, nil
}
