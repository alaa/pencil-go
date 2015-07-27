package main

import (
	"fmt"
	"github.com/alaa/pencil-go/consul"
	"github.com/alaa/pencil-go/docker"
	dockerclient "github.com/fsouza/go-dockerclient"
	consulclient "github.com/hashicorp/consul/api"
	"time"
)

func main() {
	fmt.Println("starting pencil ...\n")

	dockerAdapter := getDockerAdapter()

	c := time.Tick(docker.Interval)
	for range c {
		containers, _ := dockerAdapter.GetRunningContainers()
		fmt.Printf("containers: %v\n", containers)
		for _, container := range containers {
			fmt.Printf("%s %s %s service_name:%s \n\n", container.ID, container.ImageName, container.TCPPorts, container.ServiceName)
		}
	}
}

func getDockerAdapter() *docker.DockerAdapter {
	client, _ := dockerclient.NewClient(docker.Endpoint)
	return docker.NewDockerAdapter(client)
}

func getConsulAdapter() *consul.ServiceRepository {
	consulClient, _ := consulclient.NewClient(consulclient.DefaultConfig())
	return consul.NewServiceRepository(consulClient.Agent())
}
