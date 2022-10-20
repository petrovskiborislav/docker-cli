package command

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/petrovskiborislav/docker-cli/docker"
	"github.com/petrovskiborislav/docker-cli/logger"
	"github.com/petrovskiborislav/docker-cli/prompt"
)

// NewStartCommand creates start command which reads
// compose file and starts the selected services.
func NewStartCommand(ctx context.Context, logger logger.Logger, prompt prompt.Prompt, client docker.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "start [PATH to docker-compose file]",
		Short: "Starts the selected services listed from the specified compose file",
		Run: func(cmd *cobra.Command, args []string) {
			parsedServices, err := parseComposeFile(args)
			if err != nil {
				logger.Error("Error parsing compose file: %s\n", err)
				return
			}

			selectedServiceContainers, err := selectServiceContainers("Select services to start", prompt, parsedServices)
			if err != nil {
				logger.Error("Error selecting services: %s\n", err)
				return
			}

			for _, serviceContainer := range selectedServiceContainers {
				if err := client.ServiceProvisioning(ctx, serviceContainer); err != nil {
					logger.Error("Error starting services: %s\n", err)
					return
				}
			}
		},
	}
}
