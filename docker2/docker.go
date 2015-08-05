package docker2

import (
	"fmt"
	"github.com/alaa/pencil-go/container"
	docker "github.com/fsouza/go-dockerclient"
	"os"
	"strings"
)

var (
	endpoint          = "unix:///var/run/docker.sock"
	client, clientErr = docker.NewClient(endpoint)
)

// All returns all current running containers
func All() ([]container.Container, error) {
	dockerContainers, err := client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		return nil, err
	}

	containers := []container.Container{}
	for _, c := range dockerContainers {
		tags, err := containerTags(c)
		if err != nil {
			return nil, err
		}

		for _, port := range exposedPorts(c) {
			containers = append(containers, container.Container{
				ID:   c.ID,
				Name: withoutOrgName(c.Image),
				Port: port,
				Tags: tags,
			})
		}
	}

	return containers, nil
}

func containerTags(c docker.APIContainers) ([]string, error) {
	container, err := client.InspectContainer(c.ID)
	if err != nil {
		return nil, err
	}

	rawTags, found := container.Config.Labels["tags"]
	if !found || rawTags == "" {
		return []string{}, nil
	}

	return strings.Split(rawTags, ","), nil
}

func withoutOrgName(name string) string {
	parts := strings.Split(name, "/")
	if len(parts) == 1 {
		return name
	}
	return parts[1]
}

func exposedPorts(c docker.APIContainers) []int64 {
	ports := []int64{}
	for _, p := range c.Ports {
		if p.PublicPort > 0 {
			ports = append(ports, p.PublicPort)
		}
	}
	return ports
}

func init() {
	if clientErr != nil {
		fmt.Printf("Unable to connect to docker at %s: %v\n", endpoint, endpoint)
		os.Exit(1)
	}
}
