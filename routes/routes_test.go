package routes

import (
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
	"yt-video-transcriptor/models"
)

func Test_responseToYTVideo(t *testing.T) {
	testcases := []struct{
		response io.Reader
		data []models.YTVideo
	}
}

func Test_ValidationVideoRequest(t *testing.T) {
	t.Run("isNotValidVideoRequest: Input request params", test_isNotValidVideoRequest)
	t.Run("isValidVideoRequest: Input request params", test_isValidVideoRequest)
}

func test_isNotValidVideoRequest(t *testing.T) {
	tests := []struct {
		data        VideoRequest
		expectedRes bool
	}{
		{
			data: VideoRequest{
				VideoID:  "12345678901",
				Language: "ru",
			},
			expectedRes: false,
		},
		{
			data: VideoRequest{
				VideoID:  "1234567890",
				Language: "ru",
			},
			expectedRes: true,
		},
		{
			data: VideoRequest{
				VideoID:  "12345678901",
				Language: "r1",
			},
			expectedRes: true,
		},
		{
			data: VideoRequest{
				VideoID:  "12345678901",
				Language: "1u",
			},
			expectedRes: true,
		},
		{
			data: VideoRequest{
				VideoID:  "12345678901",
				Language: "1",
			},
			expectedRes: true,
		},
		{
			data: VideoRequest{
				VideoID:  "",
				Language: "1",
			},
			expectedRes: true,
		},
		{
			data: VideoRequest{
				VideoID:  "12345678901",
				Language: "",
			},
			expectedRes: true,
		},
		{
			data: VideoRequest{
				VideoID:  "",
				Language: "",
			},
			expectedRes: true,
		},
	}

	for _, testcase := range tests {
		var (
			request  = testcase.data
			expected = testcase.expectedRes
		)

		isNotValid, _ := isNotValidVideoRequest(request)

		assert.Equal(t, expected, isNotValid)

	}
}
func test_isValidVideoRequest(t *testing.T) {
	tests := []struct {
		data        VideoRequest
		expectedRes bool
	}{
		{
			data: VideoRequest{
				VideoID:  "12345678901",
				Language: "ru",
			},
			expectedRes: true,
		},
		{
			data: VideoRequest{
				VideoID:  "1234567890",
				Language: "ru",
			},
			expectedRes: false,
		},
		{
			data: VideoRequest{
				VideoID:  "12345678901",
				Language: "r1",
			},
			expectedRes: false,
		},
		{
			data: VideoRequest{
				VideoID:  "12345678901",
				Language: "1u",
			},
			expectedRes: false,
		},
		{
			data: VideoRequest{
				VideoID:  "12345678901",
				Language: "1",
			},
			expectedRes: false,
		},
		{
			data: VideoRequest{
				VideoID:  "",
				Language: "1",
			},
			expectedRes: false,
		},
		{
			data: VideoRequest{
				VideoID:  "12345678901",
				Language: "",
			},
			expectedRes: false,
		},
		{
			data: VideoRequest{
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
