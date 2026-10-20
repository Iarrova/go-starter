package post

import (
	"context"
	"strings"

	"github.com/Iarrova/go-starter/internal/apperr"
	"github.com/Iarrova/go-starter/internal/validation"
	"github.com/google/uuid"
)

type UpdateInput struct {
	Title *string
	Body  *string
}

func (in UpdateInput) normalize() UpdateInput {
	if in.Title != nil {
		title := strings.TrimSpace(*in.Title)
		in.Title = &title
	}
	return in
}

func (in UpdateInput) validate() error {
	if in.Title == nil && in.Body == nil {
		return apperr.Invalid("at least one field must be provided")
	}

	var v validation.Validator
	if in.Title != nil {
		checkTitle(&v, *in.Title)
	}
	if in.Body != nil {
		checkBody(&v, *in.Body)
	}
	return v.Err()
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, in UpdateInput) (Post, error) {
	in = in.normalize()
	if err := in.validate(); err != nil {
		return Post{}, err
	}

	return s.repository.Update(ctx, id, in)
}
