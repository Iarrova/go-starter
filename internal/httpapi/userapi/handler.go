package userapi

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/Iarrova/go-starter/internal/user"
)

type Handler struct {
	service *user.Service
}

func New(service *user.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "create-user",
		Method:        http.MethodPost,
		Path:          "/users",
		Summary:       "Create a user",
		Tags:          []string{"users"},
		DefaultStatus: http.StatusCreated,
		Errors:        []int{http.StatusBadRequest, http.StatusConflict},
	}, h.Create)

	huma.Register(api, huma.Operation{
		OperationID: "get-user",
		Method:      http.MethodGet,
		Path:        "/users/{id}",
		Summary:     "Get a user",
		Tags:        []string{"users"},
		Errors:      []int{http.StatusNotFound},
	}, h.Get)

	huma.Register(api, huma.Operation{
		OperationID: "update-user",
		Method:      http.MethodPatch,
		Path:        "/users/{id}",
		Summary:     "Update a user",
		Description: "Updates only the fields present in the body.",
		Tags:        []string{"users"},
		Errors:      []int{http.StatusBadRequest, http.StatusNotFound, http.StatusConflict},
	}, h.Update)
}

type output struct {
	Body response
}

type response struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func toResponse(u user.User) response {
	return response{
		ID:   u.ID,
		Name: u.Name,
	}
}
