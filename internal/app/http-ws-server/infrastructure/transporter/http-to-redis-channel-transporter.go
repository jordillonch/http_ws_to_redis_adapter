package transporter

import (
	"github.com/go-redis/redis/v8"
	"github.com/jordillonch/http_ws_to_redis_adapter/internal/app/http-ws-server/domain/topic-message-processor"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type HttpToRedisChannelTransporter struct {
	logger      *zap.Logger
	redisClient *redis.Client
	processor   topic_message_processor.TopicMessageProcessor
}

func NewHttpToRedisChannelTransporter(logger *zap.Logger, redisClient *redis.Client, processor topic_message_processor.TopicMessageProcessor) *HttpToRedisChannelTransporter {
	return &HttpToRedisChannelTransporter{logger, redisClient, processor}
}

func (t *HttpToRedisChannelTransporter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.logger.Error("Failed to read message from request body", zap.Error(err))
		http.Error(w, "Failed to read message from request body", http.StatusInternalServerError)
	}
	message := topic_message_processor.TopicMessage{
		Topic:   r.URL.Path,
		Message: body,
	}
	processedMessage, err := t.processor.Process(message)

	// Publish message to Redis
	err = t.redisClient.Publish(r.Context(), processedMessage.Topic, processedMessage.Message).Err()
	if err != nil {
		t.logger.Error("Failed to publish message to topic", zap.Error(err))
		http.Error(w, "Failed to publish message to topic", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}
