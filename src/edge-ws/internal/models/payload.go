package models

import "encoding/json"

type MessageType string

const (
	MessageTypeAuth MessageType = "auth"
	MessageTypeChat MessageType = "chat"
	MessageTypePing MessageType = "ping"
)

type InboundMessage struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type AuthPayload struct {
	Token string `json:"token"`
}

type ChatPayload struct {
	ChannelID string `json:"channel_id"`
	Text      string `json:"text"`
}
