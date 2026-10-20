package postapi

import (
	"context"

	"github.com/google/uuid"

	"github.com/Iarrova/go-starter/internal/post"
)

type updateInput struct {
	ID   uuid.UUID `path:"id" format:"uuid"`
	Body updateRequest
}

type updateRequest struct {
	Title *string `json:"title,omitempty" example:"Hello, world"`
	Body  *string `json:"body,omitempty" example:"My first post."`
}

func (h *Handler) Update(ctx context.Context, in *updateInput) (*output, error) {
	p, err := h.service.Update(ctx, in.ID, post.UpdateInput{
		Title: in.Body.Title,
		Body:  in.Body.Body,
	})
	if err != nil {
		return nil, err
	}

	return &output{Body: toResponse(p)}, nil
}
