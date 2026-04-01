package page

import (
	"fmt"
	"strings"
	"time"
)

const DefaultSize = 20

type CursorRequest struct {
	Cursor string
	Size   int
}

func (r CursorRequest) EffectiveSize() int {
	if r.Size <= 0 || r.Size > 100 {
		return DefaultSize
	}
	return r.Size
}

type CursorPage[T any] struct {
	Items      []T    `json:"items"`
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
}

// Cursor format: "2006-01-02T15:04:05.999999Z_<id>"

func EncodeCursor(createdAt time.Time, id string) string {
	return createdAt.UTC().Format(time.RFC3339Nano) + "_" + id
}

func DecodeCursor(cursor string) (createdAt time.Time, id string, err error) {
	if cursor == "" {
		return time.Time{}, "", nil
	}

	idx := strings.LastIndex(cursor, "_")
	if idx < 0 {
		return time.Time{}, "", fmt.Errorf("invalid cursor format")
	}

	createdAt, err = time.Parse(time.RFC3339Nano, cursor[:idx])
	if err != nil {
		return time.Time{}, "", fmt.Errorf("invalid cursor timestamp: %w", err)
	}

	id = cursor[idx+1:]
	if id == "" {
		return time.Time{}, "", fmt.Errorf("invalid cursor: empty id")
	}

	return createdAt, id, nil
}
