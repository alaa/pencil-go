package main

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
)

const (
	Endpoint = "unix:///var/run/docker.sock"
)

type ContainerInfo struct {
	ID    string
	Name  string
	Image string
}

type Docker struct {
	client *docker.Client
}

func NewDocker() {
	client, _ := docker.NewClient(Endpoint)
	return Docker{client: client}
}

func main() {
	containers := getRunningContainers()
	for _, container := range containers {
		fmt.Println(container.ID)
		fmt.Println(container.Name)
	}
}

func (Docker *d) buildContainerInfo(container *docker.Container) ContainerInfo {
	return ContainerInfo{
		ID:    container.ID,
		Name:  container.Name,
		Image: container.Image,
	}
}

func (Docker *d) getRunningContainers() []ContainerInfo {
	client, _ := docker.NewClient(Endpoint)

	options := docker.ListContainersOptions{}
	containers, _ := client.ListContainers(options)

	containersInfo := []ContainerInfo{}

	for _, container := range containers {
		containerData, _ := client.InspectContainer(container.ID)

		containerInfo := buildContainerInfo(containerData)
		containersInfo = append(containersInfo, containerInfo)
	}
	return containersInfo
}
