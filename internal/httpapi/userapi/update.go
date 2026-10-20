package userapi

import (
	"context"

	"github.com/google/uuid"

	"github.com/Iarrova/go-starter/internal/user"
)

type updateInput struct {
	ID   uuid.UUID `path:"id" format:"uuid"`
	Body updateRequest
}

type updateRequest struct {
	Email *string `json:"email,omitempty" example:"ann@example.com"`
	Name  *string `json:"name,omitempty" example:"Ann"`
}

func (h *Handler) Update(ctx context.Context, in *updateInput) (*output, error) {
	u, err := h.service.Update(ctx, in.ID, user.UpdateInput{
		Email: in.Body.Email,
		Name:  in.Body.Name,
	})
	if err != nil {
		return nil, err
	}

	return &output{Body: toResponse(u)}, nil
}
