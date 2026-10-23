package transaction_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/internal/transaction"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/internal/transaction/mocks"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/zaplog"
)

func TestService_GetRecentTransactions(t *testing.T) {
	t.Parallel()

	logger := zaplog.NewNop()

	sampleUserID := "user-uuid-12345"
	now := time.Now()

	sampleTransactions := []transaction.Transaction{
		{
			TransactionID: "txn-001",
			UserID:        sampleUserID,
			Name:          sql.NullString{String: "Emily", Valid: true},
			Image: sql.NullString{
				String: "https://dummyimage.com/54x54/999/fff",
				Valid:  true,
			},
			IsBank:    false,
			CreatedAt: now,
		},
		{
			TransactionID: "txn-002",
			UserID:        sampleUserID,
			Name:          sql.NullString{String: "Jone Kiersten", Valid: true},
			Image: sql.NullString{
				String: "https://dummyimage.com/54x54/999/fff",
				Valid:  true,
			},
			IsBank:    false,
			CreatedAt: now.Add(-1 * time.Hour),
		},
	}

	testCases := []struct {
		name          string
		userID        string
		mockSetup     func(mRepo *mocks.MockCachedRepository)
		expectedRes   []transaction.Transaction
		expectedError error
	}{
		{
			name:   "Success - Recent Transactions Retrieved",
			userID: sampleUserID,
			mockSetup: func(mRepo *mocks.MockCachedRepository) {
				mRepo.EXPECT().
					GetRecentTransactions(gomock.Any(), sampleUserID).
					Return(sampleTransactions, nil).
					Times(1)
			},
			expectedRes:   sampleTransactions,
			expectedError: nil,
		},
		{
			name:   "Success - Empty Transactions List",
			userID: sampleUserID,
			mockSetup: func(mRepo *mocks.MockCachedRepository) {
				mRepo.EXPECT().
					GetRecentTransactions(gomock.Any(), sampleUserID).
					Return([]transaction.Transaction{}, nil).
					Times(1)
			},
			expectedRes:   []transaction.Transaction{},
			expectedError: nil,
		},
		{
			name:   "Failure - Repository Error Propagated",
			userID: sampleUserID,
			mockSetup: func(mRepo *mocks.MockCachedRepository) {
				mRepo.EXPECT().
					GetRecentTransactions(gomock.Any(), sampleUserID).
					Return(nil, errors.New("database read timeout")).
					Times(1)
			},
			expectedRes:   nil,
			expectedError: errors.New("database read timeout"),
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

			svc := transaction.NewService(logger, mockRepo)

			ctx := context.Background()
			res, err := svc.GetRecentTransactions(ctx, tc.userID)

			if tc.expectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, res)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedRes, res)
				assert.Len(t, res, len(tc.expectedRes))
			}
		})
	}
}
