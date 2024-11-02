package di

import (
	"context"
	"github.com/jordillonch/http_ws_to_redis_adapter/configs"
	"github.com/jordillonch/http_ws_to_redis_adapter/internal/app/http-ws-server/infrastructure/http-handler"
	"github.com/jordillonch/http_ws_to_redis_adapter/internal/app/http-ws-server/infrastructure/transporter"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type HttpWsServerModule struct {
	Logger     *zap.Logger
	HttpPort   int
	HttpServer *http.Server
}

func NewHttpWsServerModule(services *CommonServices, cnf config.Config, processors *MessageProcessors) *HttpWsServerModule {
	mux := http.NewServeMux()
	m := &HttpWsServerModule{
		Logger:   services.Logger,
		HttpPort: cnf.HttpPort,
		HttpServer: &http.Server{
			Addr:    ":" + strconv.Itoa(cnf.HttpPort),
			Handler: mux,
		},
	}

	commandsProcessorAndTransporter := transporter.NewHttpToRedisChannelTransporter(services.Logger, services.RedisClient, processors.Commands)
	commandsHandler := http_handler.NewHttpHandler(commandsProcessorAndTransporter)
	commandsHandler.Use(http_handler.NewLoggingMiddleware(services.Logger))
	commandsHandler.Use(http_handler.NewAuthenticationMiddleware(services.Logger))
	mux.Handle("/", commandsHandler.WrapMiddleware())

	eventsProcessorAndTransporter := transporter.NewRedisChannelToWebsocketTransporter(services.Logger, services.RedisClient, processors.Events, "events")
	eventsHandler := http_handler.NewHttpHandler(eventsProcessorAndTransporter)
	eventsHandler.Use(http_handler.NewLoggingMiddleware(services.Logger))
	mux.Handle("/events", eventsHandler.WrapMiddleware())

	return m
}

func (s HttpWsServerModule) Start() {
	strHttpPort := strconv.Itoa(s.HttpPort)
	s.Logger.Info("Server starting on :" + strHttpPort + "...")
	err := s.HttpServer.ListenAndServe()
	if err != nil {
		s.Logger.Error("ListenAndServe failed:", zap.Error(err))
	}
}

func (s HttpWsServerModule) Stop() {
	s.Logger.Info("Server stopping...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.HttpServer.Shutdown(ctx); err != nil {
		s.Logger.Error("Shutdown failed:", zap.Error(err))
	} else {
		s.Logger.Info("Server stopped")
	}
}
