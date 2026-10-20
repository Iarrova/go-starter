package userapi

import (
	"context"

	"github.com/google/uuid"
)

type getInput struct {
	ID uuid.UUID `path:"id" format:"uuid"`
}

func (h *Handler) Get(ctx context.Context, in *getInput) (*output, error) {
	u, err := h.service.Get(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	return &output{Body: toResponse(u)}, nil
}
