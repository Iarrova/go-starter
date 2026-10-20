package postgres

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"

	ent "github.com/Iarrova/go-starter/ent/generated"
	entpost "github.com/Iarrova/go-starter/ent/generated/post"
	"github.com/Iarrova/go-starter/internal/pagination"
	"github.com/Iarrova/go-starter/internal/post"
)

// postUserFK is the posts.user_id foreign key, as named in migrations/.
const postUserFK = "posts_users_posts"

type PostRepository struct {
	client *ent.Client
}

func NewPostRepository(client *ent.Client) *PostRepository {
	return &PostRepository{client: client}
}

func (r *PostRepository) Create(ctx context.Context, in post.CreateParams) (post.Post, error) {
	p, err := clientFrom(ctx, r.client).Post.Create().
		SetUserID(in.UserID).
		SetTitle(in.Title).
		SetBody(in.Body).
		Save(ctx)

	switch {
	case isConstraintViolation(err, foreignKeyViolation, postUserFK):
		return post.Post{}, post.ErrUserNotFound
	case err != nil:
		return post.Post{}, fmt.Errorf("creating post: %w", err)
	}

	return postToDomain(p), nil
}

func (r *PostRepository) Get(ctx context.Context, id uuid.UUID) (post.Post, error) {
	p, err := clientFrom(ctx, r.client).Post.Get(ctx, id)

	switch {
	case ent.IsNotFound(err):
		return post.Post{}, post.ErrNotFound
	case err != nil:
		return post.Post{}, fmt.Errorf("retrieving post: %w", err)
	}

	return postToDomain(p), nil
}

func (r *PostRepository) Update(ctx context.Context, id uuid.UUID, in post.UpdateInput) (post.Post, error) {
	p, err := clientFrom(ctx, r.client).Post.UpdateOneID(id).
		SetNillableTitle(in.Title).
		SetNillableBody(in.Body).
		Save(ctx)

	switch {
	case ent.IsNotFound(err):
		return post.Post{}, post.ErrNotFound
	case err != nil:
		return post.Post{}, fmt.Errorf("updating post: %w", err)
	}

	return postToDomain(p), nil
}

func (r *PostRepository) List(ctx context.Context, f post.ListFilter, p pagination.Request) (pagination.Page[post.Post], error) {
	q := clientFrom(ctx, r.client).Post.Query()
	if f.UserID != uuid.Nil {
		q = q.Where(entpost.UserID(f.UserID))
	}
	if p.Cursor != nil {
		q = q.Where(entpost.IDLT(*p.Cursor))
	}

	rows, err := q.Order(entpost.ByID(sql.OrderDesc())).Limit(p.Limit + 1).All(ctx)
	if err != nil {
		return pagination.Page[post.Post]{}, fmt.Errorf("listing posts: %w", err)
	}

	posts := make([]post.Post, len(rows))
	for i, row := range rows {
		posts[i] = postToDomain(row)
	}

	return pagination.Build(posts, p.Limit, func(p post.Post) uuid.UUID { return p.ID }), nil
}

func postToDomain(p *ent.Post) post.Post {
	return post.Post{
		ID:        p.ID,
		UserID:    p.UserID,
		Title:     p.Title,
		Body:      p.Body,
		CreatedAt: p.CreatedAt.UTC(),
		UpdatedAt: p.UpdatedAt.UTC(),
	}
}
