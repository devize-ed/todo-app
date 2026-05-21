package handlers

import (
	"context"

	"github.com/devize-ed/todo-app/internal/models"
)

type Service interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
	UpdateUser(ctx context.Context, id string, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}
