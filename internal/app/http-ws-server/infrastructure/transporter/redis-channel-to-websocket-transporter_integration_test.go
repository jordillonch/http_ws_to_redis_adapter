package transporter_test

import (
	"context"
	"github.com/jordillonch/http_ws_to_redis_adapter/internal/app/http-ws-server/infrastructure/transporter"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/jordillonch/http_ws_to_redis_adapter/internal/app/http-ws-server/domain/topic-message-processor"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TODO:
//   - Refactor
//   - Improve to cover the case where the processor returns an error
func TestRedisChannelToWebsocketTransporter(t *testing.T) {
	// Arrange ---------------------------------------------------------------
	// Set up Redis client (assumes Redis is running locally on port 6379)
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer redisClient.Close()
	// Set up logger and processor
	logger := zap.NewNop()
	processor := &FakeRedisToWebsocketProcessor{}
	// Redis channel to use
	channel := "test_topic"
	// Create an instance of RedisChannelToWebsocketTransporter
	transporter := transporter.NewRedisChannelToWebsocketTransporter(logger, redisClient, processor, channel)
	// Create an HTTP server and WebSocket endpoint
	server := httptest.NewServer(transporter)
	defer server.Close()
	// Parse the WebSocket URL
	wsURL := "ws" + server.URL[4:] + "/ws"
	// Connect a WebSocket client to the transporter
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err, "Failed to connect to WebSocket server")
	defer ws.Close()

	// Act -------------------------------------------------------------------
	// Publish a message to the Redis channel
	redisMessage := "hello from redis"
	go func() {
		time.Sleep(1 * time.Second) // Ensure the WebSocket connection is ready
		err := redisClient.Publish(context.Background(), channel, redisMessage).Err()
		assert.NoError(t, err, "Failed to publish message to Redis")
	}()

	// Assert ----------------------------------------------------------------
	// Set a read deadline for the WebSocket message to avoid blocking
	ws.SetReadDeadline(time.Now().Add(5 * time.Second))
	// Read the message from WebSocket
	_, msg, err := ws.ReadMessage()
	assert.NoError(t, err, "Failed to read message from WebSocket")
	assert.Equal(t, redisMessage, string(msg), "Message received from WebSocket should match the Redis published message")
}

type FakeRedisToWebsocketProcessor struct{}

func (p *FakeRedisToWebsocketProcessor) Process(message topic_message_processor.TopicMessage) (topic_message_processor.TopicMessage, error) {
	return topic_message_processor.TopicMessage{
		Topic:   message.Topic,
		Message: append(message.Message, []byte(" processed")...),
	}, nil
}
