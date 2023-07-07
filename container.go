// Package testcontainerpostgres provides an easy way to start a postgres testcontainer using docker
package testcontainerpostgres

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Option func(testcontainers.ContainerRequest) testcontainers.ContainerRequest

// WithInit adds a path to a folder containing sql files to be executed on startup
func WithInit(init string) Option {
	return func(req testcontainers.ContainerRequest) testcontainers.ContainerRequest {
		req.Mounts = testcontainers.Mounts(testcontainers.ContainerMount{
			Source: testcontainers.GenericBindMountSource{
				HostPath: init,
			},
			Target: testcontainers.ContainerMountTarget("/docker-entrypoint-initdb.d"),
		})
		return req
	}
}

// WithDb adds a database name to the container request
func WithDb(db string) Option {
	return func(req testcontainers.ContainerRequest) testcontainers.ContainerRequest {
		if req.Env == nil {
			req.Env = make(map[string]string)
		}
		req.Env["POSTGRES_DB"] = db
		return req
	}
}

// WithAuth adds a username and password to the container request
func WithAuth(user, pass string) Option {
	return func(req testcontainers.ContainerRequest) testcontainers.ContainerRequest {
		if req.Env == nil {
			req.Env = make(map[string]string)
		}
		req.Env["POSTGRES_USER"] = user
		req.Env["POSTGRES_PASSWORD"] = pass
		return req
	}
}

// WithEnv replaces the environment variables of the container request
func WithEnv(env map[string]string) Option {
	return func(req testcontainers.ContainerRequest) testcontainers.ContainerRequest {
		req.Env = env
		return req
	}
}

// New setup a postgres testcontainer
func New(ctx context.Context, tag string, opts ...Option) (testcontainers.Container, string, error) {
	const (
		name = "test_db"
		user = "postgres"
		pass = "postgres"
	)

	// Create PostgreSQL container request
	req := testcontainers.ContainerRequest{
		Image: "postgres:" + tag,
		Env: map[string]string{
			"POSTGRES_DB":       name,
			"POSTGRES_USER":     user,
			"POSTGRES_PASSWORD": pass,
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
	}

	// Start PostgreSQL container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to start postgres container: %w", err)
	}

	// Get host and port of PostgreSQL container
	host, err := container.Host(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get host: %w", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get port: %w", err)
	}

	conn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", user, pass, host, port.Port(), name)

	// Create db connection string and connect
	return container, conn, nil
}

// Terminate terminates the container in a defer friendly way
func Terminate(ctx context.Context, c testcontainers.Container) {
	err := c.Terminate(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to terminate container: %v", err))
	}
}
