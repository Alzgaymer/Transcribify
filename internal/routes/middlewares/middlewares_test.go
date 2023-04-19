package middlewares

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"transcribify/internal/models"
	"transcribify/pkg/logging"
)

func TestLogging(t *testing.T) {

	type test struct {
		Name         string
		ExpectedCode int
		Handler      http.Handler
		Endpoint     string
		Query        string
		ExpectedLang string
	}

	logger, err := logging.New(
		logging.WithOutputPaths("stderr"),
	)
	if err != nil {
		t.Errorf("Failed to create logger:%s", err.Error())
	}
	testData := []test{
		{
			Name:         "Using url.Values",
			ExpectedCode: http.StatusOK,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				lang := r.URL.Query().Get("lang")
				if lang == "" {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
			}),
			Endpoint:     "/00000000000", //11
			Query:        "?lang=ua",
			ExpectedLang: "ua",
		},
		{
			Name:         "Using through maps",
			ExpectedCode: http.StatusOK,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				lang := r.URL.Query()["lang"][0]
				if lang == "" {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
			}),
			Endpoint:     "/00000000000", //11
			Query:        "?lang=ua",     //?lang=ua
			ExpectedLang: "ua",
		},
	}
	for _, testCase := range testData {
		t.Run(testCase.Name, func(t *testing.T) {
			url := "http://www.example.com" + testCase.Endpoint + testCase.Query
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			handlerWithMiddleware := LogVideoRequest(logger)(testCase.Handler)

			handlerWithMiddleware.ServeHTTP(w, req)

			assert.Equal(t, testCase.ExpectedCode, w.Code)
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
				Language: "ua",
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
				Language: "ua",
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
		t.Run(testcase.name, func(t *testing.T) {
			var (
				request  = testcase.data
				expected = testcase.expectedRes
			)

			isNotValid, _ := ValidateVideoRequest(request)

			assert.Equal(t, expected, isNotValid)
		})

	}
}
