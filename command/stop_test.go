package command_test

import (
	"context"
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/petrovskiborislav/docker-cli/command"
	"github.com/petrovskiborislav/docker-cli/docker"
)

type stopTestSuite struct {
	suite.Suite
	logger *mockLogger
	client *mockClient
	prompt *mockPrompt
	sut    *cobra.Command
}

func (s *stopTestSuite) SetupTest() {
	s.client = &mockClient{}
	s.prompt = &mockPrompt{}
	s.logger = &mockLogger{}
	s.sut = command.NewStopCommand(context.Background(), s.logger, s.prompt, s.client)
}

func TestSuite_Stop(t *testing.T) {
	suite.Run(t, &stopTestSuite{})
}

func (s *stopTestSuite) TestStop_WhenAllServicesSelected_ThenSuccess() {
	// Arrange
	ctx := context.Background()

	msg := "Select services to stop"
	items := []string{"all", "nginx", "db", "cache", "wordpress"}
	serviceContainer1 := docker.Container{Name: "nginx", Image: "nginx:alpine"}
	serviceContainer2 := docker.Container{
		Name:            "db",
		Image:           "mysql:latest",
		EnvironmentVars: []string{"MYSQL_ALLOW_EMPTY_PASSWORD=true"}}
	serviceContainer3 := docker.Container{Name: "cache", Image: "memcached"}
	serviceContainer4 := docker.Container{Name: "wordpress", Image: "wordpress:6.0"}

	matcher := mock.MatchedBy(matchElements(items))
	s.prompt.On("SelectPrompt", msg, matcher).Return(items[0:1], nil)

	s.client.On("ServiceDecommissioning", ctx, serviceContainer1).Return(nil).Once()
	s.client.On("ServiceDecommissioning", ctx, serviceContainer2).Return(nil).Once()
	s.client.On("ServiceDecommissioning", ctx, serviceContainer3).Return(nil).Once()
	s.client.On("ServiceDecommissioning", ctx, serviceContainer4).Return(nil).Once()

	// Act
	s.sut.Run(nil, []string{filePath})

	// Assert
	s.logger.AssertExpectations(s.T())
	s.prompt.AssertExpectations(s.T())
	s.client.AssertExpectations(s.T())
}

func (s *stopTestSuite) TestStop_WhenSingleServiceSelected_ThenSuccess() {
	// Arrange
	ctx := context.Background()

	msg := "Select services to stop"
	items := []string{"all", "nginx", "db", "cache", "wordpress"}
	serviceContainer := docker.Container{Name: "nginx", Image: "nginx:alpine"}

	matcher := mock.MatchedBy(matchElements(items))
	s.prompt.On("SelectPrompt", msg, matcher).Return(items[1:2], nil)

	s.client.On("ServiceDecommissioning", ctx, serviceContainer).Return(nil)

	// Act
	s.sut.Run(nil, []string{filePath})

	// Assert
	s.logger.AssertExpectations(s.T())
	s.prompt.AssertExpectations(s.T())
	s.client.AssertExpectations(s.T())
}

func (s *stopTestSuite) TestStop_WhenErrorOccursOnParsing_ThenFailure() {
	// Arrange
	err := errors.New("error reading YAML file: open : no such file or directory")
	s.logger.On("Error", "Error parsing compose file: %s\n", err).Return()

	// Act
	s.sut.Run(nil, []string{""})

	// Assert
	s.logger.AssertExpectations(s.T())
	s.prompt.AssertExpectations(s.T())
	s.client.AssertExpectations(s.T())
}

func (s *stopTestSuite) TestStop_WhenErrorOccursOnSelectingServices_ThenFailure() {
	// Arrange
	msg := "Select services to stop"
	items := []string{"all", "nginx", "db", "cache", "wordpress"}

	matcher := mock.MatchedBy(matchElements(items))
	s.prompt.On("SelectPrompt", msg, matcher).Return(nil, errors.New("error"))

	err := errors.New("error")
	s.logger.On("Error", "Error selecting services: %s\n", err).Return()

	// Act
	s.sut.Run(nil, []string{})

	// Assert
	s.logger.AssertExpectations(s.T())
	s.prompt.AssertExpectations(s.T())
	s.client.AssertExpectations(s.T())
}

func (s *stopTestSuite) TestStop_WhenErrorOccursOnServiceProvisioning_ThenFailure() {
	// Arrange
	ctx := context.Background()

	msg := "Select services to stop"
	items := []string{"all", "nginx", "db", "cache", "wordpress"}
	serviceContainer := docker.Container{Name: "nginx", Image: "nginx:alpine"}

	matcher := mock.MatchedBy(matchElements(items))
	s.prompt.On("SelectPrompt", msg, matcher).Return(items[1:2], nil)

	s.client.On("ServiceDecommissioning", ctx, serviceContainer).Return(errors.New("error"))

	err := errors.New("error")
	s.logger.On("Error", "Error stopping services: %s\n", err).Return()

	// Act
	s.sut.Run(nil, []string{})

	// Assert
	s.logger.AssertExpectations(s.T())
	s.prompt.AssertExpectations(s.T())
	s.client.AssertExpectations(s.T())
}
