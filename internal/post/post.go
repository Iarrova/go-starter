package post

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Iarrova/go-starter/internal/validation"
	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Title     string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func checkTitle(v *validation.Validator, title string) {
	n := utf8.RuneCountInString(title)
	v.Check(n > 0 && n <= 200, "title", "must be between 1 and 200 characters")
}

// checkBody rejects whitespace-only bodies without trimming the body itself,
// since leading and trailing whitespace can be meaningful (e.g. Markdown).
func checkBody(v *validation.Validator, body string) {
	v.Check(strings.TrimSpace(body) != "", "body", "must not be empty")
}
