package user

import (
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Iarrova/go-starter/internal/validation"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Name         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func normalizeEmail(s string) string { return strings.ToLower(strings.TrimSpace(s)) }

func checkEmail(v *validation.Validator, email string) {
	addr, err := mail.ParseAddress(email)
	v.Check(err == nil && addr.Address == email, "email", "must be a valid email address")
}

func checkName(v *validation.Validator, name string) {
	n := utf8.RuneCountInString(name)
	v.Check(n > 0 && n <= 200, "name", "is required and cannot be more than 200 characters")
}

// checkPassword counts bytes, not characters: bcrypt rejects passwords over 72
// bytes, and 64 multi-byte characters could exceed that.
func checkPassword(v *validation.Validator, password string) {
	v.Check(len(password) >= 8 && len(password) <= 64, "password", "must be between 8 and 64 characters")
}
