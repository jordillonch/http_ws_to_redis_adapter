package commands

import (
	"errors"
	"github.com/jordillonch/http_ws_to_redis_adapter/internal/app/http-ws-server/domain/topic-message-processor"
	"go.uber.org/zap"
	"strings"
)

type Processor struct {
	logger *zap.Logger
}

func NewProcessor(logger *zap.Logger) *Processor {
	return &Processor{
		logger: logger,
	}
}

func (p *Processor) Process(message topic_message_processor.TopicMessage) (topic_message_processor.TopicMessage, error) {
	pathParts := strings.Split(message.Topic, "/")
	topicName := ""
	if len(pathParts) >= 3 {
		topicName = pathParts[3]
		p.logger.Info("Publishing message to topic", zap.String("topic", topicName))
	} else {
		return message, errors.New("Invalid topic")
	}

	// TODO: Implement the logic to process the message
	p.logger.Info("Processing message", zap.String("message", string(message.Message)))
	transformedTopicMessage := topic_message_processor.TopicMessage{
		Topic:   topicName,
		Message: message.Message,
	}
	return transformedTopicMessage, nil
}
