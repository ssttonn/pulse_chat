package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

type MessageHandler func(message []byte) error

type AsyncConsumer struct {
	handler MessageHandler
}

func NewAsyncConsumer(handler MessageHandler) *AsyncConsumer {
	return &AsyncConsumer{
		handler: handler,
	}
}

func (c *AsyncConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *AsyncConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *AsyncConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		err := c.handler(msg.Value)
		if err != nil {
			log.Printf("Error processing message: %v", err)
			continue
		}

		session.MarkMessage(msg, "")
	}

	return nil
}

func NewConsumerGroup(brokers []string, groupID string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()

	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return sarama.NewConsumerGroup(brokers, groupID, config)
}
