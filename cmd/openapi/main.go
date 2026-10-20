// Command openapi writes the API's OpenAPI spec as YAML to stdout, without starting the server or connecting to a database.
package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/Iarrova/go-starter/internal/app"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "openapi:", err)
		os.Exit(1)
	}
}

func run() error {
	gin.SetMode(gin.ReleaseMode)

	server := app.NewServer(nil, nil, slog.New(slog.DiscardHandler))

	spec, err := server.OpenAPI().YAML()
	if err != nil {
		return fmt.Errorf("encoding spec: %w", err)
	}

	_, err = os.Stdout.Write(spec)
	return err
}
