package db

import (
	"context"
	"pulse/src/pkg/models"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Repository struct {
	db *dynamodb.Client
}

func (r *Repository) BatchInsertMessages(ctx context.Context, messages []models.ChatPayload) error {
	if len(messages) == 0 {
		return nil
	}

	writeRequests := make([]types.WriteRequest, 0, len(messages))

	for _, message := range messages {
		itemMap, err := attributevalue.MarshalMap(message)

		if err != nil {
			continue
		}

		writeRequests = append(writeRequests, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: itemMap,
			},
		})
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			"chat_messages": writeRequests,
		},
	}

	_, err := r.db.BatchWriteItem(ctx, input)

	return err
}
