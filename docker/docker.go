package docker

import (
	docker "github.com/fsouza/go-dockerclient"
	"strings"
	"time"
)

const (
	Endpoint = "unix:///var/run/docker.sock"
	Interval = 5 * time.Second
)

const SRV_NAME = "SRV_NAME"

type ConcreteDockerClient struct {
	client *docker.Client
}

type DockerClient interface {
	listContainers(opts docker.ListContainersOptions) ([]docker.APIContainers, error)
	inspectContainer(id string) (*docker.Container, error)
}

type TCPPorts []string

type UDPPorts []string

type ContainerInfo struct {
	ID          string
	Name        string
	ImageID     string
	ImageName   string
	ServiceName string
	Env         map[string]string
	Created     time.Time
	Config      *docker.Config
	TCPPorts    TCPPorts
	UDPPorts    UDPPorts
}

type EnvVariables map[string]string

// NewDockerClient creates an instance of DockerClient.
// Returns error when unable to connect to docker daemon.
func NewDockerClient() (DockerClient, error) {
	client, err := docker.NewClient(Endpoint)
	if err != nil {
		return &ConcreteDockerClient{}, err
	}
	return &ConcreteDockerClient{client: client}, nil
}

// GetRunningContainers finds running containers and returns specific details.
// TOOD package should return error, logging should be disabled.
func GetRunningContainers(c DockerClient) ([]ContainerInfo, error) {
	containersIDs, err := getContainersIDs(c)
	if err != nil {
		return []ContainerInfo{}, err
	}
	containersDetails, err := getContainersDetails(c, containersIDs)
	if err != nil {
		return []ContainerInfo{}, err
	}
	return containersDetails, nil
}

func (c *ConcreteDockerClient) listContainers(opts docker.ListContainersOptions) ([]docker.APIContainers, error) {
	return c.client.ListContainers(opts)
}

func (c *ConcreteDockerClient) inspectContainer(id string) (*docker.Container, error) {
	return c.client.InspectContainer(id)
}

func buildContainerInfo(container *docker.Container) ContainerInfo {
	tcpPorts, udpPorts := getExposedPorts(container)
	imageName := getImageName(container.Config.Image)
	envVars := getEnvVariables(container.Config.Env)
	srvName := serviceName(container)

	return ContainerInfo{
		ID:          container.ID,
		Name:        container.Name,
		ImageID:     container.Image,
		ImageName:   imageName,
		ServiceName: srvName,
		Env:         envVars,
		Created:     container.Created,
		Config:      container.Config,
		TCPPorts:    tcpPorts,
		UDPPorts:    udpPorts,
	}
}

func getImageName(image string) string {
	if img := strings.Split(image, "/"); len(img) > 1 {
		return img[1]
	} else {
		return img[0]
	}
}

func serviceName(container *docker.Container) string {
	serviceName := getImageName(container.Config.Image)
	serviceEnv := getEnvVariables(container.Config.Env)

	if value, exsists := serviceEnv[SRV_NAME]; exsists {
		serviceName = value
	}
	return serviceName
}

func getExposedPorts(container *docker.Container) (TCPPorts, UDPPorts) {
	tcp_list := TCPPorts{}
	udp_list := UDPPorts{}
	exposed_ports := container.Config.ExposedPorts

	for port, _ := range exposed_ports {
		if port.Proto() == "tcp" {
			tcp_list = append(tcp_list, port.Port())
		} else if port.Proto() == "udp" {
			udp_list = append(udp_list, port.Port())
		}
	}
	return tcp_list, udp_list
}

func getEnvVariables(env []string) EnvVariables {
	m := make(map[string]string)
	for _, value := range env {
		e := strings.Split(value, "=")
		m[e[0]] = e[1]
	}
	return m
}

func getContainersIDs(c DockerClient) ([]docker.APIContainers, error) {
	options := docker.ListContainersOptions{}
	containers, err := c.listContainers(options)
	if err != nil {
		return containers, err
	}
	return containers, nil
}

func getContainersDetails(c DockerClient, containers []docker.APIContainers) ([]ContainerInfo, error) {
	list := []ContainerInfo{}
	for _, container := range containers {
		i, err := inspectContainer(c, container.ID)
		if err != nil {
			return []ContainerInfo{}, err
		}
		list = append(list, i)
	}
	return list, nil
}

func inspectContainer(c DockerClient, cid string) (ContainerInfo, error) {
	data, err := c.inspectContainer(cid)
	if err != nil {
		return ContainerInfo{}, err
	}
	return buildContainerInfo(data), nil
}
