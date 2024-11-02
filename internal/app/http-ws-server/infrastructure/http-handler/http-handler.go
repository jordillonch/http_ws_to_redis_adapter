package http_handler

import (
	"net/http"
)

type HttpMiddleware interface {
	Handle(http.Handler) http.Handler
}

type HttpHandler struct {
	middlewares  []HttpMiddleware
	finalHandler http.Handler
}

func NewHttpHandler(finalHandler http.Handler) *HttpHandler {
	return &HttpHandler{
		finalHandler: finalHandler,
	}
}

func (s *HttpHandler) Use(mw HttpMiddleware) {
	s.middlewares = append(s.middlewares, mw)
}

func (s *HttpHandler) WrapMiddleware() http.Handler {
	return s.applyEachMiddlewareInReverseOrder(s.finalHandler)
}

func (s *HttpHandler) applyEachMiddlewareInReverseOrder(handler http.Handler) http.Handler {
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		handler = s.middlewares[i].Handle(handler)
	}
	return handler
}
