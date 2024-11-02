package transporter_test

import (
	"bytes"
	"context"
	"github.com/jordillonch/http_ws_to_redis_adapter/internal/app/http-ws-server/infrastructure/transporter"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jordillonch/http_ws_to_redis_adapter/internal/app/http-ws-server/domain/topic-message-processor"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TODO:
//   - Refactor
//   - Improve to cover the case where the processor returns an error
func TestIntegrationHttpToRedisChannelTransporter(t *testing.T) {
	// Arrange ---------------------------------------------------------------
	// Set up a real Redis client (assuming Redis is running on localhost:6379)
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer redisClient.Close()
	// Clear the Redis topic before the test to avoid interference from previous tests
	topic := "/publisher/topics/SOME_TOPIC"
	redisClient.Del(context.Background(), topic)
	// Initialize dependencies
	logger := zap.NewNop() // Use a no-op logger for testing
	processor := &FakeHttpToRedisProcessor{}
	transporter := transporter.NewHttpToRedisChannelTransporter(logger, redisClient, processor)
	// Set up a Redis subscriber to capture published messages
	subscriber := redisClient.Subscribe(context.Background(), topic)
	defer subscriber.Close()
	_, err := subscriber.Receive(context.Background())
	assert.NoError(t, err, "Failed to subscribe to topic")

	// Act -------------------------------------------------------------------
	// HTTP request
	body := []byte("test message")
	req := httptest.NewRequest("POST", topic, bytes.NewReader(body))
	w := httptest.NewRecorder()
	transporter.ServeHTTP(w, req)

	// Assert ----------------------------------------------------------------
	// Validate HTTP response
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	// Check that the message was published to Redis
	// Wait for the message with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg, err := subscriber.ReceiveMessage(ctx)
	assert.NoError(t, err, "Failed to receive message from Redis")
	if err == nil {
		// Validate that the message payload is what we expect
		expectedMessage := "test message processed"
		assert.Equal(t, expectedMessage, msg.Payload)
	}
}

type FakeHttpToRedisProcessor struct{}

func (p *FakeHttpToRedisProcessor) Process(message topic_message_processor.TopicMessage) (topic_message_processor.TopicMessage, error) {
	return topic_message_processor.TopicMessage{
		Topic:   message.Topic,
		Message: append(message.Message, []byte(" processed")...),
	}, nil
}
