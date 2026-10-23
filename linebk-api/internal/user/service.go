package user

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/errs"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/zaplog"
)

// Service defines the data access contract for user services.
type Service interface {
	GetUserInfo(
		ctx context.Context,
		userID string,
	) (*Info, error)
}

type service struct {
	logger   *zaplog.Logger
	userRepo Repository
}

// NewService creates a new instance of User Service with the provided dependencies.
func NewService(
	logger *zaplog.Logger,
	userRepo Repository,
) Service {
	return &service{
		logger:   logger.Named("service"),
		userRepo: userRepo,
	}
}

func (s *service) GetUserInfo(
	ctx context.Context,
	userID string,
) (*Info, error) {
	logger := s.logger.Named("get_user_info")

	userInfo, err := s.userRepo.GetUserInfo(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			logger.Warn("user info query failed: user not found",
				zap.String("user_id", userID),
				zap.String("event", "user.get.info.failed"),
			)

			return nil, errs.NewNotFoundError("user not found.")
		}

		return nil, err
	}

	return userInfo, nil
}
