package repository

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"yt-video-transcriptor/models"
	"yt-video-transcriptor/models/repository/mocks"
)

func Test_formatQuery(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected string
	}{
		{
			name:     "With \\n",
			query:    "\nCREATE DATABASE example;",
			expected: " CREATE DATABASE example;",
		},
		{
			name:     "With \\t",
			query:    "\tCREATE DATABASE example;",
			expected: "CREATE DATABASE example;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := formatQuery(tt.query)

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestYTVideoRepository_Create(t *testing.T) {
	type MockHandler func(*mocks.MockRepository, context.Context, models.YTVideo)

	testCases := []struct {
		name    string
		request models.VideoRequest
		videos  models.YTVideo
		mock    MockHandler
		wantErr bool
	}{
		{
			name: "successful create",
			videos: models.YTVideo{
				VideoRequest: models.VideoRequest{
					VideoID:  "video_id_1",
					Language: "en",
				},
				Title: "Hello, World!",
			},

			mock: func(
				mockRepo *mocks.MockRepository,
				ctx context.Context,
				videos models.YTVideo) {

				mockRepo.EXPECT().Create(ctx, videos).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "create error",
			videos: models.YTVideo{
				VideoRequest: models.VideoRequest{
					VideoID:  "video_id_2",
					Language: "en",
				},
				Title: "Hello, World!",
			},

			mock: func(
				mockRepo *mocks.MockRepository,
				ctx context.Context,
				videos models.YTVideo) {

				mockRepo.EXPECT().Create(ctx, videos).Return(errors.New("create error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockRepository(ctrl)

			tc.mock(mockRepo, ctx, tc.videos)

			err := mockRepo.Create(ctx, tc.videos)

			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}
