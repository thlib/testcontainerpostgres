# Postgres testcontainer
A postgres docker testcontainer for go


## Example usage:

```go
initPath, err := filepath.Abs("./fixtures")
if err != nil {
    log.Fatalf("%v", err)
}
ctx := context.Background()
postgresC, conn, err := postgrestestcontainer.New(ctx, "14.5-alpine", initPath)
if err != nil {
    log.Fatalf("%v", err)
}
defer Terminate(ctx, postgresC)

fmt.Println(conn)
// Output: postgres://postgres:postgres@localhost:49156/test_db
```

## Run tests

```sh
go test
```
