package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"

	dockerClient "github.com/docker/docker/client"
)

//go:generate mockery --name=Actions --structname mockActions --filename mock_actions_test.go --outpkg=docker_test --output=.

// Actions represents a set of actions that can be performed on docker engine.
// This interface is used to mock docker SDK in tests.
type Actions interface {
	CheckIfImageExists(ctx context.Context, imageName string) (bool, error)
	PullImage(ctx context.Context, imageName string) error
	CreateNetwork(ctx context.Context, networkName string) (string, error)
	CreateContainerWithNetwork(ctx context.Context, imageName, containerName, networkID string, envs []string) (string, error)
	StartContainer(ctx context.Context, containerID string) error
	StopContainer(ctx context.Context, containerName string) (string, error)
	RemoveContainer(ctx context.Context, containerID string) error
	RemoveNetwork(ctx context.Context, containerName string) error
}

type actions struct {
	client dockerClient.APIClient
}

// NewActions creates a new instance of Actions.
func NewActions(client dockerClient.APIClient) Actions {
	return actions{client: client}
}

// CheckIfImageExists checks if an image exists in the local docker.
func (a actions) CheckIfImageExists(ctx context.Context, imageName string) (bool, error) {
	filter := filters.NewArgs()
	filter.Add("reference", imageName)

	imageListOptions := types.ImageListOptions{Filters: filter}
	images, err := a.client.ImageList(ctx, imageListOptions)
	if err != nil {
		return false, err
	}

	return len(images) == 1, nil
}

// PullImage pulls an image from the docker hub.
func (a actions) PullImage(ctx context.Context, imageName string) error {
	reader, err := a.client.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()
	io.Copy(os.Stdout, reader)

	return nil
}

// CreateNetwork creates a new network.
func (a actions) CreateNetwork(ctx context.Context, networkName string) (string, error) {
	network, err := a.client.NetworkCreate(ctx, networkName, types.NetworkCreate{})
	if err != nil {
		return "", err
	}
	return network.ID, nil
}

// CreateContainerWithNetwork creates a new container and connects it to the specified network.
func (a actions) CreateContainerWithNetwork(ctx context.Context, imageName, containerName, networkID string, env []string) (string, error) {
	containerConfig := &container.Config{Image: imageName, Env: env}
	createdContainer, err := a.client.ContainerCreate(ctx, containerConfig, nil, nil, nil, containerName)
	if err != nil {
		return "", err
	}

	err = a.client.NetworkConnect(ctx, networkID, createdContainer.ID, nil)
	if err != nil {
		return "", err
	}

	return createdContainer.ID, nil
}

// StartContainer starts a container.
func (a actions) StartContainer(ctx context.Context, containerID string) error {
	return a.client.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
}

// StopContainer stops a container.
func (a actions) StopContainer(ctx context.Context, containerName string) (string, error) {
	filter := filters.NewArgs()
	filter.Add("name", containerName)

	containerListOptions := types.ContainerListOptions{Filters: filter}
	containers, err := a.client.ContainerList(ctx, containerListOptions)
	if err != nil {
		return "", err
	}

	if len(containers) == 0 {
		return "", fmt.Errorf("container %s not found", containerName)
	}

	containerID := containers[0].ID

	return containerID, a.client.ContainerStop(ctx, containerID, nil)
}

// RemoveContainer removes a container.
func (a actions) RemoveContainer(ctx context.Context, containerID string) error {
	return a.client.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{})
}

// RemoveNetwork removes a network.
func (a actions) RemoveNetwork(ctx context.Context, containerName string) error {
	filter := filters.NewArgs()
	filter.Add("name", fmt.Sprintf("%s-network", containerName))

	networkListOptions := types.NetworkListOptions{Filters: filter}
	networks, err := a.client.NetworkList(ctx, networkListOptions)
	if err != nil {
		return err
	}

	return a.client.NetworkRemove(ctx, networks[0].ID)
}
