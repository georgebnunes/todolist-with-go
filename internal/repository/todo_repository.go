package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/georgebnunes/todolist-with-go/internal/model"
	"github.com/google/uuid"
)

type TodoRepository struct {
	client    *dynamodb.Client
	tableName string
}

func New(client *dynamodb.Client, tableName string) *TodoRepository {
	return &TodoRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *TodoRepository) Create(ctx context.Context, title string) (*model.TodoItem, error) {

	now := time.Now().UTC()

	todo := &model.TodoItem{
		ID:        uuid.NewString(),
		Title:     title,
		Done:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	item, err := attributevalue.MarshalMap(todo)
	if err != nil {
		return nil, fmt.Errorf("marshling todo: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})

	if err != nil {
		return nil, fmt.Errorf("putting item in dynamodb: %w", err)
	}

	return todo, nil
}
