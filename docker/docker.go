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

type DockerClient interface {
	ListContainers(opts docker.ListContainersOptions) ([]docker.APIContainers, error)
	InspectContainer(id string) (*docker.Container, error)
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

type DockerAdapter struct {
	dockerClient DockerClient
}

func NewDockerAdapter(dockerClient DockerClient) *DockerAdapter {
	return &DockerAdapter{dockerClient: dockerClient}
}

// GetRunningContainers finds running containers and returns specific details.
// TOOD package should return error, logging should be disabled.
func (d *DockerAdapter) GetRunningContainers() ([]ContainerInfo, error) {
	containersIDs, err := d.getContainersIDs()
	if err != nil {
		return []ContainerInfo{}, err
	}
	containersDetails, err := d.getContainersDetails(containersIDs)
	if err != nil {
		return []ContainerInfo{}, err
	}
	return containersDetails, nil
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

func (d *DockerAdapter) getContainersIDs() ([]docker.APIContainers, error) {
	options := docker.ListContainersOptions{}
	return d.dockerClient.ListContainers(options)
}

func (d *DockerAdapter) getContainersDetails(containers []docker.APIContainers) ([]ContainerInfo, error) {
	list := []ContainerInfo{}
	for _, container := range containers {
		i, err := d.inspectContainer(container.ID)
		if err != nil {
			return []ContainerInfo{}, err
		}
		list = append(list, i)
	}
	return list, nil
}

func (d *DockerAdapter) inspectContainer(cid string) (ContainerInfo, error) {
	data, err := d.dockerClient.InspectContainer(cid)
	if err != nil {
		return ContainerInfo{}, err
	}
	return buildContainerInfo(data), nil
}
