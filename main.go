package main

import (
	"fmt"
	//consul "github.com/alaa/pencil-go/consul"
	docker "github.com/alaa/pencil-go/docker"
	"time"
)

func main() {
	fmt.Println("starting pencil ...\n")
	client := docker.NewDocker()
	c := time.Tick(docker.Interval)
	for range c {
		containers := client.GetRunningContainers()
		for _, container := range containers {
			fmt.Printf("%s %s \n\n", container.ID, container.TCPPorts)
		}
	}
}
