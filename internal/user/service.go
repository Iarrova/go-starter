package user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, in CreateParams) (User, error)
	Get(ctx context.Context, id uuid.UUID) (User, error)
	Update(ctx context.Context, id uuid.UUID, in UpdateInput) (User, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}
