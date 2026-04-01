package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/georgebnunes/todolist-with-go/internal/handler"
	"github.com/georgebnunes/todolist-with-go/internal/repository"
)

func main() {

	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if tableName == "" {
		log.Fatal("DYNAMODB_TABLE_NAME environment variable is not set")
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	dynamodbClient := dynamodb.NewFromConfig(cfg)
	repo := repository.New(dynamodbClient, tableName)

	h := handler.New(repo)

	lambda.Start(h.Route)
}
