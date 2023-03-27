package routes

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"yt-video-transcriptor/models"
)

func Test_isValidVideoRequest(t *testing.T) {
	tests := []struct {
		name        string
		data        models.VideoRequest
		expectedRes bool
	}{
		{
			name: "Valid input",
			data: models.VideoRequest{
				VideoID:  "12345678901",
				Language: "ru",
			},
			expectedRes: true,
		},
		{
			name: "Invalid video id (len: 10)",
			data: models.VideoRequest{
				VideoID:  "1234567890",
				Language: "ru",
			},
			expectedRes: false,
		},
		{
			name: "Invalid language (with numbers) (len: 2)",
			data: models.VideoRequest{
				VideoID:  "12345678901",
				Language: "r1",
			},
			expectedRes: false,
		},
		{
			name: "Invalid language (with numbers) (len: 1)",
			data: models.VideoRequest{
				VideoID:  "12345678901",
				Language: "1",
			},
			expectedRes: false,
		},
		{
			name: "Invalid video id (len: 0)",
			data: models.VideoRequest{
				VideoID:  "",
				Language: "ru",
			},
			expectedRes: false,
		},
		{
			name: "Invalid language (len: 0)",
			data: models.VideoRequest{
				VideoID:  "12345678901",
				Language: "",
			},
			expectedRes: false,
		},
		{
			name: "Invalid  video id & language (len: 0)",
			data: models.VideoRequest{
				VideoID:  "",
				Language: "",
			},
			expectedRes: false,
		},
	}

	for _, testcase := range tests {
		var (
			request  = testcase.data
			expected = testcase.expectedRes
		)

		isNotValid, _ := isValidVideoRequest(request)

		assert.Equal(t, expected, isNotValid)

	}
}
