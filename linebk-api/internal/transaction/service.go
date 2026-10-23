package transaction

import (
	"context"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/zaplog"
)

// Service defines the data access contract for transaction services.
type Service interface {
	GetRecentTransactions(
		ctx context.Context,
		userID string,
	) ([]Transaction, error)
}

type service struct {
	logger    *zaplog.Logger
	cacheRepo CachedRepository
}

// NewService creates a new instance of Transaction Service with the provided dependencies.
func NewService(
	logger *zaplog.Logger,
	cacheRepo Repository,
) Service {
	return &service{
		logger:    logger.Named("service"),
		cacheRepo: cacheRepo,
	}
}

func (s *service) GetRecentTransactions(
	ctx context.Context,
	userID string,
) ([]Transaction, error) {
	transactions, err := s.cacheRepo.GetRecentTransactions(ctx, userID)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
