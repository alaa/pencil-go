package main

import (
	"fmt"
	"log"
	"time"
	//consul "github.com/alaa/pencil-go/consul"
	docker "github.com/alaa/pencil-go/docker"
)

func main() {
	fmt.Println("starting pencil ...\n")
	client := docker.NewDocker()
	c := time.Tick(docker.Interval)
	for now := range c {
		containers := client.GetRunningContainers()

		for _, container := range containers {
			log.Printf("%s %s %s", now, container.ID, container.Name)
		}
	}
}
