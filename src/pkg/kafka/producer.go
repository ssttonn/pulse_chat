package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

type AsyncProducer struct {
	producer sarama.AsyncProducer
}

func NewProducer(brokers []string) (*AsyncProducer, error) {
	config := sarama.NewConfig()

	config.Producer.Return.Successes = false
	config.Producer.Return.Errors = true

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	go func() {
		for err := range producer.Errors() {
			log.Printf("Kafka async producer error: %v", err)
		}
	}()

	return &AsyncProducer{
		producer: producer,
	}, nil
}

func (ap *AsyncProducer) Publish(topic string, key string, value []byte) {
	message := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	ap.producer.Input() <- message
}

func (ap *AsyncProducer) Close() error {
	return ap.producer.Close()
}
