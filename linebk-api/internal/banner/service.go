package banner

import (
	"context"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/zaplog"
)

// Service defines the data access contract for banner services.
type Service interface {
	GetBanners(
		ctx context.Context,
		userID string,
	) ([]Banner, error)
}

type service struct {
	logger     *zaplog.Logger
	bannerRepo Repository
}

// NewService creates a new instance of Banner Service with the provided dependencies.
func NewService(
	logger *zaplog.Logger,
	bannerRepo Repository,
) Service {
	return &service{
		logger:     logger.Named("service"),
		bannerRepo: bannerRepo,
	}
}

func (s *service) GetBanners(
	ctx context.Context,
	userID string,
) ([]Banner, error) {
	banners, err := s.bannerRepo.GetBanners(ctx, userID)
	if err != nil {
		return nil, err
	}

	return banners, nil
}
