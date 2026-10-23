package account

import (
	"context"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/zaplog"
)

// Service defines the data access contract for account services.
type Service interface {
	GetAccountsByUserID(
		ctx context.Context,
		userID string,
	) ([]Info, error)
}

type service struct {
	logger    *zaplog.Logger
	cacheRepo CachedRepository
}

// NewService creates a new instance of Account Service with the provided dependencies.
func NewService(
	logger *zaplog.Logger,
	cacheRepo Repository,
) Service {
	return &service{
		logger:    logger.Named("service"),
		cacheRepo: cacheRepo,
	}
}

func (s *service) GetAccountsByUserID(
	ctx context.Context,
	userID string,
) ([]Info, error) {
	accountsInfo, err := s.cacheRepo.GetAccountsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return accountsInfo, nil
}
