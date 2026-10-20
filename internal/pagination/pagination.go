package pagination

import "github.com/google/uuid"

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

type Request struct {
	Cursor *uuid.UUID // ID of the last item on the previous page
	Limit  int
}

func (r Request) Normalize() Request {
	if r.Limit <= 0 {
		r.Limit = DefaultLimit
	}
	r.Limit = min(r.Limit, MaxLimit)
	return r
}

type Page[T any] struct {
	Items      []T
	NextCursor *uuid.UUID
}

// Build trims the extra row fetched to detect whether a next page exists.
func Build[T any](items []T, limit int, id func(T) uuid.UUID) Page[T] {
	if len(items) <= limit {
		return Page[T]{Items: items}
	}
	items = items[:limit]
	next := id(items[len(items)-1])
	return Page[T]{Items: items, NextCursor: &next}
}
