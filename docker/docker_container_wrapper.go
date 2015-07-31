package docker

import (
	docker "github.com/fsouza/go-dockerclient"
	"sort"
	"strconv"
	"strings"
)

type dockerContainerWrapper struct {
	docker.Container
}

func (c *dockerContainerWrapper) getExposedTCPPorts() (ports []int) {
	for port := range c.NetworkSettings.Ports {
		if port.Proto() == "tcp" {
			port, _ := strconv.Atoi(port.Port())
			ports = append(ports, port)
		}
	}
	sort.Ints(ports)
	return
}

func (c *dockerContainerWrapper) getTags() []string {
	tags, exist := c.Config.Labels["tags"]
	if !exist {
		return []string{}
	}
	return strings.Split(tags, ",")
}

func (c *dockerContainerWrapper) getEnv() map[string]string {
	envMap := make(map[string]string)
	for _, value := range c.Config.Env {
		envParts := strings.Split(value, "=")
		envMap[envParts[0]] = envParts[1]
	}
	return envMap
}

func (c *dockerContainerWrapper) getName() string {
	if name, exist := c.getEnv()["SRV_NAME"]; exist {
		return name
	}
	return c.getImage()
}

func (c *dockerContainerWrapper) getImage() string {
	imageParts := strings.Split(c.Config.Image, "/")
	if len(imageParts) > 1 {
		return imageParts[1]
	}
	return imageParts[0]
}
