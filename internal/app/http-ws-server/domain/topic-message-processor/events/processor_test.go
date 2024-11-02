package events_test

import (
	"github.com/jordillonch/http_ws_to_redis_adapter/internal/app/http-ws-server/domain/topic-message-processor/events"
	"testing"

	"github.com/jordillonch/http_ws_to_redis_adapter/internal/app/http-ws-server/domain/topic-message-processor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestEventsProcessorProcessSuccessfully(t *testing.T) {
	logger := zap.NewNop()
	processor := events.NewProcessor(logger)

	testMessage := topic_message_processor.TopicMessage{
		Topic:   "test_topic",
		Message: []byte("test message"),
	}
	transformedMessage, err := processor.Process(testMessage)

	require.NoError(t, err, "Process should not return an error")
	assert.Equal(t, testMessage.Topic, transformedMessage.Topic, "Topic should be unchanged")
	assert.Equal(t, testMessage.Message, transformedMessage.Message, "Message content should be unchanged")
}
