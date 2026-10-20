package post

import (
	"context"

	"github.com/Iarrova/go-starter/internal/pagination"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, in CreateParams) (Post, error)
	Get(ctx context.Context, id uuid.UUID) (Post, error)
	Update(ctx context.Context, id uuid.UUID, in UpdateInput) (Post, error)
	List(ctx context.Context, f ListFilter, p pagination.Request) (pagination.Page[Post], error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}
