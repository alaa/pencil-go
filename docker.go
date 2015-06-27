package main

import (
	"github.com/fsouza/go-dockerclient"
	"log"
	"time"
)

const (
	Endpoint = "unix:///var/run/docker.sock"
)

type DockerClient struct {
	client *docker.Client
}

func NewDocker() DockerClient {
	client, _ := docker.NewClient(Endpoint)
	return DockerClient{client: client}
}

type ContainerInfo struct {
	ID      string
	Name    string
	Image   string
	Created time.Time
	Config  *docker.Config
}

func (d *DockerClient) buildContainerInfo(container *docker.Container) ContainerInfo {
	return ContainerInfo{
		ID:      container.ID,
		Name:    container.Name,
		Image:   container.Image,
		Created: container.Created,
		Config:  container.Config,
	}
}

func main() {
	log.Print("Starting Pencil ... \n\n")
	client := NewDocker()
	s := client.getRunningContainers()
	log.Print(s)
}

// getRunningContainers finds running containers and returns specific details.
func (d *DockerClient) getRunningContainers() []ContainerInfo {
	running_containers_ids := d.getContainersIDs()
	a := d.getContainersDetails(running_containers_ids)
	return a
}

// getContainersIDs retruns a list of running docker contianers.
func (d *DockerClient) getContainersIDs() []docker.APIContainers {
	options := docker.ListContainersOptions{}
	containers, _ := d.client.ListContainers(options)
	return containers
}

// getContainersDetails iterate over a list of containers and returns a list of ContainerInfo struct.
func (d *DockerClient) getContainersDetails(containers []docker.APIContainers) []ContainerInfo {
	list := []ContainerInfo{}
	for _, c := range containers {
		list = append(list, d.inspectContainer(c.ID))
	}
	return list
}

// inspectContainer extract container info for a continer ID.
func (d *DockerClient) inspectContainer(cid string) ContainerInfo {
	data, _ := d.client.InspectContainer(cid)
	return d.buildContainerInfo(data)
}
