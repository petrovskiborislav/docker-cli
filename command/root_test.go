package command_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/petrovskiborislav/docker-cli/command"
)

func TestNewRootCommand(t *testing.T) {
	// Arrange
	// Act
	rootCommand := command.NewRootCommand()

	// Assert
	want := &cobra.Command{
		Use:   "docker-cli [OPTIONS]",
		Short: "CLI for docker",
	}

	assert.EqualValues(t, want, rootCommand)
}
