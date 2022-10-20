package docker_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/petrovskiborislav/docker-cli/docker"
)

//go:generate mockery --name=APIClient --structname mockAPIClient --filename mock_api_clinet_test.go --outpkg=docker_test --output=. --srcpkg=github.com/docker/docker/client

type actionsTestSuite struct {
	suite.Suite
	client *mockAPIClient
	sut    docker.Actions
}

func (s *actionsTestSuite) SetupTest() {
	s.client = &mockAPIClient{}
	s.sut = docker.NewActions(s.client)
}

func (s *actionsTestSuite) AfterTest(suiteName string, testName string) {
	s.client.AssertExpectations(s.T())
}

func TestSuite_Actions(t *testing.T) {
	suite.Run(t, &actionsTestSuite{})
}

func (s *actionsTestSuite) TestCheckIfImageExists_WhenImageExists_ThenSuccess() {
	// Arrange
	ctx := context.Background()
	image := "image"

	filter := filters.NewArgs()
	filter.Add("reference", image)
	imageListOptions := types.ImageListOptions{Filters: filter}

	s.client.On("ImageList", ctx, imageListOptions).Return([]types.ImageSummary{{}}, nil)

	// Act
	exists, err := s.sut.CheckIfImageExists(ctx, image)

	// Assert
	s.NoError(err)
	s.EqualValues(true, exists)
}

func (s *actionsTestSuite) TestCheckIfImageExists_WhenImageDoesNotExists_ThenSuccess() {
	// Arrange
	ctx := context.Background()
	image := "image"

	filter := filters.NewArgs()
	filter.Add("reference", image)
	imageListOptions := types.ImageListOptions{Filters: filter}

	s.client.On("ImageList", ctx, imageListOptions).Return([]types.ImageSummary{}, nil)

	// Act
	exists, err := s.sut.CheckIfImageExists(ctx, image)

	// Assert
	s.NoError(err)
	s.EqualValues(false, exists)
}

func (s *actionsTestSuite) TestCheckIfImageExists_WhenErrorOccursOnImageList_ThenFailure() {
	// Arrange
	ctx := context.Background()
	image := "image"

	filter := filters.NewArgs()
	filter.Add("reference", image)
	imageListOptions := types.ImageListOptions{Filters: filter}

	s.client.On("ImageList", ctx, imageListOptions).Return(nil, errors.New("error"))

	// Act
	exists, err := s.sut.CheckIfImageExists(ctx, image)

	// Assert
	s.Error(err)
	s.EqualValues(false, exists)
}

func (s *actionsTestSuite) TestPullImage_ThenSuccess() {
	// Arrange
	ctx := context.Background()
	image := "image"
	reader := io.NopCloser(strings.NewReader("Success"))

	s.client.On("ImagePull", ctx, image, types.ImagePullOptions{}).Return(reader, nil)

	// Act
	err := s.sut.PullImage(ctx, image)

	// Assert
	s.NoError(err)
}

func (s *actionsTestSuite) TestPullImage_ThenFailure() {
	// Arrange
	ctx := context.Background()
	image := "image"

	s.client.On("ImagePull", ctx, image, types.ImagePullOptions{}).Return(nil, errors.New("error"))

	// Act
	err := s.sut.PullImage(ctx, image)

	// Assert
	s.Error(err)
}

func (s *actionsTestSuite) TestCreateNetwork_ThenSuccess() {
	// Arrange
	ctx := context.Background()
	networkName := "network"
	networkID := "id"

	networkResp := types.NetworkCreateResponse{ID: "id"}
	s.client.On("NetworkCreate", ctx, networkName, types.NetworkCreate{}).Return(networkResp, nil)

	// Act
	id, err := s.sut.CreateNetwork(ctx, networkName)

	// Assert
	s.NoError(err)
	s.Equal(networkID, id)
}

func (s *actionsTestSuite) TestCreateNetwork_ThenFailure() {
	// Arrange
	ctx := context.Background()
	networkName := "network"

	s.client.On("NetworkCreate", ctx, networkName, types.NetworkCreate{}).Return(types.NetworkCreateResponse{}, errors.New("error"))

	// Act
	id, err := s.sut.CreateNetwork(ctx, networkName)

	// Assert
	s.Error(err)
	s.Equal("", id)
}

