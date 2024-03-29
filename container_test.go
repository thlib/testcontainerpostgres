package testcontainerpostgres_test

import (
	"context"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/thlib/testcontainerpostgres"
)

// TestNew runs an example postgresql container
func TestNew(t *testing.T) {
	schemaPath, err := filepath.Abs("./initdb")
	if err != nil {
		t.Fatalf("%v", err)
	}
	ctx := context.Background()
	postgresC, conn, err := testcontainerpostgres.New(ctx, "14.5-alpine",
		testcontainerpostgres.WithInit(schemaPath),
		testcontainerpostgres.WithDb("test_db"),
		testcontainerpostgres.WithAuth("postgres", "postgres"),
	)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer testcontainerpostgres.Terminate(ctx, postgresC)

	expected := regexp.QuoteMeta("postgres://postgres:postgres@localhost:") + "[0-9]+" + regexp.QuoteMeta("/test_db")
	rx, err := regexp.Compile(expected)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if !rx.MatchString(conn) {
		t.Errorf("Expected a connection string that looks like: %v, got: %v", expected, conn)
	}
}
