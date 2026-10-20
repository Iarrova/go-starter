package postapi

import (
	"context"

	"github.com/google/uuid"

	"github.com/Iarrova/go-starter/internal/httpapi"
	"github.com/Iarrova/go-starter/internal/post"
)

type listInput struct {
	UserID uuid.UUID `query:"user_id" format:"uuid" doc:"Only posts by this user"`
	httpapi.PageQuery
}

type pageOutput struct {
	Body httpapi.PageResponse[response]
}

func (h *Handler) List(ctx context.Context, in *listInput) (*pageOutput, error) {
	result, err := h.service.List(ctx, post.ListFilter{UserID: in.UserID}, in.Request())
	if err != nil {
		return nil, err
	}

	return &pageOutput{Body: httpapi.ToPageResponse(result, toResponse)}, nil
}
