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

//go:generate mockery --name=Logger --structname mockLogger --filename mock_logger_test.go --outpkg=command_test --output=. --srcpkg=github.com/petrovskiborislav/docker-cli/logger
//go:generate mockery --name=Client --structname mockClient --filename mock_client_test.go --outpkg=command_test --output=. --srcpkg=github.com/petrovskiborislav/docker-cli/docker
//go:generate mockery --name=Prompt --structname mockPrompt --filename mock_prompt_test.go --outpkg=command_test --output=. --srcpkg=github.com/petrovskiborislav/docker-cli/prompt

const filePath = "../default-compose.yaml"

type startTestSuite struct {
	suite.Suite
	logger *mockLogger
	client *mockClient
	prompt *mockPrompt
	sut    *cobra.Command
}

func (s *startTestSuite) SetupTest() {
	s.client = &mockClient{}
	s.prompt = &mockPrompt{}
	s.logger = &mockLogger{}
	s.sut = command.NewStartCommand(context.Background(), s.logger, s.prompt, s.client)
}

func TestSuite_Start(t *testing.T) {
	suite.Run(t, &startTestSuite{})
}

func (s *startTestSuite) TestStart_WhenAllServicesSelected_ThenSuccess() {
	// Arrange
	ctx := context.Background()

	msg := "Select services to start"
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

	s.client.On("ServiceProvisioning", ctx, serviceContainer1).Return(nil).Once()
	s.client.On("ServiceProvisioning", ctx, serviceContainer2).Return(nil).Once()
	s.client.On("ServiceProvisioning", ctx, serviceContainer3).Return(nil).Once()
	s.client.On("ServiceProvisioning", ctx, serviceContainer4).Return(nil).Once()

	// Act
	s.sut.Run(nil, []string{filePath})

	// Assert
	s.logger.AssertExpectations(s.T())
	s.prompt.AssertExpectations(s.T())
	s.client.AssertExpectations(s.T())
}

func (s *startTestSuite) TestStart_WhenSingleServiceSelected_ThenSuccess() {
	// Arrange
	ctx := context.Background()

	msg := "Select services to start"
	items := []string{"all", "nginx", "db", "cache", "wordpress"}
	serviceContainer := docker.Container{Name: "nginx", Image: "nginx:alpine"}

	matcher := mock.MatchedBy(matchElements(items))
	s.prompt.On("SelectPrompt", msg, matcher).Return(items[1:2], nil)

	s.client.On("ServiceProvisioning", ctx, serviceContainer).Return(nil)

	// Act
	s.sut.Run(nil, []string{filePath})

	// Assert
	s.logger.AssertExpectations(s.T())
	s.prompt.AssertExpectations(s.T())
	s.client.AssertExpectations(s.T())
}

func (s *startTestSuite) TestStart_WhenErrorOccursOnParsing_ThenFailure() {
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

func (s *startTestSuite) TestStart_WhenErrorOccursOnSelectingServices_ThenFailure() {
	// Arrange
	msg := "Select services to start"
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

func (s *startTestSuite) TestStart_WhenErrorOccursOnServiceProvisioning_ThenFailure() {
	// Arrange
	ctx := context.Background()

	msg := "Select services to start"
	items := []string{"all", "nginx", "db", "cache", "wordpress"}
	serviceContainer := docker.Container{Name: "nginx", Image: "nginx:alpine"}

	matcher := mock.MatchedBy(matchElements(items))
	s.prompt.On("SelectPrompt", msg, matcher).Return(items[1:2], nil)

	s.client.On("ServiceProvisioning", ctx, serviceContainer).Return(errors.New("error"))

	err := errors.New("error")
	s.logger.On("Error", "Error starting services: %s\n", err).Return()

	// Act
	s.sut.Run(nil, []string{})

	// Assert
	s.logger.AssertExpectations(s.T())
	s.prompt.AssertExpectations(s.T())
	s.client.AssertExpectations(s.T())
}

func matchElements(x []string) func(y []string) bool {
	return func(y []string) bool {
		if len(x) != len(y) {
			return false
		}
		diff := make(map[string]int, len(x))
		for _, _x := range x {
			diff[_x]++
		}
		for _, _y := range y {
			if _, ok := diff[_y]; !ok {
				return false
			}
			diff[_y] -= 1
			if diff[_y] == 0 {
				delete(diff, _y)
			}
		}
		return len(diff) == 0
	}
}
