package auth_test

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"testing"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/internal/auth"
	sessionMocks "github.com/suphawit-wt/linebk-assignment/linebk-api/internal/session/mocks"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/internal/user"
	userMocks "github.com/suphawit-wt/linebk-assignment/linebk-api/internal/user/mocks"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/config"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/security/jwt"
	jwtMocks "github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/security/jwt/mocks"
	passwdMocks "github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/security/passwd/mocks"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/zaplog"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/shared/jwtkey"
)

func TestService_SignInWithPasscode(t *testing.T) {
	t.Parallel()

	logger := zaplog.NewNop()
	sampleUserID := "user-uuid-12345"
	samplePasscode := "123456"
	sampleHash := "$2a$10$hashedpasscodestring"
	clientType := auth.ClientTypeUser
	clientKey := strings.ToLower(clientType.String())

	sampleUser := &user.User{
		UserID:       sampleUserID,
		PasscodeHash: sql.NullString{String: sampleHash, Valid: true},
	}

	jwtKeys := &jwtkey.Keys{
		User: jwtkey.UserTokens{
			Access:  jwtkey.KeyPair{Private: &rsa.PrivateKey{}},
			Refresh: jwtkey.KeyPair{Private: &rsa.PrivateKey{}},
		},
	}

	testCases := []struct {
		name          string
		userID        string
		passcode      string
		mockSetup     func(mUserRepo *userMocks.MockRepository, mPasswd *passwdMocks.MockEncoder, mJwt *jwtMocks.MockService, mSessionRepo *sessionMocks.MockRepository)
		expectedError bool
	}{
		{
			name:     "Success - Sign In With Passcode",
			userID:   sampleUserID,
			passcode: samplePasscode,
			mockSetup: func(mUserRepo *userMocks.MockRepository, mPasswd *passwdMocks.MockEncoder, mJwt *jwtMocks.MockService, mSessionRepo *sessionMocks.MockRepository) {
				mUserRepo.EXPECT().
					GetUserByUserID(gomock.Any(), sampleUserID).
					Return(sampleUser, nil).
					Times(1)

				mPasswd.EXPECT().
					Hash(gomock.Any()).
					Return("dummyhash", nil).
					AnyTimes()

				mPasswd.EXPECT().
					Verify(sampleHash, samplePasscode).
					Return(true, nil).
					Times(1)

				mSessionRepo.EXPECT().
					TrimSessionsByCapacity(gomock.Any(), sampleUserID, gomock.Any(), clientKey, int64(1)).
					Return(nil).
					Times(1)

				mJwt.EXPECT().
					SignRS256(sampleUserID, gomock.Any(), "", gomock.Any(), clientType.String(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("mock.access.token", nil).
					Times(1)

				mJwt.EXPECT().
					SignRS256(sampleUserID, gomock.Any(), "", gomock.Any(), clientType.String(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("mock.refresh.token", nil).
					Times(1)

				mSessionRepo.EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			expectedError: false,
		},
		{
			name:     "Failure - User Not Found",
			userID:   "unknown-user",
			passcode: samplePasscode,
			mockSetup: func(mUserRepo *userMocks.MockRepository, mPasswd *passwdMocks.MockEncoder, mJwt *jwtMocks.MockService, mSessionRepo *sessionMocks.MockRepository) {
				mUserRepo.EXPECT().
					GetUserByUserID(gomock.Any(), "unknown-user").
					Return(nil, user.ErrUserNotFound).
					Times(1)
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mUserRepo := userMocks.NewMockRepository(ctrl)
			mPasswd := passwdMocks.NewMockEncoder(ctrl)
			mJwt := jwtMocks.NewMockService(ctrl)
			mSessionRepo := sessionMocks.NewMockRepository(ctrl)

			if tc.mockSetup != nil {
				tc.mockSetup(mUserRepo, mPasswd, mJwt, mSessionRepo)
			}

			svc := auth.NewService(
				logger,
				&config.AuthConfig{},
				mJwt,
				jwtKeys,
				mPasswd,
				mUserRepo,
				mSessionRepo,
			)

			tokens, err := svc.SignInWithPasscode(
				context.Background(),
				tc.userID,
				tc.passcode,
				"127.0.0.1",
				"Mozilla",
			)

			if tc.expectedError {
				require.Error(t, err)
				assert.Nil(t, tokens)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, tokens)
			}
		})
	}
}

func TestService_SignInWithPassword(t *testing.T) {
	t.Parallel()

	logger := zaplog.NewNop()
	sampleUserID := "user-uuid-12345"
	samplePassword := "securePassword123"
	sampleHash := "$2a$10$hashedpasswordstring"
	clientType := auth.ClientTypeUser
	clientKey := strings.ToLower(clientType.String())

	sampleUser := &user.User{
		UserID:       sampleUserID,
		PasswordHash: sql.NullString{String: sampleHash, Valid: true},
	}

	jwtKeys := &jwtkey.Keys{
		User: jwtkey.UserTokens{
			Access:  jwtkey.KeyPair{Private: &rsa.PrivateKey{}},
			Refresh: jwtkey.KeyPair{Private: &rsa.PrivateKey{}},
		},
	}

	testCases := []struct {
		name          string
		userID        string
		password      string
		mockSetup     func(mUserRepo *userMocks.MockRepository, mPasswd *passwdMocks.MockEncoder, mJwt *jwtMocks.MockService, mSessionRepo *sessionMocks.MockRepository)
		expectedError bool
	}{
		{
			name:     "Success - Sign In With Password",
			userID:   sampleUserID,
			password: samplePassword,
			mockSetup: func(mUserRepo *userMocks.MockRepository, mPasswd *passwdMocks.MockEncoder, mJwt *jwtMocks.MockService, mSessionRepo *sessionMocks.MockRepository) {
				mUserRepo.EXPECT().
					GetUserByUserID(gomock.Any(), sampleUserID).
					Return(sampleUser, nil).
					Times(1)

				mPasswd.EXPECT().
					Verify(sampleHash, samplePassword).
					Return(true, nil).
					Times(1)

				mSessionRepo.EXPECT().
					TrimSessionsByCapacity(gomock.Any(), sampleUserID, gomock.Any(), clientKey, int64(1)).
					Return(nil).
					Times(1)

				mJwt.EXPECT().
					SignRS256(sampleUserID, gomock.Any(), "", gomock.Any(), clientType.String(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("mock.access.token", nil).
					Times(1)

				mJwt.EXPECT().
					SignRS256(sampleUserID, gomock.Any(), "", gomock.Any(), clientType.String(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("mock.refresh.token", nil).
					Times(1)

				mSessionRepo.EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			expectedError: false,
		},
		{
			name:     "Failure - User Not Found",
			userID:   "unknown-user",
			password: samplePassword,
			mockSetup: func(mUserRepo *userMocks.MockRepository, mPasswd *passwdMocks.MockEncoder, mJwt *jwtMocks.MockService, mSessionRepo *sessionMocks.MockRepository) {
				mUserRepo.EXPECT().
					GetUserByUserID(gomock.Any(), "unknown-user").
					Return(nil, user.ErrUserNotFound).
					Times(1)
			},
			expectedError: true,
		},
		{
			name:     "Failure - Invalid Password",
			userID:   sampleUserID,
			password: "wrongPassword",
			mockSetup: func(mUserRepo *userMocks.MockRepository, mPasswd *passwdMocks.MockEncoder, mJwt *jwtMocks.MockService, mSessionRepo *sessionMocks.MockRepository) {
				mUserRepo.EXPECT().
					GetUserByUserID(gomock.Any(), sampleUserID).
					Return(sampleUser, nil).
					Times(1)

				mPasswd.EXPECT().
					Verify(sampleHash, "wrongPassword").
					Return(false, nil).
					Times(1)
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mUserRepo := userMocks.NewMockRepository(ctrl)
			mPasswd := passwdMocks.NewMockEncoder(ctrl)
			mJwt := jwtMocks.NewMockService(ctrl)
			mSessionRepo := sessionMocks.NewMockRepository(ctrl)

			if tc.mockSetup != nil {
				tc.mockSetup(mUserRepo, mPasswd, mJwt, mSessionRepo)
			}

			svc := auth.NewService(
				logger,
				&config.AuthConfig{},
				mJwt,
				jwtKeys,
				mPasswd,
				mUserRepo,
				mSessionRepo,
			)

			tokens, err := svc.SignInWithPassword(
				context.Background(),
				tc.userID,
				tc.password,
				"127.0.0.1",
				"Mozilla",
			)

			if tc.expectedError {
				require.Error(t, err)
				assert.Nil(t, tokens)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, tokens)
				assert.NotEmpty(t, tokens.AccessToken)
				assert.NotEmpty(t, tokens.RefreshToken)
			}
		})
	}
}

func TestService_SignOut(t *testing.T) {
	t.Parallel()

	logger := zaplog.NewNop()

	sampleRefreshToken := "valid.refresh.token.string"
	sampleAccountID := "user-uuid-12345"
	sampleSessionID := "session-uuid-98765"
	samplePublicKey := &rsa.PublicKey{}
	clientType := auth.ClientTypeUser
	clientKey := strings.ToLower(clientType.String())

	jwtKeys := &jwtkey.Keys{
		User: jwtkey.UserTokens{
			Refresh: jwtkey.KeyPair{
				Public: samplePublicKey,
			},
		},
	}

	sampleClaims := &jwt.AuthClaims{
		RegisteredClaims: gojwt.RegisteredClaims{
			Subject: sampleAccountID,
			ID:      sampleSessionID,
		},
	}

	expectedSessionKey := fmtKey(clientKey, sampleAccountID, "session", sampleSessionID)
	expectedAccSessionsKey := fmtKey(clientKey, sampleAccountID, "sessions", "")

	testCases := []struct {
		name          string
		refreshToken  string
		mockSetup     func(mJwt *jwtMocks.MockService, mSessionRepo *sessionMocks.MockRepository)
		expectedError bool
	}{
		{
			name:         "Success - Session Successfully Revoked",
			refreshToken: sampleRefreshToken,
			mockSetup: func(mJwt *jwtMocks.MockService, mSessionRepo *sessionMocks.MockRepository) {
				mJwt.EXPECT().
					VerifyRS256(sampleRefreshToken, samplePublicKey).
					Return(sampleClaims, nil).
					Times(1)

				mSessionRepo.EXPECT().
					Revoke(gomock.Any(), sampleSessionID, expectedSessionKey, expectedAccSessionsKey).
					Return(true, nil).
					Times(1)
			},
			expectedError: false,
		},
		{
			name:         "Success - Session Already Expired or Revoked (Idempotent)",
			refreshToken: sampleRefreshToken,
			mockSetup: func(mJwt *jwtMocks.MockService, mSessionRepo *sessionMocks.MockRepository) {
				mJwt.EXPECT().
					VerifyRS256(sampleRefreshToken, samplePublicKey).
					Return(sampleClaims, nil).
					Times(1)

				mSessionRepo.EXPECT().
					Revoke(gomock.Any(), sampleSessionID, expectedSessionKey, expectedAccSessionsKey).
					Return(false, nil).
					Times(1)
			},
			expectedError: false,
		},
		{
			name:         "Failure - Invalid Refresh Token Signature",
			refreshToken: "invalid.token",
			mockSetup: func(mJwt *jwtMocks.MockService, mSessionRepo *sessionMocks.MockRepository) {
				mJwt.EXPECT().
					VerifyRS256("invalid.token", samplePublicKey).
					Return(nil, errors.New("signature mismatch")).
					Times(1)
			},
			expectedError: true,
		},
		{
			name:         "Failure - Repository Revoke Error Propagated",
			refreshToken: sampleRefreshToken,
			mockSetup: func(mJwt *jwtMocks.MockService, mSessionRepo *sessionMocks.MockRepository) {
				mJwt.EXPECT().
					VerifyRS256(sampleRefreshToken, samplePublicKey).
					Return(sampleClaims, nil).
					Times(1)

				mSessionRepo.EXPECT().
					Revoke(gomock.Any(), sampleSessionID, expectedSessionKey, expectedAccSessionsKey).
					Return(false, errors.New("redis down")).
					Times(1)
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockJwtSvc := jwtMocks.NewMockService(ctrl)
			mockSessionRepo := sessionMocks.NewMockRepository(ctrl)

			if tc.mockSetup != nil {
				tc.mockSetup(mockJwtSvc, mockSessionRepo)
			}

			svc := auth.NewService(
				logger,
				&config.AuthConfig{},
				mockJwtSvc,
				jwtKeys,
				nil,
				nil,
				mockSessionRepo,
			)

			ctx := context.Background()
			err := svc.SignOut(ctx, tc.refreshToken)

			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func fmtKey(clientType, accID, kind, sessionID string) string {
	if kind == "sessions" {
		return fmt.Sprintf("%s:%s:%s", clientType, accID, kind)
	}
	return fmt.Sprintf("%s:%s:%s:%s", clientType, accID, kind, sessionID)
}
