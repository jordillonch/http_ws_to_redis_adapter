package http_handler

import (
	"go.uber.org/zap"
	"net/http"
)

type AuthenticationMiddleware struct {
	logger *zap.Logger
}

func NewAuthenticationMiddleware(logger *zap.Logger) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		logger: logger,
	}
}

func (a *AuthenticationMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement a proper authentication mechanism
		a.logger.Info("AuthenticationMiddleware: Request received (doing nothing... TODO)")
		next.ServeHTTP(w, r)
	})
}
