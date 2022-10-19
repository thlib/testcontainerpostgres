# Postgres testcontainer
A postgres docker testcontainer for go


Example usage:

```go
schemaPath, err := filepath.Abs("example_schema.sql")
if err != nil {
    log.Fatalf("%v", err)
}
ctx := context.Background()
postgresC, conn, err := postgrestestcontainer.SetupTestDatabase(ctx, "14.5-alpine", schemaPath)
if err != nil {
    log.Fatalf("%v", err)
}
defer Terminate(ctx, postgresC)

fmt.Println(conn)
// Output: postgres://postgres:postgres@localhost:49156/test_db
```
