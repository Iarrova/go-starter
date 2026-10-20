package post

import (
	"context"
	"strings"

	"github.com/Iarrova/go-starter/internal/validation"
	"github.com/google/uuid"
)

type CreateInput struct {
	Title string
	Body  string
}

type CreateParams struct {
	UserID uuid.UUID
	Title  string
	Body   string
}

func (in CreateInput) normalize() CreateInput {
	in.Title = strings.TrimSpace(in.Title)
	return in
}

func (in CreateInput) validate() error {
	var v validation.Validator
	checkTitle(&v, in.Title)
	checkBody(&v, in.Body)
	return v.Err()
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID, in CreateInput) (Post, error) {
	in = in.normalize()
	if err := in.validate(); err != nil {
		return Post{}, err
	}

	return s.repository.Create(ctx, CreateParams{
		UserID: userID,
		Title:  in.Title,
		Body:   in.Body,
	})
}
