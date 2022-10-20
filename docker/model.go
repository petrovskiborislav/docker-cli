package docker

// Container represents a docker container.
type Container struct {
	Name            string
	Image           string
	EnvironmentVars []string
}
