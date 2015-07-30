package main

import (
	"fmt"
	"github.com/alaa/pencil-go/consul"
	"github.com/alaa/pencil-go/docker"
	"github.com/alaa/pencil-go/registry"
	dockerclient "github.com/fsouza/go-dockerclient"
	consulclient "github.com/hashicorp/consul/api"
	"log"
	"time"
)

func main() {
	fmt.Println("starting pencil ...\n")
	registry := registry.NewRegistry(getContainerRepository(), getServiceRepository())
	for range time.Tick(5 * time.Second) {
		err := registry.Synchronize()
		if err != nil {
			log.Printf("Error occured during synchronization: %v\n", err)
		}
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
