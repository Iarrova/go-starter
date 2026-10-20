package httpapi

import (
	"context"
	"errors"
	"maps"
	"net/http"
	"slices"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"

	"github.com/Iarrova/go-starter/internal/apperr"
)

// statusClientClosedRequest is nginx's non-standard 499, used when the client
// disconnects before the response; net/http has no constant for it.
const statusClientClosedRequest = 499

func init() {
	// If a handler returns a plain Go error (anything that isn't a huma.StatusError), Huma treats it as a 500 by default.
	// This acts as an override to map domain errors to their proper status.
	defaultNewError := huma.NewErrorWithContext

	// Override huma.NewErrorWithContext with our own function.
	huma.NewErrorWithContext = func(ctx huma.Context, status int, msg string, errs ...error) huma.StatusError {
		for _, err := range errs {
			var aerr *apperr.Error
			if errors.As(err, &aerr) {
				// Look for the corresponding HTTP status for the kind. Internal and unknown pass through.
				if status, ok := statusFor(aerr.Kind); ok {
					return huma.NewError(status, aerr.Message, fieldErrors(aerr.Fields)...)
				}
			}
		}

		// Hide and log server errors
		if status >= http.StatusInternalServerError {
			timedOut, canceled := false, false
			for _, err := range errs {
				if err == nil {
					continue
				}
				if ctx != nil {
					_ = humagin.Unwrap(ctx).Error(err) // Keep the real error for logging
				}
				if errors.Is(err, context.DeadlineExceeded) {
					timedOut = true
				}
				if errors.Is(err, context.Canceled) {
					canceled = true
				}
			}
			if timedOut {
				return huma.NewError(http.StatusGatewayTimeout, "request timed out")
			}
			if canceled {
				// The client is gone and won't see this; it keeps disconnects out of 5xx logs and metrics.
				return huma.NewError(statusClientClosedRequest, "client closed request")
			}
			return huma.NewError(status, msg)
		}

		return defaultNewError(ctx, status, msg, errs...)
	}
}

// statusFor returns the HTTP status for a client-facing kind.
// It reports false for KindInternal and any kind without a mapping.
func statusFor(kind apperr.Kind) (int, bool) {
	switch kind {
	case apperr.KindUnauthorized:
		return http.StatusUnauthorized, true
	case apperr.KindForbidden:
		return http.StatusForbidden, true
	case apperr.KindNotFound:
		return http.StatusNotFound, true
	case apperr.KindConflict:
		return http.StatusConflict, true
	case apperr.KindInvalid:
		return http.StatusUnprocessableEntity, true
	default:
		return 0, false
	}
}

// fieldErrors converts domain field errors into Huma error details, sorted by field so responses are stable.
// Domain validation covers request bodies only.
func fieldErrors(fields map[string]string) []error {
	errs := make([]error, 0, len(fields))
	for _, field := range slices.Sorted(maps.Keys(fields)) {
		errs = append(errs, &huma.ErrorDetail{Location: "body." + field, Message: fields[field]})
	}
	return errs
}
