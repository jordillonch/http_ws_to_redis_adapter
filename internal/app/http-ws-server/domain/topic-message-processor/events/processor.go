package events

import (
	"github.com/jordillonch/http_ws_to_redis_adapter/internal/app/http-ws-server/domain/topic-message-processor"
	"go.uber.org/zap"
)

type Processor struct {
	logger *zap.Logger
}

func NewProcessor(logger *zap.Logger) *Processor {
	return &Processor{
		logger: logger,
	}
}

func (f *Processor) Process(message topic_message_processor.TopicMessage) (topic_message_processor.TopicMessage, error) {
	f.logger.Info("Processing message", zap.String("message", string(message.Message)))
	transformedTopicMessage := topic_message_processor.TopicMessage{
		Topic:   message.Topic,
		Message: message.Message,
	}
	return transformedTopicMessage, nil
}
