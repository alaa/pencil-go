package docker

import (
	docker "github.com/fsouza/go-dockerclient"
	"log"
	"strings"
	"time"
)

const (
	Endpoint = "unix:///var/run/docker.sock"
	Interval = 5 * time.Second
)

type DockerClient struct {
	client *docker.Client
}

func NewDocker() DockerClient {
	client, err := docker.NewClient(Endpoint)
	if err != nil {
		log.Fatal(err)
	}
	return DockerClient{client: client}
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

const SRV_NAME = "SRV_NAME"

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

type EnvVariables map[string]string

func getEnvVariables(env []string) EnvVariables {
	m := make(map[string]string)
	for _, value := range env {
		e := strings.Split(value, "=")
		m[e[0]] = e[1]
	}
	return m
}

// getRunningContainers finds running containers and returns specific details.
// TOOD package should return error, logging should be disabled.
func (c *DockerClient) GetRunningContainers() []ContainerInfo {
	containersIDs, err := c.getContainersIDs()
	if err != nil {
		log.Print(err)
	}
	containersDetails := c.getContainersDetails(containersIDs)
	log.Print("Running containers count: ", len(containersDetails), "\n\n")
	return containersDetails
}

// getContainersIDs retruns a list of running docker contianers.
func (c *DockerClient) getContainersIDs() ([]docker.APIContainers, error) {
	options := docker.ListContainersOptions{}
	containers, err := c.client.ListContainers(options)
	if err != nil {
		return containers, err
	}
	return containers, nil
}

// getContainersDetails iterate over a list of containers and returns a list of ContainerInfo struct.
func (c *DockerClient) getContainersDetails(containers []docker.APIContainers) []ContainerInfo {
	list := []ContainerInfo{}
	for _, container := range containers {
		i, err := c.inspectContainer(container.ID)
		if err != nil {
			log.Print(err)
		}
		list = append(list, i)
	}
	return list
}

// inspectContainer extract container info for a continer ID.
func (c *DockerClient) inspectContainer(cid string) (ContainerInfo, error) {
	data, err := c.client.InspectContainer(cid)
	log.Printf("%v", *data.NetworkSettings)
	if err != nil {
		return ContainerInfo{}, err
	}
	return buildContainerInfo(data), nil
}
