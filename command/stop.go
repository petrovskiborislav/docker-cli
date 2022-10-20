package command

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/petrovskiborislav/docker-cli/docker"
	"github.com/petrovskiborislav/docker-cli/logger"
	"github.com/petrovskiborislav/docker-cli/prompt"
)

// NewStopCommand creates stop command which reads
// compose file and stops the selected services.
func NewStopCommand(ctx context.Context, logger logger.Logger, prompt prompt.Prompt, client docker.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "stop [PATH to docker-compose file]",
		Short: "Stops the selected services listed from the specified compose file",
		Run: func(cmd *cobra.Command, args []string) {
			parsedServices, err := parseComposeFile(args)
			if err != nil {
				logger.Error("Error parsing compose file: %s\n", err)
				return
			}

			selectedServiceContainers, err := selectServiceContainers("Select services to stop", prompt, parsedServices)
			if err != nil {
				logger.Error("Error selecting services: %s\n", err)
				return
			}

			for _, serviceContainer := range selectedServiceContainers {
				if err := client.ServiceDecommissioning(ctx, serviceContainer); err != nil {
					logger.Error("Error stopping services: %s\n", err)
					return
				}
			}
		},
	}
}