func (s *actionsTestSuite) TestCreateContainerWithNetwork_ThenSuccess() {
	// Arrange
	ctx := context.Background()
	imageName := "image"
	containerName := "container"
	networkID := "id"
	containerID := "id"

	containerConfig := &container.Config{Image: imageName}
	containerCreateCreatedBody := container.ContainerCreateCreatedBody{ID: containerID}
	s.client.On("ContainerCreate", ctx, containerConfig, mock.Anything, mock.Anything, mock.Anything, containerName).Return(containerCreateCreatedBody, nil)
	s.client.On("NetworkConnect", ctx, networkID, containerID, mock.Anything).Return(nil)

	// Act
	id, err := s.sut.CreateContainerWithNetwork(ctx, imageName, containerName, networkID, nil)

	// Assert
	s.NoError(err)
	s.Equal(containerID, id)
}

func (s *actionsTestSuite) TestCreateContainerWithNetwork_WhenErrorOccursOnContainerCreation_ThenFailure() {
	// Arrange
	ctx := context.Background()
	imageName := "image"
	containerName := "container"
	networkID := "id"

	containerConfig := &container.Config{Image: imageName}
	containerCreateCreatedBody := container.ContainerCreateCreatedBody{}
	s.client.On("ContainerCreate", ctx, containerConfig, mock.Anything, mock.Anything, mock.Anything, containerName).Return(containerCreateCreatedBody, errors.New("error"))

	// Act
	id, err := s.sut.CreateContainerWithNetwork(ctx, imageName, containerName, networkID, nil)

	// Assert
	s.Error(err)
	s.Equal("", id)
}

func (s *actionsTestSuite) TestCreateContainerWithNetwork_WhenErrorOccursOnNetowrkConnection_ThenFailure() {
	// Arrange
	ctx := context.Background()
	imageName := "image"
	containerName := "container"
	networkID := "id"
	containerID := "id"

	containerConfig := &container.Config{Image: imageName}
	containerCreateCreatedBody := container.ContainerCreateCreatedBody{ID: containerID}
	s.client.On("ContainerCreate", ctx, containerConfig, mock.Anything, mock.Anything, mock.Anything, containerName).Return(containerCreateCreatedBody, nil)
	s.client.On("NetworkConnect", ctx, networkID, containerID, mock.Anything).Return(errors.New("error"))

	// Act
	id, err := s.sut.CreateContainerWithNetwork(ctx, imageName, containerName, networkID, nil)

	// Assert
	s.Error(err)
	s.Equal("", id)
}

func (s *actionsTestSuite) TestStartContainer_ThenSuccess() {
	// Arrange
	ctx := context.Background()
	containerID := "id"

	s.client.On("ContainerStart", ctx, containerID, types.ContainerStartOptions{}).Return(nil)

	// Act
	err := s.sut.StartContainer(ctx, containerID)

	// Assert
	s.NoError(err)
}

func (s *actionsTestSuite) TestStartContainer_ThenFailure() {
	// Arrange
	ctx := context.Background()
	containerID := "id"

	s.client.On("ContainerStart", ctx, containerID, types.ContainerStartOptions{}).Return(errors.New("error"))

	// Act
	err := s.sut.StartContainer(ctx, containerID)

	// Assert
	s.Error(err)
}

func (s *actionsTestSuite) TestStopContainer_ThenSuccess() {
	// Arrange
	ctx := context.Background()
	containerName := "container"
	containerID := "id"

	filter := filters.NewArgs()
	filter.Add("name", containerName)
	containerListOptions := types.ContainerListOptions{Filters: filter}
	containers := []types.Container{{ID: containerID}}
	s.client.On("ContainerList", ctx, containerListOptions).Return(containers, nil)
	s.client.On("ContainerStop", ctx, containerID, mock.Anything).Return(nil)

	// Act
	id, err := s.sut.StopContainer(ctx, containerName)

	// Assert
	s.NoError(err)
	s.Equal(containerID, id)
}

func (s *actionsTestSuite) TestStopContainer_WhenErrorOccursOnContainerListing_ThenFailure() {
	// Arrange
	ctx := context.Background()
	containerName := "container"

	filter := filters.NewArgs()
	filter.Add("name", containerName)
	containerListOptions := types.ContainerListOptions{Filters: filter}
	s.client.On("ContainerList", ctx, containerListOptions).Return(nil, errors.New("error"))

	// Act
	id, err := s.sut.StopContainer(ctx, containerName)

	// Assert
	s.Error(err)
	s.Equal("", id)
}

