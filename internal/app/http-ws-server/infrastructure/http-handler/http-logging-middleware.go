package http_handler

import (
	"go.uber.org/zap"
	"net/http"
)

type LoggingMiddleware struct {
	logger *zap.Logger
}

func NewLoggingMiddleware(logger *zap.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

func (l *LoggingMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.logger.Info("Logging request", zap.String("path", r.URL.Path))
		next.ServeHTTP(w, r)
	})
}
