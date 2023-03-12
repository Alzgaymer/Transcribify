package middlewares

import (
	"bytes"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strings"
	"yt-video-transcriptor/routes"
)

func Logging(logger *zap.Logger) func(next http.Handler) http.Handler {
	//middleware
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				videoRequest = routes.VideoRequest{}
			)
			// Copying url
			query, err := copyUrl(r.URL.RawQuery)
			if err != nil {
				w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
				return
			}
			videoRequest.VideoID = query.Get("v")
			videoRequest.Language = query.Get("lang")

			logger.Info("Input parameter",
				zap.String("Video ID", videoRequest.VideoID),
				zap.String("Language", videoRequest.Language),
				zap.String("URL", r.URL.Path),
			)

			next.ServeHTTP(w, r)
		})

	}
}

func copyUrl(uri string) (url.Values, error) {
	var (
		buf      []byte
		writeBuf = bytes.NewBuffer(buf)
	)
	_, err := io.Copy(writeBuf, strings.NewReader(uri))
	if err != nil {
		return nil, err
	}

	return url.ParseQuery(writeBuf.String())
}
