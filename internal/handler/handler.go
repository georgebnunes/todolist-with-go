package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/georgebnunes/todolist-with-go/internal/model"
	"github.com/georgebnunes/todolist-with-go/internal/repository"
)

type Handler struct {
	repo *repository.TodoRepository
}

func New(repo *repository.TodoRepository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Route(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch {
	case req.HTTPMethod == http.MethodPost && req.Resource == "/todos":
		return h.CreateTodo(ctx, req)
	default:
		return response(http.StatusNotFound, map[string]string{"error": "route not found"})
	}
}

func (h *Handler) CreateTodo(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var body = model.CreateTodoRequest{}

	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return response(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if body.Title == "" || body.UserID == "" {
		return response(http.StatusBadRequest, map[string]string{"error": "title and user ID are required"})
	}

	todo, err := h.repo.Create(ctx, body.Title)

	if err != nil {
		return response(http.StatusInternalServerError, map[string]string{"error": "failed to create todo"})
	}

	// Implementation for creating a todo
	return response(http.StatusCreated, todo)
}

func response(statusCode int, body any) (events.APIGatewayProxyResponse, error) {
	var bodyStr string

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500}, err
		}
		bodyStr = string(b)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       bodyStr,
	}, nil
}
