package main

import (
	"fmt"
	//consul "github.com/alaa/pencil-go/consul"
	"github.com/alaa/pencil-go/docker"
	dockerclient "github.com/fsouza/go-dockerclient"
	"time"
)

func main() {
	fmt.Println("starting pencil ...\n")
	client, _ := dockerclient.NewClient(docker.Endpoint)
	adapter := docker.NewDockerAdapter(client)
	c := time.Tick(docker.Interval)
	for range c {
		containers, _ := adapter.GetRunningContainers()
		fmt.Printf("containers: %v\n", containers)
		for _, container := range containers {
			fmt.Printf("%s %s %s service_name:%s \n\n", container.ID, container.ImageName, container.TCPPorts, container.ServiceName)
		}
	}
}
