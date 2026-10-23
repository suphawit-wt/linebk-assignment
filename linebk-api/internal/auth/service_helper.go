package auth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"strings"
	"time"
	"uuid"

	"go.uber.org/zap"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/internal/session"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/internal/user"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/errs"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/security/secret"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/zaplog"
)

// KeyPair holds the private signing keys for a client type.
type KeyPair struct {
	AccessPrivateKey  *rsa.PrivateKey
	RefreshPrivateKey *rsa.PrivateKey
}

func (s *service) validateUserPasscode(
	u *user.User,
	passcode, ipAddress, userAgent string,
	logger *zaplog.Logger,
) error {
	if !u.PasscodeHash.Valid {
		logger.Warn(
			"user sign-in failed: account has no passcode set",
			zap.String("user_id", u.UserID),
			zap.String("ip_address", ipAddress),
		)
		return errs.NewUnauthorizedError("Invalid user id or passcode.")
	}

	verified, err := s.passwordEncoder.Verify(u.PasscodeHash.String, passcode)
	if err != nil {
		logger.Error(
			"failed to execute passcode verification",
			zap.String("user_id", u.UserID),
			zap.Error(err),
		)
		return fmt.Errorf("auth_service.SignInWithPasscode: %w", err)
	}

	if !verified {
		logger.Warn(
			"user sign-in failed: invalid passcode provided",
			zap.String("user_id", u.UserID),
			zap.String("ip_address", ipAddress),
			zap.String("user_agent", userAgent),
		)
		return errs.NewUnauthorizedError("Invalid user id or passcode.")
	}

	return nil
}

func (s *service) createSessionAndTokens(
	ctx context.Context,
	u *user.User,
	keys KeyPair,
	maxSessions int64,
	ipAddress, userAgent string,
	logger *zaplog.Logger,
) (*Token, error) {
	clientType := ClientTypeUser

	accID := u.UserID
	clientKey := strings.ToLower(clientType.String())
	sessionID := uuid.New().String()
	sessionKey := fmt.Sprintf("%s:%s:session:%s", clientKey, accID, sessionID)
	accSessionsKey := fmt.Sprintf("%s:%s:sessions", clientKey, accID)

	if err := s.sessionRepo.TrimSessionsByCapacity(
		ctx,
		accID,
		accSessionsKey,
		clientKey,
		maxSessions,
	); err != nil {
		logger.Warn(
			"failed to trim session capacity",
			zap.String("account_id", accID),
			zap.Error(err),
		)
	}

	userRoles := []string{"USER"}

	accessToken, err := s.jwtSvc.SignRS256(
		accID,
		sessionID,
		"",
		userRoles,
		clientType.String(),
		s.authCfg.AccessTTL,
		s.authCfg.Audience,
		keys.AccessPrivateKey,
	)
	if err != nil {
		logger.Error("failed to sign access token", zap.String("account_id", accID), zap.Error(err))
		return nil, err
	}

	refreshToken, err := s.jwtSvc.SignRS256(
		accID,
		sessionID,
		"",
		userRoles,
		clientType.String(),
		s.authCfg.RefreshTTL,
		s.authCfg.Audience,
		keys.RefreshPrivateKey,
	)
	if err != nil {
		logger.Error(
			"failed to sign refresh token",
			zap.String("account_id", accID),
			zap.Error(err),
		)
		return nil, err
	}

	sessionMeta := session.Metadata{
		SessionID: sessionID,
		TokenHash: secret.HashToken(refreshToken),
		IPAddress: ipAddress,
		UserAgent: userAgent,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	if err := s.sessionRepo.Create(
		ctx,
		sessionID,
		sessionKey,
		accSessionsKey,
		&sessionMeta,
		s.authCfg.RefreshTTL,
	); err != nil {
		logger.Error("failed to store device session in repository",
			zap.String("account_id", accID),
			zap.String("session_id", sessionID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("authService.createSessionAndTokens: %w", err)
	}

	logger.Info("authenticated successfully",
		zap.String("account_id", accID),
		zap.String("session_id", sessionID),
		zap.String("client_type", clientType.String()),
	)

	return &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
