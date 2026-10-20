// Package app wires repositories, services and handlers into the HTTP server.
package app

import (
	"log/slog"

	ent "github.com/Iarrova/go-starter/ent/generated"
	"github.com/Iarrova/go-starter/internal/httpapi"
	"github.com/Iarrova/go-starter/internal/httpapi/postapi"
	"github.com/Iarrova/go-starter/internal/httpapi/userapi"
	"github.com/Iarrova/go-starter/internal/post"
	"github.com/Iarrova/go-starter/internal/postgres"
	"github.com/Iarrova/go-starter/internal/user"
)

func NewServer(client *ent.Client, db httpapi.Pinger, logger *slog.Logger) *httpapi.HTTPServer {
	userService := user.NewService(postgres.NewUserRepository(client))
	postService := post.NewService(postgres.NewPostRepository(client))

	return httpapi.NewHTTPServer(
		logger,
		httpapi.NewHealthHandler(db, logger),
		userapi.New(userService),
		postapi.New(postService),
	)
}
