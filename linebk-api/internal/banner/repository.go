package banner

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

// Repository defines data access methods for managing banner
type Repository interface {
	GetBanners(ctx context.Context, userID string) ([]Banner, error)
}

type repositoryDB struct {
	db *sqlx.DB
}

// NewRepositoryDB creates a new instance of banner repositoryDB with the provided database connection pool.
func NewRepositoryDB(db *sqlx.DB) Repository {
	return &repositoryDB{db: db}
}

func (r *repositoryDB) GetBanners(
	ctx context.Context,
	userID string,
) ([]Banner, error) {
	query := `
		SELECT 
			title,
			description,
			image
		FROM banners
		WHERE user_id = ?;
	`

	var banners []Banner
	if err := r.db.SelectContext(ctx, &banners, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Banner{}, nil
		}

		return nil, err
	}

	return banners, nil
}
