package post

import (
	"context"

	"github.com/Iarrova/go-starter/internal/pagination"
	"github.com/google/uuid"
)

type ListFilter struct {
	// Empty means no filter
	UserID uuid.UUID
}

func (s *Service) List(ctx context.Context, f ListFilter, p pagination.Request) (pagination.Page[Post], error) {
	return s.repository.List(ctx, f, p.Normalize())
}
