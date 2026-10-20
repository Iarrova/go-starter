package user

import (
	"context"

	"github.com/google/uuid"
)

func (s *Service) Get(ctx context.Context, id uuid.UUID) (User, error) {
	return s.repository.Get(ctx, id)
}
