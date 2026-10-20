package postapi

import (
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/Iarrova/go-starter/internal/post"
)

type Handler struct {
	service *post.Service
}

func New(service *post.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "create-post",
		Method:        http.MethodPost,
		Path:          "/posts",
		Summary:       "Create a post",
		Tags:          []string{"posts"},
		DefaultStatus: http.StatusCreated,
		Errors:        []int{http.StatusBadRequest},
	}, h.Create)

	huma.Register(api, huma.Operation{
		OperationID: "get-post",
		Method:      http.MethodGet,
		Path:        "/posts/{id}",
		Summary:     "Get a post",
		Tags:        []string{"posts"},
		Errors:      []int{http.StatusNotFound},
	}, h.Get)

	huma.Register(api, huma.Operation{
		OperationID: "update-post",
		Method:      http.MethodPatch,
		Path:        "/posts/{id}",
		Summary:     "Update a post",
		Description: "Updates only the fields present in the body.",
		Tags:        []string{"posts"},
		Errors:      []int{http.StatusBadRequest, http.StatusNotFound},
	}, h.Update)

	huma.Register(api, huma.Operation{
		OperationID: "list-posts",
		Method:      http.MethodGet,
		Path:        "/posts",
		Summary:     "List posts",
		Description: "Returns posts newest first. Pass next_cursor as cursor to get the next page.",
		Tags:        []string{"posts"},
		Errors:      []int{http.StatusUnprocessableEntity},
	}, h.List)
}

type output struct {
	Body response
}

type response struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toResponse(p post.Post) response {
	return response{
		ID:        p.ID,
		UserID:    p.UserID,
		Title:     p.Title,
		Body:      p.Body,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
