package docker

import (
	"errors"
	"github.com/alaa/pencil-go/registry"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	client := mockDockerClient{}
	containerRepository := NewContainerRepository(&client)

	client.On("ListContainers", docker.ListContainersOptions{}).Return([]docker.APIContainers{}, nil)

	expectedContainers := []registry.Container{}

	containers, err := containerRepository.GetAll()

	assert.Nil(t, err)
	assert.Equal(t, expectedContainers, containers)
}

func TestGetRunningContainersWithTwoContainers(t *testing.T) {
	client := mockDockerClient{}
	repository := NewContainerRepository(&client)

	client.On("ListContainers", docker.ListContainersOptions{}).Return([]docker.APIContainers{containerA, containerB}, nil)
	client.On("InspectContainer", "bd1d34c0ebeeb62dfdcc57327aca15d2ef3cbc39a60e44aecb7085a8d1f89fd9").Return(&containerADetails, nil)
	client.On("InspectContainer", "f717f795bcccd674628b92f77a72f4b80b2c6b5da289846a0edbd21fb4c462db").Return(&containerBDetails, nil)

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

	containers, _ := repository.GetAll()
	sort.Sort(byID{containers})

	assert.Equal(t, expectedContainers, containers)
}

func TestGetAllWhenListContainersFails(t *testing.T) {
	client := mockDockerClient{}
	containerRepository := NewContainerRepository(&client)
	expectedError := errors.New("foo")

	client.On("ListContainers", docker.ListContainersOptions{}).Return([]docker.APIContainers{}, expectedError)

	_, err := containerRepository.GetAll()
	assert.Equal(t, expectedError, err)
}

func TestGetAllWhenInspectContainerFails(t *testing.T) {
	client := mockDockerClient{}
	containerRepository := NewContainerRepository(&client)
	expectedError := errors.New("bar")

	client.On("ListContainers", docker.ListContainersOptions{}).Return([]docker.APIContainers{containerA}, nil)
	client.On("InspectContainer", "bd1d34c0ebeeb62dfdcc57327aca15d2ef3cbc39a60e44aecb7085a8d1f89fd9").Return(&docker.Container{}, expectedError)

	_, err := containerRepository.GetAll()
	assert.Equal(t, expectedError, err)
}

type mockDockerClient struct {
	mock.Mock
}

func (c *mockDockerClient) ListContainers(opts docker.ListContainersOptions) ([]docker.APIContainers, error) {
	args := c.Called(opts)
	return args.Get(0).([]docker.APIContainers), args.Error(1)
}

func (c *mockDockerClient) InspectContainer(id string) (*docker.Container, error) {
	args := c.Called(id)
	return args.Get(0).(*docker.Container), args.Error(1)
}
