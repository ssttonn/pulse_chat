package main

import (
	"context"
	"encoding/json"
	"log"
	"pulse/src/message-router/internal/db"
	"pulse/src/message-router/internal/worker"
	"pulse/src/pkg/config"
	"pulse/src/pkg/dynamo"
	"pulse/src/pkg/kafka"
	"pulse/src/pkg/models"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Cannot initialize config: %v", err)
	}

	ctx := context.Background()

	dynamoClient, err := dynamo.NewClient(ctx, "http://localhost:8000")
	if err != nil {
		log.Fatalf("Cannot initialize DynamoDB Client: %v", err)
	}

	repo := db.NewRepository(dynamoClient.DB)
	batcher := worker.NewBatcher(1000, 25, repo)
	batcher.Start(ctx)

	callback := func(message []byte) error {
		var payload models.ChatPayload
		if err := json.Unmarshal(message, &payload); err != nil {
			return err
		}

		batcher.Add(payload)

		return nil
	}

	handler := kafka.NewAsyncConsumer(callback)
	cg, err := kafka.NewConsumerGroup([]string{cfg.KafkaBrokers}, "message-router-group")

	for {
		if err := cg.Consume(ctx, []string{"chat.inbound"}, handler); err != nil {
			log.Fatalf("Error from consumer: %v", err)
		}
	}
}
