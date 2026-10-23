package auth

import (
	"context"
	"crypto/rsa"
	"crypto/subtle"
	"errors"
	"fmt"
	"strings"
	"time"
	"uuid"

	"go.uber.org/zap"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/internal/session"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/internal/user"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/errs"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/security/jwt"
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

func (s *service) validateUserPassword(
	u *user.User,
	password, identifier, ipAddress, userAgent string,
	logger *zaplog.Logger,
) error {
	if !u.PasswordHash.Valid {
		logger.Warn(
			"user sign-in failed: account has no password set",
			zap.String("user_id", u.UserID),
			zap.String("identifier", identifier),
			zap.String("ip_address", ipAddress),
		)
		return errs.NewUnauthorizedError("Invalid user id or password.")
	}

	verified, err := s.passwordEncoder.Verify(u.PasswordHash.String, password)
	if err != nil {
		logger.Error(
			"failed to execute password verification",
			zap.String("user_id", u.UserID),
			zap.Error(err),
		)
		return fmt.Errorf("auth_service.SignInWithPassword: %w", err)
	}

	if !verified {
		logger.Warn(
			"user sign-in failed: invalid password provided",
			zap.String("user_id", u.UserID),
			zap.String("identifier", identifier),
			zap.String("ip_address", ipAddress),
			zap.String("user_agent", userAgent),
		)
		return errs.NewUnauthorizedError("Invalid user id or password.")
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

	userID := u.UserID
	clientKey := strings.ToLower(clientType.String())
	sessionID := uuid.New().String()
	sessionKey := fmt.Sprintf("%s:%s:session:%s", clientKey, userID, sessionID)
	accSessionsKey := fmt.Sprintf("%s:%s:sessions", clientKey, userID)

	if err := s.sessionRepo.TrimSessionsByCapacity(
		ctx,
		userID,
		accSessionsKey,
		clientKey,
		maxSessions,
	); err != nil {
		logger.Warn(
			"failed to trim session capacity",
			zap.String("user_id", userID),
			zap.Error(err),
		)
	}

	userRoles := []string{"USER"}

	accessToken, err := s.jwtSvc.SignRS256(
		userID,
		sessionID,
		"",
		userRoles,
		clientType.String(),
		s.authCfg.AccessTTL,
		s.authCfg.Audience,
		keys.AccessPrivateKey,
	)
	if err != nil {
		logger.Error("failed to sign access token", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}

	refreshToken, err := s.jwtSvc.SignRS256(
		userID,
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
			zap.String("user_id", userID),
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
			zap.String("user_id", userID),
			zap.String("session_id", sessionID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("authService.createSessionAndTokens: %w", err)
	}

	logger.Info("authenticated successfully",
		zap.String("user_id", userID),
		zap.String("session_id", sessionID),
		zap.String("client_type", clientType.String()),
	)

	return &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) validateAndFetchSession(
	ctx context.Context,
	sessionKey, accSessionsKey, userID, sessionID, refreshToken string,
	logger *zaplog.Logger,
) (*session.Metadata, error) {
	sessionMeta, err := s.sessionRepo.GetMetadata(ctx, sessionKey)
	if err != nil {
		if errors.Is(err, session.ErrTokenNotFound) {
			logger.Warn("refresh attempt on non-existent or revoked session",
				zap.String("user_id", userID),
				zap.String("session_id", sessionID),
			)
			return nil, errs.NewUnauthorizedError("Session has expired or been revoked.")
		}

		logger.Error("failed to retrieve session from repository",
			zap.String("user_id", userID),
			zap.String("session_id", sessionID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("authService.RotateTokens: %w", err)
	}

	incomingHash := secret.HashToken(refreshToken)
	if sessionMeta.TokenHash != "" &&
		subtle.ConstantTimeCompare([]byte(sessionMeta.TokenHash), []byte(incomingHash)) != 1 {
		logger.Warn("token reuse or mismatch detected",
			zap.String("user_id", userID),
			zap.String("session_id", sessionID),
		)

		_, _ = s.sessionRepo.Revoke(ctx, sessionID, sessionKey, accSessionsKey)
		return nil, errs.NewUnauthorizedError("Invalid refresh token.")
	}

	return sessionMeta, nil
}

func (s *service) issueRotatedTokens(
	ctx context.Context,
	claims *jwt.AuthClaims,
	sessionMeta *session.Metadata,
	sessionID, sessionKey, accSessionsKey string,
	keys RotateKeys,
	logger *zaplog.Logger,
) (*Token, error) {
	clientType := ClientTypeUser
	userID := claims.Subject

	newAccessToken, err := s.jwtSvc.SignRS256(
		userID,
		sessionID,
		claims.Email,
		claims.Roles,
		clientType.String(),
		s.authCfg.AccessTTL,
		s.authCfg.Audience,
		keys.AccessPrivateKey,
	)
	if err != nil {
		logger.Error("failed to sign new access token during refresh",
			zap.String("user_id", userID),
			zap.String("session_id", sessionID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("authService.RotateTokens: %w", err)
	}

	newRefreshToken, err := s.jwtSvc.SignRS256(
		userID,
		sessionID,
		claims.Email,
		claims.Roles,
		clientType.String(),
		s.authCfg.RefreshTTL,
		s.authCfg.Audience,
		keys.RefreshPrivateKey,
	)
	if err != nil {
		logger.Error("failed to sign new refresh token during refresh",
			zap.String("user_id", userID),
			zap.String("session_id", sessionID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("authService.RotateTokens: %w", err)
	}

	sessionMeta.TokenHash = secret.HashToken(newRefreshToken)
	sessionMeta.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	if err := s.sessionRepo.Create(
		ctx,
		sessionID,
		sessionKey,
		accSessionsKey,
		sessionMeta,
		s.authCfg.RefreshTTL,
	); err != nil {
		logger.Error("failed to create session during rotate tokens",
			zap.String("user_id", userID),
			zap.String("session_id", sessionID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("authService.RotateTokens: %w", err)
	}

	logger.Info("tokens refreshed successfully",
		zap.String("user_id", userID),
		zap.String("session_id", sessionID),
	)

	return &Token{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
