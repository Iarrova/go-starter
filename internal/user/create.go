package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/Iarrova/go-starter/internal/validation"
)

type CreateInput struct {
	Email    string
	Password string
	Name     string
}

func (in CreateInput) normalize() CreateInput {
	in.Email = normalizeEmail(in.Email)
	in.Name = strings.TrimSpace(in.Name)
	return in
}

func (in CreateInput) validate() error {
	var v validation.Validator
	checkEmail(&v, in.Email)
	checkName(&v, in.Name)
	checkPassword(&v, in.Password)
	return v.Err()
}

type CreateParams struct {
	Email        string
	PasswordHash string
	Name         string
}

func (s *Service) Create(ctx context.Context, in CreateInput) (User, error) {
	in = in.normalize()
	if err := in.validate(); err != nil {
		return User{}, err
	}

	passwordHash, err := hashPassword(in.Password)
	if err != nil {
		return User{}, fmt.Errorf("hashing password: %w", err)
	}

	return s.repository.Create(ctx, CreateParams{
		Email:        in.Email,
		PasswordHash: passwordHash,
		Name:         in.Name,
	})
}
