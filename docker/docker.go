package docker

import (
	"github.com/alaa/pencil-go/registry"
	docker "github.com/fsouza/go-dockerclient"
)

type dockerClient interface {
	ListContainers(opts docker.ListContainersOptions) ([]docker.APIContainers, error)
	InspectContainer(id string) (*docker.Container, error)
}

// ContainerRepository is docker-based implementation of registry.ContainerRepository
type ContainerRepository struct {
	dockerClient dockerClient
}

// NewContainerRepository creates new instance of ContainerRepository structure
func NewContainerRepository(dockerClient dockerClient) *ContainerRepository {
	return &ContainerRepository{dockerClient: dockerClient}
}

// GetAll returns list of all running docker containers
func (cr *ContainerRepository) GetAll() ([]registry.Container, error) {
	containersIDs, err := cr.getContainersIDs()
	if err != nil {
		return nil, err
	}
	containers, err := cr.getContainers(containersIDs)
	if err != nil {
		return nil, err
	}
	return containers, nil
}

func (cr *ContainerRepository) getContainersIDs() ([]string, error) {
	containersIDs := []string{}
	containers, err := cr.dockerClient.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		return nil, err
	}
	for _, container := range containers {
		containersIDs = append(containersIDs, container.ID)
	}
	return containersIDs, nil
}

func (cr *ContainerRepository) getContainers(containersIDs []string) ([]registry.Container, error) {
	containers := []registry.Container{}
	for _, containerID := range containersIDs {
		containerDetails, err := cr.dockerClient.InspectContainer(containerID)
		if err != nil {
			return nil, err
		}
		containers = append(containers, buildContainers(containerDetails)...)
	}
	return containers, nil
}

func buildContainers(container *docker.Container) []registry.Container {
	containerWrapper := dockerContainerWrapper{*container}
	containers := []registry.Container{}

	for _, port := range containerWrapper.getExposedTCPPorts() {
		container := registry.Container{
			ID:   containerWrapper.ID,
			Name: containerWrapper.getName(),
			Tags: containerWrapper.getTags(),
			Port: port,
		}
		containers = append(containers, container)
	}
	return containers
}
