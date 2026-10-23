package card

import (
	"context"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/zaplog"
)

// Service defines the data access contract for card services.
type Service interface {
	GetDebitCards(
		ctx context.Context,
		userID string,
	) ([]DebitCardInfo, error)
}

type service struct {
	logger    *zaplog.Logger
	cacheRepo CachedRepository
}

// NewService creates a new instance of Card Service with the provided dependencies.
func NewService(
	logger *zaplog.Logger,
	cacheRepo Repository,
) Service {
	return &service{
		logger:    logger.Named("service"),
		cacheRepo: cacheRepo,
	}
}

func (s *service) GetDebitCards(
	ctx context.Context,
	userID string,
) ([]DebitCardInfo, error) {
	debitCards, err := s.cacheRepo.GetDebitCards(ctx, userID)
	if err != nil {
		return nil, err
	}

	return debitCards, nil
}
