package docker

import (
	docker "github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	// container A section
	containerA = docker.APIContainers{
		ID:         "bd1d34c0ebeeb62dfdcc57327aca15d2ef3cbc39a60e44aecb7085a8d1f89fd9",
		Image:      "33e91d7bac3e",
		Command:    "/usr/sbin/sshd -D -o UseDNS=no -o UsePAM=no -o PasswordAuthentication=yes -o UsePrivilegeSeparation=no -o PidFile=/tmp/sshd.pid",
		Created:    1436541119,
		Status:     "Up 7 days",
		SizeRw:     0,
		SizeRootFs: 0,
		Names:      []string{"/elated_kirch"},
	}

	containerAConfig = docker.Config{
		Hostname:     "bd1d34c0ebee",
		ExposedPorts: map[docker.Port]struct{}{"22/tcp": {}},
		Env:          []string{"SRV_NAME=microservice1", "SRV_TAG=tag1"},
		Image:        "brainly/eve-landing-pages",
	}

	containerADetails = docker.Container{
		ID:              "bd1d34c0ebeeb62dfdcc57327aca15d2ef3cbc39a60e44aecb7085a8d1f89fd9",
		Config:          &containerAConfig,
		Image:           "33e91d7bac3e",
		NetworkSettings: nil,
		Name:            "/elated_kirch",
		Created:         time.Unix(1436541119, 0),
	}

	// container B section
	containerB = docker.APIContainers{
		ID:         "f717f795bcccd674628b92f77a72f4b80b2c6b5da289846a0edbd21fb4c462db",
		Image:      "e5d4de01ea02",
		Command:    "/usr/local/bin/wiyd --daddy=John",
		Created:    1436541319,
		Status:     "Up 9 days",
		SizeRw:     0,
		SizeRootFs: 0,
		Names:      []string{"/naughty_heisenberg"},
	}

	containerBConfig = docker.Config{
		Hostname:     "f717f795bccc",
		ExposedPorts: map[docker.Port]struct{}{"9000/tcp": {}},
		Env:          []string{"SRV_NAME=microservice2", "SRV_TAG=tag2"},
		Image:        "brainly/eve-who-is-your-daddy",
	}

	containerBDetails = docker.Container{
		ID:              "f717f795bcccd674628b92f77a72f4b80b2c6b5da289846a0edbd21fb4c462db",
		Config:          &containerBConfig,
		Image:           "e5d4de01ea02",
		NetworkSettings: nil,
		Name:            "/naughty_heisenberg",
		Created:         time.Unix(1436541319, 0),
	}
)

func TestGetRunningContainers(t *testing.T) {
	client := fakeDockerClient{
		fakeContainers: []docker.APIContainers{containerA},
		fakeContainerDetails: map[string]*docker.Container{
			"bd1d34c0ebeeb62dfdcc57327aca15d2ef3cbc39a60e44aecb7085a8d1f89fd9": &containerADetails,
		},
	}

	expectedContainers := []ContainerInfo{
		ContainerInfo{
			ID:          "bd1d34c0ebeeb62dfdcc57327aca15d2ef3cbc39a60e44aecb7085a8d1f89fd9",
			Name:        "/elated_kirch",
			ImageName:   "eve-landing-pages",
			ImageID:     "33e91d7bac3e",
			ServiceName: "microservice1",
			Env: map[string]string{
				"SRV_NAME": "microservice1",
				"SRV_TAG":  "tag1",
			},
			Created:  time.Unix(1436541119, 0),
			Config:   &containerAConfig,
			TCPPorts: TCPPorts{"22"},
			UDPPorts: UDPPorts{},
		},
	}

	containers, err := GetRunningContainers(client)
	assert.Nil(t, err)

	assert.Equal(t, expectedContainers, containers)
}

func TestGetRunningContainersWithNoContainers(t *testing.T) {
	client := fakeDockerClient{}

	expectedContainers := []ContainerInfo{}

	containers, err := GetRunningContainers(client)
	assert.Nil(t, err)

	assert.Equal(t, expectedContainers, containers)
}

func TestGetRunningContainersWithTwoContainers(t *testing.T) {
	client := fakeDockerClient{
		fakeContainers: []docker.APIContainers{
			containerA,
			containerB,
		},
		fakeContainerDetails: map[string]*docker.Container{
			"bd1d34c0ebeeb62dfdcc57327aca15d2ef3cbc39a60e44aecb7085a8d1f89fd9": &containerADetails,
			"f717f795bcccd674628b92f77a72f4b80b2c6b5da289846a0edbd21fb4c462db": &containerBDetails,
		},
	}

	expectedContainers := []ContainerInfo{
		ContainerInfo{
			ID:          "bd1d34c0ebeeb62dfdcc57327aca15d2ef3cbc39a60e44aecb7085a8d1f89fd9",
			Name:        "/elated_kirch",
			ImageName:   "eve-landing-pages",
			ImageID:     "33e91d7bac3e",
			ServiceName: "microservice1",
			Env: map[string]string{
				"SRV_NAME": "microservice1",
				"SRV_TAG":  "tag1",
			},
			Created:  time.Unix(1436541119, 0),
			Config:   &containerAConfig,
			TCPPorts: TCPPorts{"22"},
			UDPPorts: UDPPorts{},
		},
		ContainerInfo{
			ID:          "f717f795bcccd674628b92f77a72f4b80b2c6b5da289846a0edbd21fb4c462db",
			Name:        "/naughty_heisenberg",
			ImageName:   "eve-who-is-your-daddy",
			ImageID:     "e5d4de01ea02",
			ServiceName: "microservice2",
			Env: map[string]string{
				"SRV_NAME": "microservice2",
				"SRV_TAG":  "tag2",
			},
			Created:  time.Unix(1436541319, 0),
			Config:   &containerBConfig,
			TCPPorts: TCPPorts{"9000"},
			UDPPorts: UDPPorts{},
		},
	}

	containers, err := GetRunningContainers(client)
	assert.Nil(t, err)

	assert.Equal(t, expectedContainers, containers)
}

type fakeDockerClient struct {
	fakeContainers []docker.APIContainers
	fakeContainerDetails map[string]*docker.Container
}



func (c fakeDockerClient) listContainers(opts docker.ListContainersOptions) ([]docker.APIContainers, error) {
	return c.fakeContainers, nil
}

func (c fakeDockerClient) inspectContainer(id string) (*docker.Container, error) {
	return c.fakeContainerDetails[id], nil
}
