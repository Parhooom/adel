package judge

import (
	"github.com/docker/docker/client"

	"context"
	"fmt"
	"log"
	"time"
)

type DockerClient struct {
	client *client.Client
	logger *log.Logger
}

func NewDockerClient(logger *log.Logger) *DockerClient {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(fmt.Errorf("failed to create docker client: %w", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = cli.Ping(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to ping docker: %w", err))
	}

	return &DockerClient{
		client: cli,
		logger: logger,
	}
}
