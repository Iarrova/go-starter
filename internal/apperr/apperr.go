package apperr

// Kind classifies an error so transports (HTTP, gRPC, CLI) can map it without knowing about specific domain errors.
// Errors inside each domain package define their errors using these types.
type Kind int

const (
	KindInternal Kind = iota // Zero value, so an Error without a Kind is hidden from clients
	KindUnauthorized
	KindForbidden
	KindNotFound
	KindConflict
	KindInvalid
)

// Error is an expected, client-facing error. Message is safe to return to callers.
type Error struct {
	Kind    Kind
	Message string
	Fields  map[string]string // set only for KindInvalid with field-level errors
}

func (e *Error) Error() string { return e.Message }

func Unauthorized(msg string) *Error { return &Error{Kind: KindUnauthorized, Message: msg} }
func Forbidden(msg string) *Error    { return &Error{Kind: KindForbidden, Message: msg} }
func NotFound(msg string) *Error     { return &Error{Kind: KindNotFound, Message: msg} }
func Conflict(msg string) *Error     { return &Error{Kind: KindConflict, Message: msg} }
func Invalid(msg string) *Error      { return &Error{Kind: KindInvalid, Message: msg} }

// InvalidFields is a KindInvalid error carrying per-field messages, keyed by field name.
func InvalidFields(fields map[string]string) *Error {
	return &Error{Kind: KindInvalid, Message: "invalid input", Fields: fields}
}
