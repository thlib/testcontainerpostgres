// Package postgrescontainer provides an easy way to start a postgres testcontainer using docker
package postgrescontainer

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// SetupTestDatabase setup a postgres testcontainer
func SetupTestDatabase(ctx context.Context, init string) (testcontainers.Container, *pgxpool.Pool, error) {
	const (
		name = "test_db"
		user = "postgres"
		pass = "postgres"
	)

	// Create PostgreSQL container request
	req := testcontainers.ContainerRequest{
		Image: "postgres:14.5-alpine",
		Env: map[string]string{
			"POSTGRES_DB":       name,
			"POSTGRES_USER":     user,
			"POSTGRES_PASSWORD": pass,
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
	}
	if init != "" {
		req.Mounts = testcontainers.Mounts(testcontainers.ContainerMount{
			Source: testcontainers.GenericBindMountSource{
				HostPath: init,
			},
			Target: testcontainers.ContainerMountTarget("/docker-entrypoint-initdb.d/init.sql"),
		})
	}

	// Start PostgreSQL container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	// Get host and port of PostgreSQL container
	host, err := container.Host(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get host: %w", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get port: %w", err)
	}

	conn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", user, pass, host, port.Port(), name)
	pool, err := pgxpool.Connect(ctx, conn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	// Create db connection string and connect
	return container, pool, nil
}

// Terminate terminates the container in a defer friendly way
func Terminate(ctx context.Context, c testcontainers.Container) {
	err := c.Terminate(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to terminate container: %v", err))
	}
}