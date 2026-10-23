package banner_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/internal/banner"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/internal/banner/mocks"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/zaplog"
)

func TestService_GetBanners(t *testing.T) {
	t.Parallel()

	logger := zaplog.NewNop()

	sampleUserID := "user-uuid-12345"
	now := time.Now()

	sampleBanners := []banner.Banner{
		{
			BannerID:    "banner-001",
			UserID:      sampleUserID,
			Title:       sql.NullString{String: "Want some money?", Valid: true},
			Description: sql.NullString{String: "You can start apply 'Clare'", Valid: true},
			Image: sql.NullString{
				String: "https://dummyimage.com/54x54/999/fff",
				Valid:  true,
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	testCases := []struct {
		name          string
		userID        string
		mockSetup     func(mRepo *mocks.MockRepository)
		expectedRes   []banner.Banner
		expectedError error
	}{
		{
			name:   "Success - Banners Retrieved Successfully",
			userID: sampleUserID,
			mockSetup: func(mRepo *mocks.MockRepository) {
				mRepo.EXPECT().
					GetBanners(gomock.Any(), sampleUserID).
					Return(sampleBanners, nil).
					Times(1)
			},
			expectedRes:   sampleBanners,
			expectedError: nil,
		},
		{
			name:   "Success - Empty Banners List",
			userID: sampleUserID,
			mockSetup: func(mRepo *mocks.MockRepository) {
				mRepo.EXPECT().
					GetBanners(gomock.Any(), sampleUserID).
					Return([]banner.Banner{}, nil).
					Times(1)
			},
			expectedRes:   []banner.Banner{},
			expectedError: nil,
		},
		{
			name:   "Failure - Repository Error Propagated",
			userID: sampleUserID,
			mockSetup: func(mRepo *mocks.MockRepository) {
				mRepo.EXPECT().
					GetBanners(gomock.Any(), sampleUserID).
					Return(nil, errors.New("database read error")).
					Times(1)
			},
			expectedRes:   nil,
			expectedError: errors.New("database read error"),
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

			svc := banner.NewService(logger, mockRepo)

			ctx := context.Background()
			res, err := svc.GetBanners(ctx, tc.userID)

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
