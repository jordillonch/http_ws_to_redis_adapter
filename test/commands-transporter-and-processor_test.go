package test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO:
//   - Refactor
//   - Improve to cover the case where the processor returns an error
func TestHttpToRedisChannelTransporter_PublishToRedis(t *testing.T) {
	// Arrange -----------------------------------------------------------------
	setUp()
	go common.HttpWsServerServices.Start()
	// Set up a real Redis client (assumes Redis is running on localhost:6379)
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer redisClient.Close()
	// Set up Redis subscriber for the channel we want to test
	channel := "SOME_TOPIC"
	subscriber := redisClient.Subscribe(context.Background(), channel)
	defer subscriber.Close()
	// Wait for the subscription to be established
	_, err := subscriber.Receive(context.Background())
	require.NoError(t, err, "Failed to subscribe to Redis channel")
	// Define the message body to send
	body := []byte("test message")

	// Act ---------------------------------------------------------------------
	// Perform a POST request to the transporter endpoint
	url := fmt.Sprintf("http://localhost:8080/publisher/topics/%s", channel)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	require.NoError(t, err, "Failed to create HTTP request")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "Failed to perform HTTP request")
	defer resp.Body.Close()

	// Assert ------------------------------------------------------------------
	// Assert that the response status code is 201 Created
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	// Wait for a message from Redis with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// Receive the message from Redis
	msg, err := subscriber.ReceiveMessage(ctx)
	require.NoError(t, err, "Failed to receive message from Redis")
	// Assert that the message payload matches the request body
	assert.Equal(t, string(body), msg.Payload)
}
