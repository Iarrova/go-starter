package postapi

import (
	"context"

	"github.com/google/uuid"

	"github.com/Iarrova/go-starter/internal/post"
)

type createInput struct {
	Body createRequest
}

type createRequest struct {
	// TODO: Take the author from the authenticated user once auth exists.
	UserID uuid.UUID `json:"user_id" format:"uuid"`
	Title  string    `json:"title" example:"Hello, world"`
	Body   string    `json:"body" example:"My first post."`
}

func (h *Handler) Create(ctx context.Context, in *createInput) (*output, error) {
	p, err := h.service.Create(ctx, in.Body.UserID, post.CreateInput{
		Title: in.Body.Title,
		Body:  in.Body.Body,
	})
	if err != nil {
		return nil, err
	}

	return &output{Body: toResponse(p)}, nil
}
