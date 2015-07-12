package docker

import (
	docker "github.com/fsouza/go-dockerclient"
	"log"
	"time"
)

const (
	Endpoint = "unix:///var/run/docker.sock"
	Interval = 5 * time.Second
)

type DockerClient struct {
	client *docker.Client
}

func NewDocker() DockerClient {
	client, err := docker.NewClient(Endpoint)
	if err != nil {
		log.Fatal(err)
	}
	return DockerClient{client: client}
}

type ContainerInfo struct {
	ID      string
	Name    string
	Image   string
	Created time.Time
	Config  *docker.Config
}

func buildContainerInfo(container *docker.Container) ContainerInfo {
	return ContainerInfo{
		ID:      container.ID,
		Name:    container.Name,
		Image:   container.Image,
		Created: container.Created,
		Config:  container.Config,
	}
}

// getRunningContainers finds running containers and returns specific details.
func (c *DockerClient) getRunningContainers() []ContainerInfo {
	containersIDs, err := c.getContainersIDs()
	if err != nil {
		log.Print(err)
	}
	containersDetails := c.getContainersDetails(containersIDs)
	log.Print("Running containers count: ", len(containersDetails), "\n\n")
	return containersDetails
}

// getContainersIDs retruns a list of running docker contianers.
func (c *DockerClient) getContainersIDs() ([]docker.APIContainers, error) {
	options := docker.ListContainersOptions{}
	containers, err := c.client.ListContainers(options)
	if err != nil {
		return containers, err
	}
	return containers, nil
}

// getContainersDetails iterate over a list of containers and returns a list of ContainerInfo struct.
func (c *DockerClient) getContainersDetails(containers []docker.APIContainers) []ContainerInfo {
	list := []ContainerInfo{}
	for _, container := range containers {
		i, err := c.inspectContainer(container.ID)
		if err != nil {
			log.Print(err)
		}
		list = append(list, i)
	}
	return list
}

// inspectContainer extract container info for a continer ID.
func (c *DockerClient) inspectContainer(cid string) (ContainerInfo, error) {
	data, err := c.client.InspectContainer(cid)
	if err != nil {
		return ContainerInfo{}, err
	}
	return buildContainerInfo(data), nil
}

func main() {
	client := NewDocker()
	c := time.Tick(Interval)
	for now := range c {
		containers := client.getRunningContainers()
		log.Print(now, containers)
	}
}
