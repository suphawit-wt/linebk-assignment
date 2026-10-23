package auth

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/internal/session"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/internal/user"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/config"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/errs"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/security/jwt"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/security/passwd"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/zaplog"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/shared/jwtkey"
)

// Token represents the access and refresh tokens for authentication
type Token struct {
	AccessToken  string
	RefreshToken string
}

// Service defines the data access contract for authentication services.
type Service interface {
	SignInWithPasscode(
		ctx context.Context,
		userID, passcode, ipAddress, userAgent string,
	) (*Token, error)
}

type service struct {
	logger          *zaplog.Logger
	authCfg         *config.AuthConfig
	jwtSvc          jwt.Service
	jwtKeys         *jwtkey.Keys
	passwordEncoder passwd.Encoder
	userRepo        user.Repository
	sessionRepo     session.Repository
}

// NewService creates a new instance of Auth Service with the provided dependencies.
func NewService(
	logger *zaplog.Logger,
	authCfg *config.AuthConfig,
	jwtSvc jwt.Service,
	jwtKeys *jwtkey.Keys,
	passwordEncoder passwd.Encoder,
	userRepo user.Repository,
	sessionRepo session.Repository,
) Service {
	return &service{
		logger:          logger.Named("service"),
		authCfg:         authCfg,
		jwtSvc:          jwtSvc,
		jwtKeys:         jwtKeys,
		passwordEncoder: passwordEncoder,
		userRepo:        userRepo,
		sessionRepo:     sessionRepo,
	}
}

func (s *service) SignInWithPasscode(
	ctx context.Context,
	userID, passcode, ipAddress, userAgent string,
) (*Token, error) {
	logger := s.logger.Named("signin_with_passcode")

	u, err := s.userRepo.GetUserByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			logger.Warn("user sign-in failed: credential not found",
				zap.String("identifier", userID),
				zap.String("event", "auth.signin.passcode.failed"),
			)

			return nil, errs.NewUnauthorizedError("Invalid user id or passcode.")
		}

		return nil, err
	}

	if err := s.validateUserPasscode(
		u,
		passcode,
		ipAddress,
		userAgent,
		logger,
	); err != nil {
		return nil, err
	}

	keys := KeyPair{
		AccessPrivateKey:  s.jwtKeys.User.Access.Private,
		RefreshPrivateKey: s.jwtKeys.User.Refresh.Private,
	}

	return s.createSessionAndTokens(
		ctx,
		u,
		keys,
		1,
		ipAddress,
		userAgent,
		logger,
	)
}
