package docker

import (
	"context"
	"fmt"

	"github.com/petrovskiborislav/docker-cli/logger"
)

// Client provides interactions with the docker SDK.
type Client interface {
	ServiceProvisioning(ctx context.Context, container Container) error
	ServiceDecommissioning(ctx context.Context, container Container) error
}

type client struct {
	logger  logger.Logger
	actions Actions
}

// NewClient creates a new docker client.
func NewClient(logger logger.Logger, actions Actions) Client {
	return &client{logger: logger, actions: actions}
}

// ServiceProvisioning creates and run a service within a container with isolated network.
func (c client) ServiceProvisioning(ctx context.Context, container Container) error {
	err := c.pullImageIfNotExists(ctx, container.Image)
	if err != nil {
		return err
	}

	networkName := fmt.Sprintf("%s-network", container.Name)
	newNetworkID, err := c.actions.CreateNetwork(ctx, networkName)
	if err != nil {
		return err
	}
	c.logger.Info("Successfully created network %s \n", networkName)

	containerID, err := c.actions.CreateContainerWithNetwork(ctx, container.Image, container.Name, newNetworkID, container.EnvironmentVars)
	if err != nil {
		return err
	}
	c.logger.Info("Successfully created container %s\n", container.Name)

	err = c.actions.StartContainer(ctx, containerID)
	if err != nil {
		return err
	}

	c.logger.Info("Successfully started container %s\n", container.Name)

	return nil
}

// ServiceDecommissioning stops and removes a service container and its isolated network.
func (c client) ServiceDecommissioning(ctx context.Context, container Container) error {
	containerID, err := c.actions.StopContainer(ctx, container.Name)
	if err != nil {
		return err
	}

	if containerID == "" {
		c.logger.Warn("Container not found skipping\n")
		return nil
	}

	c.logger.Info("Successfully stopped container %s\n", container.Name)

	err = c.actions.RemoveContainer(ctx, containerID)
	if err != nil {
		return err
	}
	c.logger.Info("Successfully removed container %s\n", container.Name)

	err = c.actions.RemoveNetwork(ctx, container.Name)
	if err != nil {
		return err
	}
	c.logger.Info("Successfully removed network %s\n", container.Name)

	return nil
}

func (c client) pullImageIfNotExists(ctx context.Context, image string) error {
	exists, err := c.actions.CheckIfImageExists(ctx, image)
	if err != nil {
		return err
	}

	if exists {
		c.logger.Warn("Image already exists skipping\n")
		return nil
	}

	err = c.actions.PullImage(ctx, image)
	if err != nil {
		return err
	}
	c.logger.Info("Successfully pulled image %s\n", image)

	return nil
}
