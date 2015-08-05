package docker2

import (
	"github.com/alaa/pencil-go/container"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestAll(t *testing.T) {
	c1 := runContainer(t, _name("helloworld"), []string{}, []string{"49822:22", "49900:9000", "5555"})
	defer removeContainer(c1)

	c2 := runContainer(t, _name("some_test"), []string{"tag1", "tag2"}, []string{"49800:8000"})
	defer removeContainer(c2)

	containers, err := All()

	assert.Nil(t, err)

	assertIncludesContainer(t, containers, container.Container{
		ID:   c1.ID,
		Name: "helloworld",
		Port: 49822,
		Tags: []string{},
	})

	assertIncludesContainer(t, containers, container.Container{
		ID:   c1.ID,
		Name: "helloworld",
		Port: 49900,
		Tags: []string{},
	})

	assertDoesNotIncludeContainer(t, containers, container.Container{
		ID:   c1.ID,
		Name: "helloworld",
		Port: 5555,
		Tags: []string{},
	})

	assertIncludesContainer(t, containers, container.Container{
		ID:   c2.ID,
		Name: "some_test",
		Port: 49800,
		Tags: []string{"tag1", "tag2"},
	})

	removeContainer(c1)

	containers, err = All()

	assert.Nil(t, err)

	assertDoesNotIncludeContainer(t, containers, container.Container{
		ID:   c1.ID,
		Name: "helloworld",
		Port: 49822,
		Tags: []string{},
	})

	assertDoesNotIncludeContainer(t, containers, container.Container{
		ID:   c1.ID,
		Name: "helloworld",
		Port: 49900,
		Tags: []string{},
	})

	assertIncludesContainer(t, containers, container.Container{
		ID:   c2.ID,
		Name: "some_test",
		Port: 49800,
		Tags: []string{"tag1", "tag2"},
	})
}

func assertIncludesContainer(t *testing.T, containers []container.Container, expected container.Container) {
	if !includesContainer(containers, expected) {
		t.Errorf("\n  Expected %v\nto include %v\n", containers, expected)
	}
}

func assertDoesNotIncludeContainer(t *testing.T, containers []container.Container, expected container.Container) {
	if includesContainer(containers, expected) {
		t.Errorf("\n      Expected %v\nto not include %v\n", containers, expected)
	}
}

func includesContainer(containers []container.Container, expected container.Container) bool {
	for _, c := range containers {
		if c.IsEqual(expected) {
			return true
		}
	}
	return false
}

func runContainer(t *testing.T, name string, tags []string, ports []string) *docker.Container {
	client.RemoveContainer(docker.RemoveContainerOptions{
		ID:    name,
		Force: true,
	})

	client.TagImage("alpine:3.2", docker.TagImageOptions{
		Repo:  name,
		Tag:   "latest",
		Force: true,
	})

	hostConfig := &docker.HostConfig{
		PortBindings: portBindingsFrom(ports),
	}

	labels := map[string]string{
		"tags": strings.Join(tags, ","),
	}

	container, err := client.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:        name,
			Cmd:          []string{"sleep", "5"},
			Labels:       labels,
			PortSpecs:    ports,
			ExposedPorts: exposedPortsFrom(ports),
		},
		HostConfig: hostConfig,
	})

	assert.Nil(t, err)

	client.StartContainer(container.ID, hostConfig)
	return container
}

func removeContainer(c *docker.Container) {
	client.RemoveContainer(docker.RemoveContainerOptions{
		ID:    c.ID,
		Force: true,
	})
}

func exposedPortsFrom(ports []string) map[docker.Port]struct{} {
	result := make(map[docker.Port]struct{})
	withExposedPorts(ports, func(parts []string) {
		result[docker.Port(parts[1]+"/tcp")] = struct{}{}
	})
	return result
}

func portBindingsFrom(ports []string) map[docker.Port][]docker.PortBinding {
	result := make(map[docker.Port][]docker.PortBinding)
	withExposedPorts(ports, func(parts []string) {
		result[docker.Port(parts[1]+"/tcp")] = []docker.PortBinding{
			docker.PortBinding{HostPort: parts[0]},
		}
	})
	return result
}

func withExposedPorts(ports []string, fn func([]string)) {
	for _, p := range ports {
		parts := strings.Split(p, ":")
		if len(parts) != 2 {
			continue
		}
		fn(parts)
	}
}

func _name(name string) string {
	return "pencil_go_test/" + name
}
