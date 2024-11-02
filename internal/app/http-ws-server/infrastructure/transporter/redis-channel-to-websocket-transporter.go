package transporter

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/jordillonch/http_ws_to_redis_adapter/internal/app/http-ws-server/domain/topic-message-processor"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

type RedisChannelToWebsocketTransporter struct {
	ctx         context.Context
	logger      *zap.Logger
	redisClient *redis.Client
	processor   topic_message_processor.TopicMessageProcessor
	clients     sync.Map
	broadcast   chan string
	upgrader    websocket.Upgrader
}

func NewRedisChannelToWebsocketTransporter(logger *zap.Logger, redisClient *redis.Client, processor topic_message_processor.TopicMessageProcessor, channel string) *RedisChannelToWebsocketTransporter {
	t := &RedisChannelToWebsocketTransporter{
		ctx:         context.Background(),
		logger:      logger,
		redisClient: redisClient,
		processor:   processor,
		clients:     sync.Map{},
		broadcast:   make(chan string),
		upgrader:    websocket.Upgrader{},
	}
	go t.broadcastToWebSocketClients()
	go t.listenRedisChannel(channel)

	return t
}

func (t *RedisChannelToWebsocketTransporter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := t.upgrader.Upgrade(w, r, nil)
	if err != nil {
		t.logger.Error("Failed to upgrade WebSocket", zap.Error(err))
		return
	}
	defer ws.Close()

	t.clients.Store(ws, true)

	// Wait for broadcast messages and send them to the WebSocket client
	for {
		msg, ok := <-t.broadcast
		if !ok {
			return
		}
		if err := ws.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			t.logger.Error("Error sending message to WebSocket client", zap.Error(err))
			ws.Close()
			t.clients.Delete(ws)
			return
		}
	}
}

func (t *RedisChannelToWebsocketTransporter) listenRedisChannel(channel string) {
	pubsub := t.redisClient.Subscribe(t.ctx, channel)
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(t.ctx)
		if err != nil {
			t.logger.Error("Error receiving message from Redis", zap.Error(err))
			continue
		}
		t.broadcast <- msg.Payload
	}
}

func (t *RedisChannelToWebsocketTransporter) broadcastToWebSocketClients() {
	for msg := range t.broadcast {
		t.clients.Range(func(client, _ interface{}) bool {
			if err := client.(*websocket.Conn).WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				t.logger.Error("Error broadcasting to WebSocket client:", zap.Error(err))
				client.(*websocket.Conn).Close()
				t.clients.Delete(client)
			}
			return true
		})
	}
}
