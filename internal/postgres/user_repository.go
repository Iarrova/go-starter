package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	ent "github.com/Iarrova/go-starter/ent/generated"
	"github.com/Iarrova/go-starter/internal/user"
)

// userEmailUnique is the unique index on users.email, as named in migrations/.
const userEmailUnique = "users_email_key"

type UserRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{client: client}
}

func (r *UserRepository) Create(ctx context.Context, in user.CreateParams) (user.User, error) {
	u, err := clientFrom(ctx, r.client).User.Create().
		SetEmail(in.Email).
		SetPasswordHash(in.PasswordHash).
		SetName(in.Name).
		Save(ctx)

	switch {
	case isConstraintViolation(err, uniqueViolation, userEmailUnique):
		return user.User{}, user.ErrEmailTaken
	case err != nil:
		return user.User{}, fmt.Errorf("creating user: %w", err)
	}

	return userToDomain(u), nil
}

func (r *UserRepository) Get(ctx context.Context, id uuid.UUID) (user.User, error) {
	u, err := clientFrom(ctx, r.client).User.Get(ctx, id)

	switch {
	case ent.IsNotFound(err):
		return user.User{}, user.ErrNotFound
	case err != nil:
		return user.User{}, fmt.Errorf("retrieving user: %w", err)
	}

	return userToDomain(u), nil
}

func (r *UserRepository) Update(ctx context.Context, id uuid.UUID, in user.UpdateInput) (user.User, error) {
	u, err := clientFrom(ctx, r.client).User.UpdateOneID(id).
		SetNillableEmail(in.Email).
		SetNillableName(in.Name).
		Save(ctx)

	switch {
	case ent.IsNotFound(err):
		return user.User{}, user.ErrNotFound
	case isConstraintViolation(err, uniqueViolation, userEmailUnique):
		return user.User{}, user.ErrEmailTaken
	case err != nil:
		return user.User{}, fmt.Errorf("updating user: %w", err)
	}

	return userToDomain(u), nil
}

func userToDomain(u *ent.User) user.User {
	return user.User{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Name:         u.Name,
		CreatedAt:    u.CreatedAt.UTC(),
		UpdatedAt:    u.UpdatedAt.UTC(),
	}
}
