package post

import "github.com/Iarrova/go-starter/internal/apperr"

var (
	ErrNotFound     = apperr.NotFound("post not found")
	ErrUserNotFound = apperr.InvalidFields(map[string]string{"user_id": "user does not exist"})
)
