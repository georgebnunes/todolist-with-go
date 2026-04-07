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


func (r *TodoRepository) ListTodos(ctx context.Context) ([]model.TodoItem, error){
	output, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(r.tableName)
	})

	if err != nil {
		return nil, fmt.Errorf("scanning table: %w", err)
	}

	var items []model.TodoItem

	for _, it range output.Items {
		var item model.TodoItem
		err := attributevalue.UnmarshalMap(it, &item)

		if err != nil {
			return nil, fmt.Errorf("unmarshaling item: %w", err)
		}

		items = append(items, item)
	}

	return items, nil
}