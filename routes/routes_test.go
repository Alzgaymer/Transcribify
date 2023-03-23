package routes

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
	"yt-video-transcriptor/models"
)

// Tests only unsuccessful cases
func Test_responseToYTVideo(t *testing.T) {

	testcases := []struct {
		name          string
		inputResponse io.Reader
		data          models.YTVideo
		expectedErr   error
	}{
		{
			name:          "String input (len: 0)",
			inputResponse: strings.NewReader(""),
			data:          models.YTVideo{},
			expectedErr:   io.EOF,
		},
		{
			name:          "Nil reader",
			inputResponse: nil,
			data:          models.YTVideo{},
			expectedErr:   errors.New("io.Reader is nil"),
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {

			var (
				input    = testcase.inputResponse
				expected = testcase.data
			)

			actual, err := responseToYTVideo(input)
			assert.Equal(t, testcase.expectedErr, err)

			assert.Equal(t, expected, actual,
				fmt.Sprintf("Failed %s. \nExpected: %v.\nGot: %v", testcase.name, expected, actual))

		})
	}
}

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
