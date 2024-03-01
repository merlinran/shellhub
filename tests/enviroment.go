package main

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/testcontainers/testcontainers-go"
	dockercompose "github.com/testcontainers/testcontainers-go/modules/compose"
)

const (
	serviceGateway = "gateway"
	serviceAgent   = "agent"
	serviceAPI     = "api"
	serviceCLI     = "cli"
	serviceSSH     = "ssh"
	serviceUI      = "ui"
)

type DockerCompose struct {
	envs     map[string]string
	services map[string]*testcontainers.DockerContainer
}

func NewDockerCompose() *DockerCompose {
	envs, _ := godotenv.Read("../.env")

	// envs["SHELLHUB_API_PORT"] = getFreePort()
	envs["SHELLHUB_HTTP_PORT"] = getFreePort()
	envs["SHELLHUB_SSH_PORT"] = getFreePort()

	return &DockerCompose{
		envs:     envs,
		services: make(map[string]*testcontainers.DockerContainer),
	}
}

func (e *DockerCompose) WithEnv(key, val string) *DockerCompose {
	e.envs[key] = val

	return e
}

func (e *DockerCompose) Clone() *DockerCompose {
	clonedEnv := DockerCompose{
		envs:     make(map[string]string),
		services: make(map[string]*testcontainers.DockerContainer),
	}

	for k, v := range e.envs {
		clonedEnv.envs[k] = v
	}

	// TODO: this does not make a deep copy of v
	for k, v := range e.services {
		clonedEnv.services[k] = v
	}

	// ensures that ports are always unique
	// clonedEnv.envs["SHELLHUB_API_PORT"] = getFreePort()
	clonedEnv.envs["SHELLHUB_HTTP_PORT"] = getFreePort()
	clonedEnv.envs["SHELLHUB_SSH_PORT"] = getFreePort()

	return &clonedEnv
}

// Start initiates the docker-compose environment, ensuring that all services are up and healthy before
// populating the service pointers (e.g., ComposeEnvironment.GetServiceAPI()). It returns a cleanup
// function, which should be invoked when the environment is no longer required, along with any potential
// errors encountered.
func (e *DockerCompose) Start() (func() error, error) {
	ctx := context.TODO()

	dockerCompose, err := dockercompose.NewDockerCompose("../docker-compose.yml", "../docker-compose.dev.yml")
	if err != nil {
		return nil, err
	}

	err = dockerCompose.WithEnv(e.envs).Up(ctx, dockercompose.Wait(true))
	if err != nil {
		return nil, err
	}

	/**
	 * The reason NewDockerCompose returns a private type that does not implement the DockerCompose interface
	 * is unclear. Consequently, we are unable to pass the dockerCompose variable as a parameter. Therefore,
	 * all actions involving it must be executed within this function.
	 */

	e.services[serviceGateway], err = dockerCompose.ServiceContainer(ctx, "gateway")
	if err != nil {
		return nil, err
	}

	e.services[serviceAgent], err = dockerCompose.ServiceContainer(ctx, "agent")
	if err != nil {
		return nil, err
	}

	e.services[serviceAPI], err = dockerCompose.ServiceContainer(ctx, "api")
	if err != nil {
		return nil, err
	}

	e.services[serviceCLI], err = dockerCompose.ServiceContainer(ctx, "cli")
	if err != nil {
		return nil, err
	}

	e.services[serviceSSH], err = dockerCompose.ServiceContainer(ctx, "ssh")
	if err != nil {
		return nil, err
	}

	e.services[serviceUI], err = dockerCompose.ServiceContainer(ctx, "ui")
	if err != nil {
		return nil, err
	}

	cleanup := func() error {
		err := dockerCompose.Down(context.Background(), dockercompose.RemoveOrphans(true), dockercompose.RemoveImagesLocal)
		// Clear the service pointers to prevent potential errors when accessing these services after the
		// cleanup process.
		for k := range e.services {
			e.services[k] = nil
		}

		return err
	}

	return cleanup, nil
}

// GetEnv retrieves a environment variable with the specified key.
func (e *DockerCompose) GetEnv(key string) string {
	return e.envs[key]
}

// GetServiceGateway retrieves the gateway service.
func (e *DockerCompose) GetServiceGateway() *testcontainers.DockerContainer {
	return e.services[serviceGateway]
}

// GetServiceAgent retrieves the agent service.
func (e *DockerCompose) GetServiceAgent() *testcontainers.DockerContainer {
	return e.services[serviceAgent]
}

// GetServiceAPI retrieves the api service.
func (e *DockerCompose) GetServiceAPI() *testcontainers.DockerContainer {
	return e.services[serviceAPI]
}

// GetServiceCLI retrieves the cli service.
func (e *DockerCompose) GetServiceCLI() *testcontainers.DockerContainer {
	return e.services[serviceCLI]
}

// GetServiceSSH retrieves the ssh service.
func (e *DockerCompose) GetServiceSSH() *testcontainers.DockerContainer {
	return e.services[serviceSSH]
}

// GetServiceUI retrieves the ui service.
func (e *DockerCompose) GetServiceUI() *testcontainers.DockerContainer {
	return e.services[serviceUI]
}
