package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"strings"

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

// RotateKeys holds the cryptographic RSA keys required during a key rotation phase.
type RotateKeys struct {
	AccessPrivateKey  *rsa.PrivateKey
	RefreshPrivateKey *rsa.PrivateKey
	RefreshPublicKey  *rsa.PublicKey
}

// Service defines the data access contract for authentication services.
type Service interface {
	SignInWithPasscode(
		ctx context.Context,
		userID, passcode, ipAddress, userAgent string,
	) (*Token, error)
	SignInWithPassword(
		ctx context.Context,
		userID, passcode, ipAddress, userAgent string,
	) (*Token, error)
	SignOut(
		ctx context.Context,
		refreshToken string,
	) error
	RotateTokens(
		ctx context.Context,
		refreshToken string,
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

func (s *service) SignInWithPassword(
	ctx context.Context,
	userID, password, ipAddress, userAgent string,
) (*Token, error) {
	logger := s.logger.Named("signin_with_password")

	u, err := s.userRepo.GetUserByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			logger.Warn("user sign-in failed: credential not found",
				zap.String("identifier", userID),
				zap.String("event", "auth.signin.password.failed"),
			)

			return nil, errs.NewUnauthorizedError("Invalid user id or password.")
		}

		return nil, err
	}

	if err := s.validateUserPassword(
		u,
		password,
		userID,
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

func (s *service) SignOut(
	ctx context.Context,
	refreshToken string,
) error {
	logger := s.logger.Named("sign_out")

	claims, err := s.jwtSvc.VerifyRS256(refreshToken, s.jwtKeys.User.Refresh.Public)
	if err != nil {
		return errs.NewUnauthorizedError("Invalid Refresh Token!")
	}

	accID := claims.Subject
	sessionID := claims.ID
	clientType := ClientTypeUser
	clientKey := strings.ToLower(clientType.String())
	sessionKey := fmt.Sprintf("%s:%s:session:%s", clientKey, accID, sessionID)
	accSessionsKey := fmt.Sprintf("%s:%s:sessions", clientKey, accID)

	deleted, err := s.sessionRepo.Revoke(ctx, sessionID, sessionKey, accSessionsKey)
	if err != nil {
		logger.Error("failed to revoke device session in storage",
			zap.String("account_id", accID),
			zap.String("session_id", sessionID),
			zap.Error(err),
		)

		return err
	}

	if !deleted {
		logger.Info("session already revoked or expired",
			zap.String("account_id", accID),
			zap.String("session_id", sessionID),
		)

		return nil
	}

	logger.Info("device session revoked successfully",
		zap.String("account_id", accID),
		zap.String("session_id", sessionID),
	)

	return nil
}

func (s *service) RotateTokens(
	ctx context.Context,
	refreshToken string,
) (*Token, error) {
	logger := s.logger.Named("rotate_tokens")

	keys := RotateKeys{
		AccessPrivateKey:  s.jwtKeys.User.Access.Private,
		RefreshPrivateKey: s.jwtKeys.User.Refresh.Private,
		RefreshPublicKey:  s.jwtKeys.User.Refresh.Public,
	}

	claims, err := s.jwtSvc.VerifyRS256(refreshToken, keys.RefreshPublicKey)
	if err != nil {
		logger.Warn("refresh token validation failed", zap.Error(err))
		return nil, errs.NewUnauthorizedError("Invalid or expired refresh token.")
	}

	accID := claims.Subject
	sessionID := claims.ID
	clientType := "user"
	clientKey := strings.ToLower(clientType)
	sessionKey := fmt.Sprintf("%s:%s:session:%s", clientKey, accID, sessionID)
	accSessionsKey := fmt.Sprintf("%s:%s:sessions", clientKey, accID)

	sessionMeta, err := s.validateAndFetchSession(
		ctx,
		sessionKey,
		accSessionsKey,
		accID,
		sessionID,
		refreshToken,
		logger,
	)
	if err != nil {
		return nil, err
	}

	return s.issueRotatedTokens(
		ctx,
		claims,
		sessionMeta,
		sessionID,
		sessionKey,
		accSessionsKey,
		keys,
		logger,
	)
}
