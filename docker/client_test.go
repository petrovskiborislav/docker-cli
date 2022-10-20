package docker_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/petrovskiborislav/docker-cli/docker"
)

//go:generate mockery --name=Logger --structname mockLogger --filename mock_logger_test.go --outpkg=docker_test --output=. --srcpkg=github.com/petrovskiborislav/docker-cli/logger

type clientTestSuite struct {
	suite.Suite
	actions *mockActions
	logger  *mockLogger
	sut     docker.Client
}

func (s *clientTestSuite) SetupTest() {
	s.actions = &mockActions{}
	s.logger = &mockLogger{}
	s.sut = docker.NewClient(s.logger, s.actions)
}

func (s *clientTestSuite) AfterTest(suiteName string, testName string) {
	s.logger.AssertExpectations(s.T())
	s.actions.AssertExpectations(s.T())
}

func TestSuite_Client(t *testing.T) {
	suite.Run(t, &clientTestSuite{})
}

func (s *clientTestSuite) TestServiceProvisioning_WhenImageExists_ThenSuccess() {
	// Arrange
	ctx := context.Background()
	container := docker.Container{Name: "name", Image: "image"}
	networkID := "networkID"
	containerID := "containerID"
	networkName := fmt.Sprintf("%s-network", container.Name)

	s.actions.On("CheckIfImageExists", ctx, container.Image).Return(true, nil)
	s.actions.On("CreateNetwork", ctx, networkName).Return(networkID, nil)
	s.actions.On("CreateContainerWithNetwork", ctx, container.Image, container.Name, networkID, container.EnvironmentVars).Return(containerID, nil)
	s.actions.On("StartContainer", ctx, "containerID").Return(nil)

	s.logger.On("Warn", "Image already exists skipping\n").Return().Once()
	s.logger.On("Info", "Successfully created network %s \n", networkName).Return().Once()
	s.logger.On("Info", "Successfully created container %s\n", container.Name).Return().Once()
	s.logger.On("Info", "Successfully started container %s\n", container.Name).Return().Once()

	// Act
	err := s.sut.ServiceProvisioning(ctx, container)

	// Assert
	s.NoError(err)
}

func (s *clientTestSuite) TestServiceProvisioning_WhenImageDoesNotExists_ThenSuccess() {
	// Arrange
	ctx := context.Background()
	container := docker.Container{Name: "name", Image: "image"}
	networkID := "networkID"
	containerID := "containerID"
	networkName := fmt.Sprintf("%s-network", container.Name)

	s.actions.On("CheckIfImageExists", ctx, container.Image).Return(false, nil)
	s.actions.On("PullImage", ctx, container.Image).Return(nil)
	s.actions.On("CreateNetwork", ctx, networkName).Return(networkID, nil)
	s.actions.On("CreateContainerWithNetwork", ctx, container.Image, container.Name, networkID, container.EnvironmentVars).Return(containerID, nil)
	s.actions.On("StartContainer", ctx, "containerID").Return(nil)

	s.logger.On("Info", "Successfully pulled image %s\n", container.Image).Return().Once()
	s.logger.On("Info", "Successfully created network %s \n", networkName).Return().Once()
	s.logger.On("Info", "Successfully created container %s\n", container.Name).Return().Once()
	s.logger.On("Info", "Successfully started container %s\n", container.Name).Return().Once()

	// Act
	err := s.sut.ServiceProvisioning(ctx, container)

	// Assert
	s.NoError(err)
}

func (s *clientTestSuite) TestServiceProvisioning_WhenErrorOccursOnCheckingImageExists_ThenFailure() {
	// Arrange
	ctx := context.Background()
	container := docker.Container{Name: "name", Image: "image"}

	s.actions.On("CheckIfImageExists", ctx, container.Image).Return(false, errors.New("error"))

	// Act
	err := s.sut.ServiceProvisioning(ctx, container)

	// Assert
	s.Error(err)
}

func (s *clientTestSuite) TestServiceProvisioning_WhenErrorOccursOnPullingImage_ThenFailure() {
	// Arrange
	ctx := context.Background()
	container := docker.Container{Name: "name", Image: "image"}

	s.actions.On("CheckIfImageExists", ctx, container.Image).Return(false, nil)
	s.actions.On("PullImage", ctx, container.Image).Return(errors.New("error"))

	// Act
	err := s.sut.ServiceProvisioning(ctx, container)

	// Assert
	s.Error(err)
}

func (s *clientTestSuite) TestServiceProvisioning_WhenErrorOccursOnCreationOfNetwork_ThenFailure() {
	// Arrange
	ctx := context.Background()
	container := docker.Container{Name: "name", Image: "image"}
	networkName := fmt.Sprintf("%s-network", container.Name)

	s.actions.On("CheckIfImageExists", ctx, container.Image).Return(true, nil)
	s.actions.On("CreateNetwork", ctx, networkName).Return("", errors.New("error"))

	s.logger.On("Warn", "Image already exists skipping\n").Return().Once()

	// Act
	err := s.sut.ServiceProvisioning(ctx, container)

	// Assert
	s.Error(err)
}

