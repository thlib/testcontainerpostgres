# Postgres testcontainer
A postgres docker testcontainer for go


## Example usage:

```go
initPath, err := filepath.Abs("./initdb")
if err != nil {
    log.Fatalf("%v", err)
}
ctx := context.Background()
postgresC, conn, err := testcontainerpostgres.New(ctx, "14.5-alpine",
    testcontainerpostgres.WithInit(initPath),
    testcontainerpostgres.WithDb("test_db"),
    testcontainerpostgres.WithAuth("postgres", "postgres"),
)
if err != nil {
    log.Fatalf("%v", err)
}
defer Terminate(ctx, postgresC)

fmt.Println(conn)
// Output: postgres://postgres:postgres@localhost:49156/test_db
```

## Run tests

``sh
go test
``
