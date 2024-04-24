package events

type RoutedMessage struct {
	ID        string `json:"id"`
	SenderID  string `json:"sender_id"`
	ChannelID string `json:"channel_id"`
	Text      string `json:"text"`
	Timestamp int64  `json:"timestamp"`
}
