package httpapi

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type Pinger interface {
	PingContext(ctx context.Context) error
}

type HealthHandler struct {
	db     Pinger
	logger *slog.Logger
}

func NewHealthHandler(db Pinger, logger *slog.Logger) *HealthHandler {
	return &HealthHandler{db: db, logger: logger}
}

func (h *HealthHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-liveness",
		Method:      http.MethodGet,
		Path:        "/healthz",
		Summary:     "Liveness probe",
		Description: "Reports that the process is up and serving requests.",
		Tags:        []string{"health"},
	}, h.Live)

	huma.Register(api, huma.Operation{
		OperationID: "get-readiness",
		Method:      http.MethodGet,
		Path:        "/readyz",
		Summary:     "Readiness probe",
		Description: "Reports whether the service can handle traffic, i.e. reach the database.",
		Tags:        []string{"health"},
		Errors:      []int{http.StatusServiceUnavailable},
	}, h.Ready)
}

type healthOutput struct {
	Body healthResponse
}

type healthResponse struct {
	Status string `json:"status" example:"ok"`
}

func (h *HealthHandler) Live(ctx context.Context, _ *struct{}) (*healthOutput, error) {
	return &healthOutput{Body: healthResponse{Status: "ok"}}, nil
}

func (h *HealthHandler) Ready(ctx context.Context, _ *struct{}) (*healthOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		h.logger.ErrorContext(ctx, "database ping failed", "error", err)
		return nil, huma.Error503ServiceUnavailable("database unavailable")
	}
	return &healthOutput{Body: healthResponse{Status: "ok"}}, nil
}
