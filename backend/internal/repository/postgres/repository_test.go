// File: repository_test.go
// Purpose: Verify PostgreSQL repository local helper behavior without DB dependency.
// Module: backend/internal/repository/postgres, repository unit test layer.
// Related: repository.go helper functions.
package postgres

import (
	"errors"
	"testing"
	"time"
)

func TestNormalizePagination(t *testing.T) {
	p, ps := normalizePagination(0, 1000)
	if p != 1 || ps != 100 {
		t.Fatalf("unexpected normalized values: page=%d pageSize=%d", p, ps)
	}
}

func TestNullableString(t *testing.T) {
	if nullableString(" ") != nil {
		t.Fatalf("expected nil for blank string")
	}
	if v, ok := nullableString("x").(string); !ok || v != "x" {
		t.Fatalf("expected passthrough string")
	}
}

func TestNullableUUID(t *testing.T) {
	if nullableUUID(" ") != nil {
		t.Fatalf("expected nil for blank uuid")
	}
	if nullableUUID("admin-1") != nil {
		t.Fatalf("expected nil for invalid uuid")
	}
	if v, ok := nullableUUID("550e8400-e29b-41d4-a716-446655440000").(string); !ok || v == "" {
		t.Fatalf("expected valid uuid passthrough")
	}
}

func TestNullableTime(t *testing.T) {
	if nullableTime(time.Time{}) != nil {
		t.Fatalf("expected nil for zero time")
	}
	now := time.Now()
	if nullableTime(now) == nil {
		t.Fatalf("expected non-nil for non-zero time")
	}
}

func TestViolationHelpers(t *testing.T) {
	if !isUniqueViolation(errors.New("duplicate key value violates unique constraint")) {
		t.Fatalf("expected unique violation detection")
	}
	if !isForeignKeyViolation(errors.New("violates foreign key constraint")) {
		t.Fatalf("expected foreign key violation detection")
	}
}

func TestEnsureJSONText(t *testing.T) {
	if got := string(ensureJSONText("")); got != "{}" {
		t.Fatalf("expected empty to become {}, got %s", got)
	}
	if got := string(ensureJSONText("not-json")); got != "{}" {
		t.Fatalf("expected invalid json to become {}, got %s", got)
	}
	if got := string(ensureJSONText(`{"a":1}`)); got != `{"a":1}` {
		t.Fatalf("expected valid json unchanged, got %s", got)
	}
}
