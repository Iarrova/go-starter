package userapi

import (
	"context"

	"github.com/Iarrova/go-starter/internal/user"
)

type createInput struct {
	Body createRequest
}

type createRequest struct {
	Email    string `json:"email" example:"ann@example.com"`
	Password string `json:"password" example:"correct-horse-battery"`
	Name     string `json:"name" example:"Ann"`
}

func (h *Handler) Create(ctx context.Context, in *createInput) (*output, error) {
	u, err := h.service.Create(ctx, user.CreateInput{
		Email:    in.Body.Email,
		Password: in.Body.Password,
		Name:     in.Body.Name,
	})
	if err != nil {
		return nil, err
	}

	return &output{Body: toResponse(u)}, nil
}
