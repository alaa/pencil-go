package main

import (
	"fmt"
	"github.com/brainly/pencil-go/consul"
	"github.com/brainly/pencil-go/docker"
	"github.com/brainly/pencil-go/registry"
	dockerclient "github.com/fsouza/go-dockerclient"
	consulclient "github.com/hashicorp/consul/api"
	"time"
)

func main() {
	fmt.Println("starting pencil ...\n")
	registry := registry.NewRegistry(getContainerRepository(), getServiceRepository())
	for range time.Tick(5 * time.Second) {
		registry.Synchronize()
	}
}

func getContainerRepository() registry.ContainerRepository {
	client, _ := dockerclient.NewClientFromEnv()
	return docker.NewContainerRepository(client)
}

func getServiceRepository() registry.ServiceRepository {
	consulClient, _ := consulclient.NewClient(consulclient.DefaultConfig())
	return consul.NewServiceRepository(consulClient.Agent())
}
