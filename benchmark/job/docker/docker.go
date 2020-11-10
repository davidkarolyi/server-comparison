package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// Docker implements the IDocker interface,
// and uses the official docker client to communicate with the docker daemon.
type Docker struct {
	serverName        string
	client            *client.Client
	context           context.Context
	cancel            context.CancelFunc
	serverContainerID string
}

// New initialises a new Docker service.
func New(serverName string) (*Docker, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	return &Docker{
		serverName: serverName,
		client:     cli,
		context:    ctx,
		cancel:     cancel,
	}, nil
}

// BuildServerImage builds an image with the named after the server.
func (docker *Docker) BuildServerImage() error {
	buildCommand := fmt.Sprintf("docker build -t %[1]s -f %[1]s/Dockerfile .", docker.serverName)
	cmd := exec.CommandContext(docker.context, "bash", "-c", buildCommand)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// StartServerContainer will start a container with the given server inside.
// and returns the container ID
func (docker *Docker) StartServerContainer() (chan *Stats, error) {
	resp, err := docker.client.ContainerCreate(
		docker.context,
		docker.containerConfig(),
		docker.hostConfig(), nil, "")
	if err != nil {
		return nil, err
	}
	docker.serverContainerID = resp.ID

	err = docker.client.ContainerStart(
		docker.context,
		docker.serverContainerID,
		types.ContainerStartOptions{},
	)
	if err != nil {
		return nil, err
	}

	statsChannel, err := docker.setupMonitoring()
	if err != nil {
		return nil, err
	}

	return statsChannel, nil
}

// Clean cleans up container(s) created by the instance.
func (docker *Docker) Clean() error {
	err := docker.client.ContainerRemove(
		context.Background(),
		docker.serverContainerID,
		types.ContainerRemoveOptions{Force: true},
	)
	if err != nil {
		return err
	}
	docker.cancel()
	return nil
}

func (docker *Docker) setupMonitoring() (chan *Stats, error) {
	var err error = nil
	statsChannel := make(chan *Stats)
	writer := newStatsWriter(statsChannel)

	go func() {
		stats, localErr := docker.client.ContainerStats(
			docker.context,
			docker.serverContainerID,
			true,
		)
		io.Copy(writer, stats.Body)
		stats.Body.Close()
		err = localErr
	}()
	if err != nil {
		return nil, err
	}

	return statsChannel, nil
}

func (docker *Docker) hostConfig() *container.HostConfig {
	portMap := make(nat.PortMap)
	portMap["3000/tcp"] = []nat.PortBinding{
		{
			HostIP:   "0.0.0.0",
			HostPort: "3000",
		},
	}

	return &container.HostConfig{
		PortBindings: portMap,
	}
}

func (docker *Docker) containerConfig() *container.Config {
	return &container.Config{
		Image: docker.serverName,
		ExposedPorts: nat.PortSet{
			"3000/tcp": struct{}{},
		},
	}
}
