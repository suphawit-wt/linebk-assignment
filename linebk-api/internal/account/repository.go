package account

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

// Info aggregates account, balance, details, and flag data for query results.
type Info struct {
	AccountID     string          `db:"account_id"`
	Type          string          `db:"type"`
	Amount        decimal.Decimal `db:"amount"`
	Currency      string          `db:"currency"`
	AccountNumber string          `db:"account_number"`
	Issuer        string          `db:"issuer"`
	Progress      *uint8          `db:"progress"`
	Color         sql.NullString  `db:"color"`
	IsMainAccount bool            `db:"is_main_account"`
	RawFlags      sql.NullString  `db:"raw_flags"`
}

// Repository defines data access methods for managing accounts.
type Repository interface {
	GetAccountsByUserID(ctx context.Context, userID string) ([]Info, error)
}

type repositoryDB struct {
	db *sqlx.DB
}

// NewRepositoryDB creates a new instance of account repositoryDB with the provided database connection pool.
func NewRepositoryDB(db *sqlx.DB) Repository {
	return &repositoryDB{db: db}
}

func (r *repositoryDB) GetAccountsByUserID(ctx context.Context, userID string) ([]Info, error) {
	query := `
		SELECT 
			a.account_id,
			a.type,
			ab.amount,
			a.currency,
			a.account_number,
			a.issuer,
			ad.color,
			ad.is_main_account,
			ad.progress,
			GROUP_CONCAT(af.flag_value SEPARATOR ',') AS raw_flags
		FROM accounts a
		LEFT JOIN account_balances ab ON a.account_id = ab.account_id
		LEFT JOIN account_details ad ON a.account_id = ad.account_id
		LEFT JOIN account_flags af ON a.account_id = af.account_id
		WHERE a.user_id = ?
		GROUP BY a.account_id
	`

	var accountsInfo []Info
	if err := r.db.SelectContext(ctx, &accountsInfo, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Info{}, nil
		}

		return nil, err
	}

	return accountsInfo, nil
}
