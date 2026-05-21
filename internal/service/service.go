package service

import (
	"context"

	"github.com/devize-ed/todo-app/internal/models"
)

type Service struct {
	userRepository Repository
}

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
	UpdateUser(ctx context.Context, id string, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
}

func NewService(userRepository Repository) *Service {
	return &Service{userRepository: userRepository}
}
