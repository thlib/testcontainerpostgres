package postgrestestcontainer

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
)

// ExampleContainer runs an example postgresql container
func ExampleContainer() {
	schemaPath, err := filepath.Abs("example_schema.sql")
	if err != nil {
		log.Fatalf("%v", err)
	}
	ctx := context.Background()
	postgresC, conn, err := SetupTestDatabase(ctx, "14.5-alpine", schemaPath)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer Terminate(ctx, postgresC)

	fmt.Println(conn)
	// Output: postgres://postgres:postgres@localhost:49156/test_db
}
