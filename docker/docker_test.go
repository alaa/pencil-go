package docker

import (
	"github.com/brainly/pencil-go/registry"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

var (
	// container A section
	containerA = docker.APIContainers{
		ID: "bd1d34c0ebeeb62dfdcc57327aca15d2ef3cbc39a60e44aecb7085a8d1f89fd9",
	}

	containerAConfig = docker.Config{
		Env:   []string{"SRV_TAG=tag1"},
		Image: "brainly/eve-landing-pages",
	}

	containerADetails = docker.Container{
		ID:     "bd1d34c0ebeeb62dfdcc57327aca15d2ef3cbc39a60e44aecb7085a8d1f89fd9",
		Config: &containerAConfig,
		NetworkSettings: &docker.NetworkSettings{
			Ports: map[docker.Port][]docker.PortBinding{
				"22/tcp":   []docker.PortBinding{},
				"8000/tcp": []docker.PortBinding{},
			},
		},
	}

	// container B section
	containerB = docker.APIContainers{
		ID: "f717f795bcccd674628b92f77a72f4b80b2c6b5da289846a0edbd21fb4c462db",
	}

	containerBConfig = docker.Config{
		Env:   []string{"SRV_NAME=microservice2", "SRV_TAG=tag2"},
		Image: "brainly/eve-who-is-your-daddy",
	}

	containerBDetails = docker.Container{
		ID:     "f717f795bcccd674628b92f77a72f4b80b2c6b5da289846a0edbd21fb4c462db",
		Config: &containerBConfig,
		NetworkSettings: &docker.NetworkSettings{
			Ports: map[docker.Port][]docker.PortBinding{
				"9000/tcp": []docker.PortBinding{},
			},
		},
	}
)

type containers []registry.Container
type byID struct{ containers }

func (s containers) Len() int      { return len(s) }
func (s containers) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byID) Less(i, j int) bool  { return s.containers[i].ID < s.containers[j].ID }

func TestProperlyWrapContainers(t *testing.T) {
	testWrapperForContainer(t, containerADetails, "eve-landing-pages", []int{22, 8000})
	testWrapperForContainer(t, containerBDetails, "microservice2", []int{9000})
}

func testWrapperForContainer(t *testing.T, container docker.Container, expectedName string, expectedPorts []int) {
	dockerContainerWrapper := dockerContainerWrapper{container}
	assert.Equal(t, expectedName, dockerContainerWrapper.getName())
	assert.Equal(t, expectedPorts, dockerContainerWrapper.getExposedTCPPorts())
}

func TestGetAllWhenNoContainersAreRunning(t *testing.T) {
	client := fakeDockerClient{}

	expectedContainers := []registry.Container{}

	containerRepository := NewContainerRepository(client)
	containers := containerRepository.GetAll()

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

	expectedContainers := []registry.Container{
		registry.Container{
			ID:   "bd1d34c0ebeeb62dfdcc57327aca15d2ef3cbc39a60e44aecb7085a8d1f89fd9",
			Name: "eve-landing-pages",
			Port: 22,
		},
		registry.Container{
			ID:   "bd1d34c0ebeeb62dfdcc57327aca15d2ef3cbc39a60e44aecb7085a8d1f89fd9",
			Name: "eve-landing-pages",
			Port: 8000,
		},
		registry.Container{
			ID:   "f717f795bcccd674628b92f77a72f4b80b2c6b5da289846a0edbd21fb4c462db",
			Name: "microservice2",
			Port: 9000,
		},
	}

	adapter := NewContainerRepository(client)
	containers := adapter.GetAll()
	sort.Sort(byID{containers})

	assert.Equal(t, expectedContainers, containers)
}

type fakeDockerClient struct {
	fakeContainers       []docker.APIContainers
	fakeContainerDetails map[string]*docker.Container
}

func (c fakeDockerClient) ListContainers(opts docker.ListContainersOptions) ([]docker.APIContainers, error) {
	return c.fakeContainers, nil
}

func (c fakeDockerClient) InspectContainer(id string) (*docker.Container, error) {
	return c.fakeContainerDetails[id], nil
}
