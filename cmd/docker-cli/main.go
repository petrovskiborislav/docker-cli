package main

import (
	"context"
	"os"

	"github.com/docker/docker/client"

	"github.com/petrovskiborislav/docker-cli/command"
	"github.com/petrovskiborislav/docker-cli/docker"
	"github.com/petrovskiborislav/docker-cli/logger"
	"github.com/petrovskiborislav/docker-cli/prompt"
)

func main() {
	ctx := context.Background()
	log := logger.NewLogger()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Error("error creating docker client: %s", err)
	}

	actions := docker.NewActions(cli)
	dockerClient := docker.NewClient(log, actions)

	pr := prompt.NewPrompt()

	startCmd := command.NewStartCommand(ctx, log, pr, dockerClient)
	stopCmd := command.NewStopCommand(ctx, log, pr, dockerClient)

	rootCmd := command.NewRootCommand()
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.AddCommand(startCmd, stopCmd)

	if err = rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
