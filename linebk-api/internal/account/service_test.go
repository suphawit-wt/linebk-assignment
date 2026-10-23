package account_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/internal/account"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/internal/account/mocks"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/zaplog"
)

func TestService_GetAccountsByUserID(t *testing.T) {
	t.Parallel()

	logger := zaplog.NewNop()

	sampleUserID := "user-uuid-12345"
	var sampleProgress uint8 = 100
	sampleAccounts := []account.Info{
		{
			AccountID:     "acc-001",
			Type:          "saving-account",
			Amount:        decimal.NewFromFloat(5000.50),
			Currency:      "THB",
			AccountNumber: "568-2-81740-9",
			Issuer:        "TestLab",
			Progress:      &sampleProgress,
			IsMainAccount: true,
		},
	}

	testCases := []struct {
		name          string
		userID        string
		mockSetup     func(mRepo *mocks.MockCachedRepository)
		expectedRes   []account.Info
		expectedError error
	}{
		{
			name:   "Success - Accounts Found",
			userID: sampleUserID,
			mockSetup: func(mRepo *mocks.MockCachedRepository) {
				mRepo.EXPECT().
					GetAccountsByUserID(gomock.Any(), sampleUserID).
					Return(sampleAccounts, nil).
					Times(1)
			},
			expectedRes:   sampleAccounts,
			expectedError: nil,
		},
		{
			name:   "Success - Empty Accounts List",
			userID: sampleUserID,
			mockSetup: func(mRepo *mocks.MockCachedRepository) {
				mRepo.EXPECT().
					GetAccountsByUserID(gomock.Any(), sampleUserID).
					Return([]account.Info{}, nil).
					Times(1)
			},
			expectedRes:   []account.Info{},
			expectedError: nil,
		},
		{
			name:   "Failure - Repository Error Propagated",
			userID: sampleUserID,
			mockSetup: func(mRepo *mocks.MockCachedRepository) {
				mRepo.EXPECT().
					GetAccountsByUserID(gomock.Any(), sampleUserID).
					Return(nil, errors.New("database connection timeout")).
					Times(1)
			},
			expectedRes:   nil,
			expectedError: errors.New("database connection timeout"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockCachedRepository(ctrl)
			if tc.mockSetup != nil {
				tc.mockSetup(mockRepo)
			}

			svc := account.NewService(logger, mockRepo)

			ctx := context.Background()
			res, err := svc.GetAccountsByUserID(ctx, tc.userID)

			// Assertions using testify
			if tc.expectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, res)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedRes, res)
			}
		})
	}
}
