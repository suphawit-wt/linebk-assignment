package transaction

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

// Repository defines data access methods for managing transaction
type Repository interface {
	GetRecentTransactions(ctx context.Context, userID string) ([]Transaction, error)
}

type repositoryDB struct {
	db *sqlx.DB
}

// NewRepositoryDB creates a new instance of transaction repositoryDB with the provided database connection pool.
func NewRepositoryDB(db *sqlx.DB) Repository {
	return &repositoryDB{db: db}
}

func (r *repositoryDB) GetRecentTransactions(
	ctx context.Context,
	userID string,
) ([]Transaction, error) {
	query := `
		SELECT 
			name,
			image,
			isBank,
			created_at
		FROM transactions
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT 10;
	`

	var transactions []Transaction
	if err := r.db.SelectContext(ctx, &transactions, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Transaction{}, nil
		}

		return nil, err
	}

	return transactions, nil
}
