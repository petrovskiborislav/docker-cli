package command

import (
	"github.com/spf13/cobra"

	"github.com/petrovskiborislav/docker-cli/docker"
	"github.com/petrovskiborislav/docker-cli/prompt"
	"github.com/petrovskiborislav/docker-cli/yaml"
)

const (
	allPromptOption        = "all"
	defaultComposeFilePath = "../default-compose.yaml"
)

// NewRootCommand creates the base command when called without any subcommands.
func NewRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "docker-cli [OPTIONS]",
		Short: "CLI for docker",
	}
}

func parseComposeFile(args []string) (map[string]yaml.Service, error) {
	filePath := defaultComposeFilePath
	if len(args) > 0 {
		filePath = args[0]
	}

	return yaml.ParseComposeFile(filePath)
}

func selectServiceContainers(label string, prompt prompt.Prompt, services map[string]yaml.Service) ([]docker.Container, error) {
	promptOptions := []string{allPromptOption}
	for name, _ := range services {
		promptOptions = append(promptOptions, name)
	}

	selectedServices, err := prompt.SelectPrompt(label, promptOptions)
	if err != nil {
		return nil, err
	}

	return selectedServicesToContainers(selectedServices, services), nil
}

func selectedServicesToContainers(selectedServices []string, services map[string]yaml.Service) []docker.Container {
	var containers []docker.Container
	for _, serviceName := range selectedServices {
		if serviceName == allPromptOption {
			for name, service := range services {
				containers = append(containers, newDockerContainer(name, service))
			}
			break
		}

		if val, ok := services[serviceName]; ok {
			containers = append(containers, newDockerContainer(serviceName, val))
		}
	}

	return containers
}

func newDockerContainer(name string, service yaml.Service) docker.Container {
	var envs []string
	for key, value := range service.EnvironmentVars {
		envs = append(envs, key+"="+value)
	}

	return docker.Container{Name: name, Image: service.Image, EnvironmentVars: envs}
}
