package handler

import (
	"context"
	"encoding/json"
	"log"
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

func (h *Handler) Route(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	method := req.RequestContext.HTTP.Method // ← fix principal
	path := req.RequestContext.HTTP.Path     // ← fix principal
	id := req.PathParameters["id"]

	log.Printf("method=%s path=%s id=%s", method, path, id)

	switch method {
	case http.MethodPost:
		return h.CreateTodo(ctx, req)

	default:
		return response(http.StatusNotFound, map[string]string{"error": "route not found"})
	}
}

func (h *Handler) CreateTodo(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

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

func response(statusCode int, body any) (events.APIGatewayV2HTTPResponse, error) {
	var bodyStr string

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return events.APIGatewayV2HTTPResponse{StatusCode: 500}, err
		}
		bodyStr = string(b)
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: statusCode,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       bodyStr,
	}, nil
}
