package card_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/internal/card"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/internal/card/mocks"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/zaplog"
)

func TestService_GetDebitCards(t *testing.T) {
	t.Parallel()

	logger := zaplog.NewNop()

	sampleUserID := "user-uuid-12345"

	sampleDebitCards := []card.DebitCardInfo{
		{
			CardID:      "card-001",
			UserID:      sampleUserID,
			Name:        sql.NullString{String: "My Salary", Valid: true},
			Status:      sql.NullString{String: "in-progress", Valid: true},
			Issuer:      sql.NullString{String: "TestLab", Valid: true},
			Number:      sql.NullString{String: "", Valid: false},
			Color:       sql.NullString{String: "#00a1e2", Valid: true},
			BorderColor: sql.NullString{String: "", Valid: false},
		},
		{
			CardID:      "card-002",
			UserID:      sampleUserID,
			Name:        sql.NullString{String: "My Debit card", Valid: true},
			Status:      sql.NullString{String: "Active", Valid: true},
			Issuer:      sql.NullString{String: "TestLab", Valid: true},
			Number:      sql.NullString{String: "9440 7841 2222 3115", Valid: true},
			Color:       sql.NullString{String: "#ffffff", Valid: true},
			BorderColor: sql.NullString{String: "#f2f3f7", Valid: true},
		},
	}

	testCases := []struct {
		name          string
		userID        string
		mockSetup     func(mRepo *mocks.MockCachedRepository)
		expectedRes   []card.DebitCardInfo
		expectedError error
	}{
		{
			name:   "Success - Debit Cards Retrieved Successfully",
			userID: sampleUserID,
			mockSetup: func(mRepo *mocks.MockCachedRepository) {
				mRepo.EXPECT().
					GetDebitCards(gomock.Any(), sampleUserID).
					Return(sampleDebitCards, nil).
					Times(1)
			},
			expectedRes:   sampleDebitCards,
			expectedError: nil,
		},
		{
			name:   "Success - Empty Debit Cards List",
			userID: sampleUserID,
			mockSetup: func(mRepo *mocks.MockCachedRepository) {
				mRepo.EXPECT().
					GetDebitCards(gomock.Any(), sampleUserID).
					Return([]card.DebitCardInfo{}, nil).
					Times(1)
			},
			expectedRes:   []card.DebitCardInfo{},
			expectedError: nil,
		},
		{
			name:   "Failure - Repository Error Propagated",
			userID: sampleUserID,
			mockSetup: func(mRepo *mocks.MockCachedRepository) {
				mRepo.EXPECT().
					GetDebitCards(gomock.Any(), sampleUserID).
					Return(nil, errors.New("database connection failed")).
					Times(1)
			},
			expectedRes:   nil,
			expectedError: errors.New("database connection failed"),
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

			svc := card.NewService(logger, mockRepo)

			ctx := context.Background()
			res, err := svc.GetDebitCards(ctx, tc.userID)

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
