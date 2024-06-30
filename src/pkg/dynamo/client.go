package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Client struct {
	DB *dynamodb.Client
}

func NewClient(ctx context.Context, endpointURL string) (*Client, error) {
	var cfg aws.Config
	var err error
	if endpointURL != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           endpointURL,
				SigningRegion: "us-east-1",
			}, nil
		})

		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"), config.WithEndpointResolverWithOptions(customResolver), config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "Hard-coded fake credentials for local DynamoDB",
			},
		}))
	} else {
		cfg, err = config.LoadDefaultConfig(ctx)
	}

	if err != nil {
		return nil, err
	}

	return &Client{
		DB: dynamodb.NewFromConfig(cfg),
	}, nil
}
