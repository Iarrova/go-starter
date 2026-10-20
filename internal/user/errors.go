package user

import "github.com/Iarrova/go-starter/internal/apperr"

var (
	ErrNotFound   = apperr.NotFound("user not found")
	ErrEmailTaken = apperr.Conflict("email already registered")
)