func (s *clientTestSuite) TestServiceProvisioning_WhenErrorOccursOnCreatingContainerWithNetwork_ThenFailure() {
	// Arrange
	ctx := context.Background()
	container := docker.Container{Name: "name", Image: "image"}
	networkID := "networkID"
	networkName := fmt.Sprintf("%s-network", container.Name)

	s.actions.On("CheckIfImageExists", ctx, container.Image).Return(true, nil)
	s.actions.On("CreateNetwork", ctx, networkName).Return(networkID, nil)
	s.actions.On("CreateContainerWithNetwork", ctx, container.Image, container.Name, networkID, container.EnvironmentVars).Return("", errors.New("error"))

	s.logger.On("Warn", "Image already exists skipping\n").Return().Once()
	s.logger.On("Info", "Successfully created network %s \n", networkName).Return().Once()

	// Act

	err := s.sut.ServiceProvisioning(ctx, container)

	// Assert
	s.Error(err)
}

func (s *clientTestSuite) TestServiceProvisioning_WhenErrorOccursOnStartingContainer_ThenFailure() {
	// Arrange
	ctx := context.Background()
	container := docker.Container{Name: "name", Image: "image"}
	networkID := "networkID"
	containerID := "containerID"
	networkName := fmt.Sprintf("%s-network", container.Name)

	s.actions.On("CheckIfImageExists", ctx, container.Image).Return(true, nil)
	s.actions.On("CreateNetwork", ctx, networkName).Return(networkID, nil)
	s.actions.On("CreateContainerWithNetwork", ctx, container.Image, container.Name, networkID, container.EnvironmentVars).Return(containerID, nil)
	s.actions.On("StartContainer", ctx, "containerID").Return(errors.New("error"))

	s.logger.On("Warn", "Image already exists skipping\n").Return().Once()
	s.logger.On("Info", "Successfully created network %s \n", networkName).Return().Once()
	s.logger.On("Info", "Successfully created container %s\n", container.Name).Return().Once()

	// Act
	err := s.sut.ServiceProvisioning(ctx, container)

	// Assert
	s.Error(err)
}

func (s *clientTestSuite) TestServiceDecommissioning_WhenContainerExists_ThenSuccess() {
	// Arrange
	ctx := context.Background()
	container := docker.Container{Name: "name", Image: "image"}
	containerID := "containerID"

	s.actions.On("StopContainer", ctx, container.Name).Return(containerID, nil)
	s.actions.On("RemoveContainer", ctx, containerID).Return(nil)
	s.actions.On("RemoveNetwork", ctx, container.Name).Return(nil)

	s.logger.On("Info", "Successfully stopped container %s\n", container.Name).Return().Once()
	s.logger.On("Info", "Successfully removed container %s\n", container.Name).Return().Once()
	s.logger.On("Info", "Successfully removed network %s\n", container.Name).Return().Once()

	// Act
	err := s.sut.ServiceDecommissioning(ctx, container)

	// Assert
	s.NoError(err)
}

func (s *clientTestSuite) TestServiceDecommissioning_WhenDoesNotExists_ThenSuccess() {
	// Arrange
	ctx := context.Background()
	container := docker.Container{Name: "name", Image: "image"}

	s.actions.On("StopContainer", ctx, container.Name).Return("", nil)

	s.logger.On("Warn", "Container not found skipping\n").Return().Once()

	// Act
	err := s.sut.ServiceDecommissioning(ctx, container)

	// Assert
	s.NoError(err)
}

func (s *clientTestSuite) TestServiceDecommissioning_WhenErrorOccursOnStoppingContainer_ThenFailure() {
	// Arrange
	ctx := context.Background()
	container := docker.Container{Name: "name", Image: "image"}

	s.actions.On("StopContainer", ctx, container.Name).Return("", errors.New("error"))

	// Act
	err := s.sut.ServiceDecommissioning(ctx, container)

	// Assert
	s.Error(err)
}

func (s *clientTestSuite) TestServiceDecommissioning_WhenErrorOccursOnRemovingContainer_ThenFailure() {
	// Arrange
	ctx := context.Background()
	container := docker.Container{Name: "name", Image: "image"}
	containerID := "containerID"

	s.actions.On("StopContainer", ctx, container.Name).Return(containerID, nil)
	s.actions.On("RemoveContainer", ctx, containerID).Return(errors.New("error"))

	s.logger.On("Info", "Successfully stopped container %s\n", container.Name).Return().Once()

	// Act
	err := s.sut.ServiceDecommissioning(ctx, container)

	// Assert
	s.Error(err)
}

func (s *clientTestSuite) TestServiceDecommissioning_WhenErrorOccursOnRemovingNetwork_ThenFailure() {
	// Arrange
	ctx := context.Background()
	container := docker.Container{Name: "name", Image: "image"}
	containerID := "containerID"

	s.actions.On("StopContainer", ctx, container.Name).Return(containerID, nil)
	s.actions.On("RemoveContainer", ctx, containerID).Return(nil)
	s.actions.On("RemoveNetwork", ctx, container.Name).Return(errors.New("error"))

	s.logger.On("Info", "Successfully stopped container %s\n", container.Name).Return().Once()
	s.logger.On("Info", "Successfully removed container %s\n", container.Name).Return().Once()

	// Act
	err := s.sut.ServiceDecommissioning(ctx, container)

	// Assert
	s.Error(err)
}
