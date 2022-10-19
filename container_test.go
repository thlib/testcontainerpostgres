package testcontainerpostgres

import (
	"context"
	"path/filepath"
	"regexp"
	"testing"
)

// TestNew runs an example postgresql container
func TestNew(t *testing.T) {
	schemaPath, err := filepath.Abs("example_schema.sql")
	if err != nil {
		t.Fatalf("%v", err)
	}
	ctx := context.Background()
	postgresC, conn, err := New(ctx, "14.5-alpine", schemaPath)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer Terminate(ctx, postgresC)

	expected := regexp.QuoteMeta("postgres://postgres:postgres@localhost:") + "[0-9]+" + regexp.QuoteMeta("/test_db")
	rx, err := regexp.Compile(expected)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if !rx.MatchString(conn) {
		t.Errorf("Expected a connection string that looks like: %v, got: %v", expected, conn)
	}
}
