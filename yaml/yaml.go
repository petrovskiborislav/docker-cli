package yaml

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// Services is a struct which represents the composer YAML file.
type Services struct {
	Services map[string]Service `yaml:"services"`
}

// Service is a struct which represents a service in a composer YAML file.
type Service struct {
	Image           string            `yaml:"image"`
	EnvironmentVars map[string]string `yaml:"environment"`
}

// ParseComposeFile parses a composer YAML file and returns a map of services.
func ParseComposeFile(path string) (map[string]Service, error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %s", err)
	}

	yamlServices := Services{}
	err = yaml.Unmarshal(yamlFile, &yamlServices)
	if err != nil {
		return nil, err
	}

	return yamlServices.Services, nil
}
