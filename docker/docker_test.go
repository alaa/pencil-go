package docker

import (
	docker "github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
	"testing"
)

var containers = []docker.APIContainers{
	docker.APIContainers{
		ID:      "bd1d34c0ebeeb62dfdcc57327aca15d2ef3cbc39a60e44aecb7085a8d1f89fd9",
		Image:   "33e91d7bac3e",
		Command: "/usr/sbin/sshd -D -o UseDNS=no -o UsePAM=no -o PasswordAuthentication=yes -o UsePrivilegeSeparation=no -o PidFile=/tmp/sshd.pid",
		Created: 1436541119,
		Status:  "Up 7 days",
		//Ports:      []docker.APIPort{PrivatePort: 22, PublicPort: 32782, Type: "tcp", IP: "0.0.0.0"},
		SizeRw:     0,
		SizeRootFs: 0,
		Names:      []string{"/elated_kirch"},
	},
}

// 2015/07/17 18:06:05 {bd1d34c0ebee   0 0 0  false false false [] map[22/tcp:{}] false false false [PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin DEBIAN_FRONTEND=noninteractive] [/usr/sbin/sshd -D -o UseDNS=no -o UsePAM=no -o PasswordAuthentication=yes -o UsePrivilegeSeparation=no -o PidFile=/tmp/sshd.pid] [] 33e91d7bac3e map[]   [] false [] [] map[]}

var containerConfig = docker.Config{
	Hostname:     "bd1d34c0ebee",
	ExposedPorts: map[docker.Port]struct{}{"22/tcp": {}},
	Env:          []string{"SRV_NAME=microservice1", "SRV_TAG=tag1"},
	Image:        "33e91d7bac3e",
}

var containerDetails = docker.Container{
	ID:              "bd1d34c0ebeeb62dfdcc57327aca15d2ef3cbc39a60e44aecb7085a8d1f89fd9",
	Config:          &containerConfig,
	Image:           "brainly/eve-landing-pages",
	NetworkSettings: nil,
	Name:            "/elated_kirch",
}

var containerNetworkSettings = docker.NetworkSettings{
	IPAddress:   "172.17.1.225",
	IPPrefixLen: 16,
	Gateway:     "172.17.42.1",
	Bridge:      "",
	//PortMapping map[string]PortMapping `json:"PortMapping,omitempty" yaml:"PortMapping,omitempty"`
	//Ports       map[Port][]PortBinding,
}

func TestGetRunningContainers(t *testing.T) {
	client := fakeDockerClient{}

	containers, err := GetRunningContainers(client)
	assert.Nil(t, err)

	assert.Equal(t, len(containers), 1)
	assert.Equal(t, containers[0].Name, "/elated_kirch")
}

type fakeDockerClient struct{}

func (c fakeDockerClient) listContainers(opts docker.ListContainersOptions) ([]docker.APIContainers, error) {
	return containers, nil
}

func (c fakeDockerClient) inspectContainer(id string) (*docker.Container, error) {
	return &containerDetails, nil
}
