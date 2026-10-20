package post

import (
	"context"

	"github.com/google/uuid"
)

func (s *Service) Get(ctx context.Context, id uuid.UUID) (Post, error) {
	return s.repository.Get(ctx, id)
}
