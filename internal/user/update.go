package user

import (
	"context"
	"strings"

	"github.com/Iarrova/go-starter/internal/apperr"
	"github.com/Iarrova/go-starter/internal/validation"
	"github.com/google/uuid"
)

// UpdateInput holds a partial update. A nil field is left unchanged.
type UpdateInput struct {
	Email *string
	Name  *string
}

func (in UpdateInput) normalize() UpdateInput {
	if in.Email != nil {
		email := normalizeEmail(*in.Email)
		in.Email = &email
	}
	if in.Name != nil {
		name := strings.TrimSpace(*in.Name)
		in.Name = &name
	}
	return in
}

func (in UpdateInput) validate() error {
	if in.Email == nil && in.Name == nil {
		return apperr.Invalid("at least one field must be provided")
	}

	var v validation.Validator
	if in.Email != nil {
		checkEmail(&v, *in.Email)
	}
	if in.Name != nil {
		checkName(&v, *in.Name)
	}
	return v.Err()
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, in UpdateInput) (User, error) {
	in = in.normalize()
	if err := in.validate(); err != nil {
		return User{}, err
	}

	return s.repository.Update(ctx, id, in)
}
