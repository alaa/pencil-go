package main

import (
	consul "github.com/alaa/pencil-go/consul2"
	docker "github.com/alaa/pencil-go/docker2"
	"log"
	"time"
)

func main() {
	log.Println("starting pencil ...\n")

	for range time.Tick(5 * time.Second) {
		containers, err := docker.All()
		if err != nil {
			log.Printf("Unable to get list of all containers: %v\n", err)
			continue
		}

		err = consul.Resync(containers)
		if err != nil {
			log.Printf("Unable to resync services: %v\n", err)
			continue
		}
	}
}
