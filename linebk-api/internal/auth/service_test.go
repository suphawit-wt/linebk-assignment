package auth_test

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/internal/auth"
	sessionMocks "github.com/suphawit-wt/linebk-assignment/linebk-api/internal/session/mocks"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/internal/user"
	userMocks "github.com/suphawit-wt/linebk-assignment/linebk-api/internal/user/mocks"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/config"
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
