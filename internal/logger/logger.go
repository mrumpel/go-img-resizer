package logger

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Logger struct {
	logrus.Logger
}

func NewLogger(lvl string) *Logger {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		level = logrus.TraceLevel
	}

	l := *logrus.New()
	l.Level = level
	//nolint
	return &Logger{l}
}

func (l *Logger) GetRequestLoggingHandler() func(handler http.Handler) http.Handler {
	return func(wrapped http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			wrapped.ServeHTTP(w, r)
			l.WithFields(map[string]interface{}{
				"method":         r.Method,
				"path":           r.URL.Path,
				"duration":       time.Since(now),
				"remote_address": r.RemoteAddr,
				"user_agent":     r.UserAgent(),
			}).Info()
		})
	}
}
