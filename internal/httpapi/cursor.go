package httpapi

import (
	"encoding/base64"
	"errors"

	"github.com/google/uuid"
)

// Cursors are opaque to clients: a version byte followed by the 16 bytes of
// the last item's ID, base64url-encoded. The version lets the format change
// (e.g. to also carry created_at) without misreading older cursors.
const cursorV1 byte = 1

var errInvalidCursor = errors.New("invalid cursor")

func encodeCursor(id uuid.UUID) string {
	b := make([]byte, 0, 1+len(id))
	b = append(b, cursorV1)
	b = append(b, id[:]...)
	return base64.RawURLEncoding.EncodeToString(b)
}

func decodeCursor(s string) (uuid.UUID, error) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil || len(b) != 1+len(uuid.UUID{}) || b[0] != cursorV1 {
		return uuid.Nil, errInvalidCursor
	}
	return uuid.UUID(b[1:]), nil
}
