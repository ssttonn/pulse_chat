#!/bin/bash

echo "Creating local DynamoDB table: chat_messages..."

aws dynamodb create-table \
    --endpoint-url http://localhost:8000 \
    --region us-east-1 \
    --table-name chat_messages \
    --attribute-definitions \
        AttributeName=channel_id,AttributeType=S \
        AttributeName=created_at,AttributeType=S \
    --key-schema \
        AttributeName=channel_id,KeyType=HASH \
        AttributeName=created_at,KeyType=RANGE \
    --billing-mode PAY_PER_REQUEST

echo -e "\nTable created successfully!"