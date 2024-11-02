package commands_test

import (
	"github.com/jordillonch/http_ws_to_redis_adapter/internal/app/http-ws-server/domain/topic-message-processor/commands"
	"testing"

	"github.com/jordillonch/http_ws_to_redis_adapter/internal/app/http-ws-server/domain/topic-message-processor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCommandsProcessorProcessSuccesfully(t *testing.T) {
	logger := zap.NewNop()
	processor := commands.NewProcessor(logger)

	testMessage := topic_message_processor.TopicMessage{
		Topic:   "/publisher/topics/test_topic",
		Message: []byte("test message"),
	}
	transformedMessage, err := processor.Process(testMessage)

	require.NoError(t, err, "Process should not return an error")
	assert.Equal(t, "test_topic", transformedMessage.Topic, "Topic should be unchanged")
	assert.Equal(t, testMessage.Message, transformedMessage.Message, "Message content should be unchanged")
}
