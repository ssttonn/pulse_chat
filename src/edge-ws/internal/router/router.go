package router

import (
	"encoding/json"
	"log"
	"net"
	"sync"

	"pulse/src/edge-ws/internal/connection"
	"pulse/src/edge-ws/internal/models"
	"pulse/src/pkg/kafka"
)

type Router struct {
	producer    *kafka.AsyncProducer
	connections sync.Map
	secret      string
}

func NewRouter(producer *kafka.AsyncProducer, secret string) *Router {
	return &Router{
		producer: producer,
		secret:   secret,
	}
}

func (r *Router) HandleConnection(conn net.Conn) {
	userID, ok := r.connections.Load(conn)
	if !ok {
		id, err := connection.HandleAuthFrame(conn, r.secret)
		if err != nil {
			log.Printf("Auth failed for %s, closing connection: %v", conn.RemoteAddr(), err)
			conn.Close()
			return
		}

		r.connections.Store(conn, id)
		log.Printf("User %s authenticated successfully on %s!", id, conn.RemoteAddr())
		return
	}

	r.handleChatFrame(conn, userID.(string))
}

func (r *Router) handleChatFrame(conn net.Conn, userID string) {
	payload, release, err := connection.ReadFrame(conn)
	if err != nil {
		log.Printf("Failed to read frame for user %s: %v", userID, err)
		r.connections.Delete(conn)
		conn.Close()
		return
	}
	defer release() // The buffer goes back to the pool exactly when this function exits!

	// Decode wrapper message
	var inbound models.InboundMessage
	if err := json.Unmarshal(payload, &inbound); err != nil {
		log.Printf("Invalid JSON from user %s: %v", userID, err)
		return
	}

	// Route Chat messages to Kafka
	if inbound.Type == models.MessageTypeChat {
		var chat models.ChatPayload
		if err := json.Unmarshal(inbound.Payload, &chat); err == nil {
			r.producer.Publish("chat.inbound", chat.ChannelID, inbound.Payload)
		} else {
			log.Printf("Invalid chat payload from user %s: %v", userID, err)
		}
	}
}
