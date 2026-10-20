package httpapi

import (
	"log/slog"
	"net/http"
	"path"
	"reflect"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	engine *gin.Engine
	api    huma.API
}

// Handler registers its operations on the Huma API.
type Handler interface {
	Register(api huma.API)
}

// Timeouts that depend on each other live together, so changing requestTimeout
// keeps the others consistent.
const (
	// requestTimeout is the deadline for handler work (database queries etc.).
	requestTimeout = 10 * time.Second
	// writeTimeout leaves time to send the 504 once requestTimeout has passed,
	// instead of net/http closing the connection first.
	writeTimeout = requestTimeout + 5*time.Second
	// ShutdownTimeout lets requests in flight when shutdown starts run for
	// their full requestTimeout.
	ShutdownTimeout = requestTimeout + 5*time.Second
)

// NewHTTPServer registers health probes at the root, since they belong to the
// infrastructure rather than a versioned API, and handlers under /api/v1.
func NewHTTPServer(logger *slog.Logger, health *HealthHandler, handlers ...Handler) *HTTPServer {
	engine := gin.New()
	engine.Use(
		requestID(),
		timeout(requestTimeout),
		logRequests(logger, "/healthz", "/readyz"),
		recoverPanics(logger),
	)

	// Huma also limits request bodies (1 MB per operation by default) and
	// rejects unknown fields.
	config := huma.DefaultConfig("go-starter", "1.0.0")
	config.CreateHooks = nil // Don't add $schema fields and Link headers to responses
	config.Components.Schemas = huma.NewMapRegistry("#/components/schemas/", schemaName)
	root := humagin.New(engine, config)
	health.Register(root)

	v1 := huma.NewGroup(root, "/api/v1")
	for _, handler := range handlers {
		handler.Register(v1)
	}

	return &HTTPServer{
		engine: engine,
		api:    root,
	}
}

func (s *HTTPServer) Handler() http.Handler {
	return s.engine
}

// Server returns an *http.Server that serves s on addr.
func (s *HTTPServer) Server(addr string) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           s.engine,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       60 * time.Second,
	}
}

// OpenAPI returns the spec built from every registered operation.
func (s *HTTPServer) OpenAPI() *huma.OpenAPI {
	return s.api.OpenAPI()
}

// schemaName works like huma.DefaultSchemaNamer, but prefixes types from
// resource sub-packages (httpapi/postapi, httpapi/userapi, ...) with the
// resource name. Huma drops the package from schema names, so without this
// postapi.createRequest and userapi.createRequest would collide.
func schemaName(t reflect.Type, hint string) string {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Name() == "" || t.PkgPath() == "" {
		return huma.DefaultSchemaNamer(t, hint)
	}

	// A generic type's Name includes its fully qualified type arguments, e.g.
	// PageResponse[github.com/.../httpapi/postapi.response].
	full := t.PkgPath() + "." + t.Name()
	var b strings.Builder
	for _, part := range strings.FieldsFunc(full, func(r rune) bool {
		return r == '[' || r == ']' || r == '*' || r == ','
	}) {
		dot := strings.LastIndex(part, ".")
		if dot >= 0 {
			if pkg := part[:dot]; path.Base(path.Dir(pkg)) == "httpapi" {
				b.WriteString(upperFirst(strings.TrimSuffix(path.Base(pkg), "api")))
			}
		}
		b.WriteString(upperFirst(part[dot+1:]))
	}
	return b.String()
}

func upperFirst(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}
