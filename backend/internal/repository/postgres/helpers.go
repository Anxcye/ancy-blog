// File: helpers.go
// Purpose: Provide shared pagination, nullable conversion, and database error helpers.
// Module: backend/internal/repository/postgres, persistence utility layer.
// Related: all postgres repository method files.
package postgres

import (
	"encoding/json"
	"strings"
	"time"
)

func normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func nullableString(v string) any {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	return v
}

func nullableUUID(v string) any {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	return v
}

func nullableTime(t time.Time) any {
	if t.IsZero() {
		return nil
	}
	return t
}

func isUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}

func isForeignKeyViolation(err error) bool {
	return strings.Contains(err.Error(), "violates foreign key constraint")
}

func ensureJSONText(v string) []byte {
	if strings.TrimSpace(v) == "" {
		return []byte("{}")
	}
	raw := []byte(v)
	if json.Valid(raw) {
		return raw
	}
	return []byte("{}")
}
