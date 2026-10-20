package httpapi

import (
	"github.com/danielgtaylor/huma/v2"

	"github.com/Iarrova/go-starter/internal/pagination"
)

// PageQuery is the cursor/limit query embedded by every list endpoint.
// The limit tags must match pagination.DefaultLimit and MaxLimit, see pagination_test.go.
type PageQuery struct {
	Cursor string `query:"cursor" doc:"next_cursor from the previous page"`
	Limit  int    `query:"limit" default:"20" minimum:"1" maximum:"100" doc:"Page size"`

	request pagination.Request // set by Resolve
}

// Resolve decodes the cursor. Huma calls it after parsing and validating the
// query, and responds 422 with the returned errors.
func (q *PageQuery) Resolve(huma.Context) []error {
	q.request = pagination.Request{Limit: q.Limit}
	if q.Cursor == "" {
		return nil
	}

	id, err := decodeCursor(q.Cursor)
	if err != nil {
		return []error{&huma.ErrorDetail{Location: "query.cursor", Message: err.Error(), Value: q.Cursor}}
	}
	q.request.Cursor = &id
	return nil
}

func (q PageQuery) Request() pagination.Request {
	return q.request
}

type PageResponse[T any] struct {
	Items      []T     `json:"items"`
	NextCursor *string `json:"next_cursor,omitempty"`
}

func ToPageResponse[D, R any](p pagination.Page[D], mapItem func(D) R) PageResponse[R] {
	items := make([]R, len(p.Items)) // never nil, so it encodes as [] not null
	for i, item := range p.Items {
		items[i] = mapItem(item)
	}

	resp := PageResponse[R]{Items: items}
	if p.NextCursor != nil {
		cursor := encodeCursor(*p.NextCursor)
		resp.NextCursor = &cursor
	}
	return resp
}
