package middlewares

import (
	"fmt"
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
				lang := r.URL.Query().Get(models.LanguageTag)
				if lang == "" {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
			}),
			Endpoint:     "/00000000000",                            //11
			Query:        fmt.Sprintf("?%s=ua", models.LanguageTag), //?lang=ua
			ExpectedLang: "ua",
		},
		{
			Name:         "Using through maps",
			ExpectedCode: http.StatusOK,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				lang := r.URL.Query()[models.LanguageTag][0]
				if lang == "" {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
			}),
			Endpoint:     "/00000000000",                            //11
			Query:        fmt.Sprintf("?%s=ua", models.LanguageTag), //?lang=ua
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