func (s *actionsTestSuite) TestStopContainer_WhenNoContainersFound_ThenFailure() {
	// Arrange
	ctx := context.Background()
	containerName := "container"

	filter := filters.NewArgs()
	filter.Add("name", containerName)
	containerListOptions := types.ContainerListOptions{Filters: filter}
	containers := []types.Container{}
	s.client.On("ContainerList", ctx, containerListOptions).Return(containers, nil)

	// Act
	id, err := s.sut.StopContainer(ctx, containerName)

	// Assert
	s.Error(err)
	s.Equal("", id)
}

func (s *actionsTestSuite) TestStopContainer_WhenErrorOccursOnContainerStopping_ThenFailure() {
	// Arrange
	ctx := context.Background()
	containerName := "container"
	containerID := "id"

	filter := filters.NewArgs()
	filter.Add("name", containerName)
	containerListOptions := types.ContainerListOptions{Filters: filter}
	containers := []types.Container{{ID: containerID}}
	s.client.On("ContainerList", ctx, containerListOptions).Return(containers, nil)
	s.client.On("ContainerStop", ctx, containerID, mock.Anything).Return(errors.New("error"))

	// Act
	_, err := s.sut.StopContainer(ctx, containerName)

	// Assert
	s.Error(err)
}

func (s *actionsTestSuite) TestRemoveContainer_ThenSuccess() {
	// Arrange
	ctx := context.Background()
	containerID := "id"

	s.client.On("ContainerRemove", ctx, containerID, types.ContainerRemoveOptions{}).Return(nil)

	// Act
	err := s.sut.RemoveContainer(ctx, containerID)

	// Assert
	s.NoError(err)
}

func (s *actionsTestSuite) TestRemoveContainer_ThenFailure() {
	// Arrange
	ctx := context.Background()
	containerID := "id"

	s.client.On("ContainerRemove", ctx, containerID, types.ContainerRemoveOptions{}).Return(errors.New("error"))

	// Act
	err := s.sut.RemoveContainer(ctx, containerID)

	// Assert
	s.Error(err)
}

func (s *actionsTestSuite) TestRemoveNetwork_ThenSuccess() {
	// Arrange
	ctx := context.Background()
	containerName := "container"
	networkID := "id"

	filter := filters.NewArgs()
	filter.Add("name", fmt.Sprintf("%s-network", containerName))

	networkListOptions := types.NetworkListOptions{Filters: filter}
	networks := []types.NetworkResource{{ID: networkID}}
	s.client.On("NetworkList", ctx, networkListOptions).Return(networks, nil)
	s.client.On("NetworkRemove", ctx, networkID).Return(nil)

	// Act
	err := s.sut.RemoveNetwork(ctx, containerName)

	// Assert
	s.NoError(err)
}

func (s *actionsTestSuite) TestRemoveNetwork_WhenErrorOccursOnNetworkList_ThenFailure() {
	// Arrange
	ctx := context.Background()
	containerName := "container"

	filter := filters.NewArgs()
	filter.Add("name", fmt.Sprintf("%s-network", containerName))

	networkListOptions := types.NetworkListOptions{Filters: filter}
	s.client.On("NetworkList", ctx, networkListOptions).Return(nil, errors.New("error"))

	// Act
	err := s.sut.RemoveNetwork(ctx, containerName)

	// Assert
	s.Error(err)
}

func (s *actionsTestSuite) TestRemoveNetwork_WhenErrorOccursOnNetworkRemove_ThenFailure() {
	// Arrange
	ctx := context.Background()
	containerName := "container"
	networkID := "id"

	filter := filters.NewArgs()
	filter.Add("name", fmt.Sprintf("%s-network", containerName))

	networkListOptions := types.NetworkListOptions{Filters: filter}
	networks := []types.NetworkResource{{ID: networkID}}
	s.client.On("NetworkList", ctx, networkListOptions).Return(networks, nil)
	s.client.On("NetworkRemove", ctx, networkID).Return(errors.New("error"))

	// Act
	err := s.sut.RemoveNetwork(ctx, containerName)

	// Assert
	s.Error(err)
}
