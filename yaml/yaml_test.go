package yaml_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/petrovskiborislav/docker-cli/yaml"
)

const (
	path = "../default-compose.yaml"
)

func TestParseComposeFile_ThenSuccess(t *testing.T) {
	// Arrange

	// Act
	result, err := yaml.ParseComposeFile(path)

	// Assert
	want := map[string]yaml.Service{
		"nginx": {Image: "nginx:alpine"},
		"db": {
			Image:           "mysql:latest",
			EnvironmentVars: map[string]string{"MYSQL_ALLOW_EMPTY_PASSWORD": "true"},
		},
		"cache":     {Image: "memcached"},
		"wordpress": {Image: "wordpress:6.0"},
	}

	assert.NoError(t, err)
	assert.EqualValues(t, want, result)
}

func TestParseComposeFile_ThenFailure(t *testing.T) {
	// Arrange

	// Act
	result, err := yaml.ParseComposeFile("invalid-path")

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
}
