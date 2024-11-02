package test

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO:
//   - Refactor
//   - Improve to cover the case where the processor returns an error
func TestRedisChannelToWebsocketTransporter_MessageReceivedInWebSocket(t *testing.T) {
	// Arrange --------------------------------------------------------------
	setUp()
	go common.HttpWsServerServices.Start()
	// Set up a real Redis client (assumes Redis is running on localhost:6379)
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer redisClient.Close()
	// Parse the WebSocket URL
	wsURL := "ws://localhost:8080/events"
	// Connect a WebSocket client to the transporter
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err, "Failed to connect to WebSocket server")
	defer ws.Close()
	// Define the message to publish to Redis
	publishMessage := "hello from redis"

	// Act ------------------------------------------------------------------
	// Publish the message to the Redis channel in a separate goroutine
	go func() {
		time.Sleep(1 * time.Second) // Ensure the WebSocket connection is ready
		channel := "events"
		err := redisClient.Publish(context.Background(), channel, publishMessage).Err()
		require.NoError(t, err, "Failed to publish message to Redis")
	}()

	// Assert ---------------------------------------------------------------
	// Set a read deadline for the WebSocket to avoid blocking indefinitely
	ws.SetReadDeadline(time.Now().Add(5 * time.Second))
	// Read the message from WebSocket
	_, msg, err := ws.ReadMessage()
	require.NoError(t, err, "Failed to read message from WebSocket")
	assert.Equal(t, publishMessage, string(msg), "Message received from WebSocket should match the Redis published message")
}
