package user_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/internal/user"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/internal/user/mocks"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/zaplog"
)

func TestService_GetUserInfo(t *testing.T) {
	t.Parallel()

	logger := zaplog.NewNop()

	sampleUserID := "user-uuid-12345"
	sampleUserInfo := &user.Info{
		Name:            sql.NullString{String: "Clare", Valid: true},
		GreetingMessage: sql.NullString{String: "Have a nice day Clare", Valid: true},
	}

	testCases := []struct {
		name          string
		userID        string
		mockSetup     func(mRepo *mocks.MockRepository)
		expectedRes   *user.Info
		expectedError bool
	}{
		{
			name:   "Success - User Info Retrieved Successfully",
			userID: sampleUserID,
			mockSetup: func(mRepo *mocks.MockRepository) {
				mRepo.EXPECT().
					GetUserInfo(gomock.Any(), sampleUserID).
					Return(sampleUserInfo, nil).
					Times(1)
			},
			expectedRes:   sampleUserInfo,
			expectedError: false,
		},
		{
			name:   "Failure - User Not Found Handled and Translated",
			userID: sampleUserID,
			mockSetup: func(mRepo *mocks.MockRepository) {
				mRepo.EXPECT().
					GetUserInfo(gomock.Any(), sampleUserID).
					Return(nil, user.ErrUserNotFound).
					// Triggers business error mapping in service layer[cite: 18, 20]
					Times(1)
			},
			expectedRes:   nil,
			expectedError: true,
		},
		{
			name:   "Failure - General Repository Error Propagated",
			userID: sampleUserID,
			mockSetup: func(mRepo *mocks.MockRepository) {
				mRepo.EXPECT().
					GetUserInfo(gomock.Any(), sampleUserID).
					Return(nil, errors.New("database connection failed")).
					Times(1)
			},
			expectedRes:   nil,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockRepository(ctrl)
			if tc.mockSetup != nil {
				tc.mockSetup(mockRepo)
			}

			svc := user.NewService(logger, mockRepo)

			ctx := context.Background()
			res, err := svc.GetUserInfo(ctx, tc.userID)

			if tc.expectedError {
				require.Error(t, err)
				assert.Nil(t, res)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedRes, res)
			}
		})
	}
}
