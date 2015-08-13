package consul

import (
	"github.com/brainly/pencil-go/registry"
	consul "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sort"
	"testing"
)

func TestThatRegisterCallConsulApiRegister(t *testing.T) {
	consulAgent := new(MockConsulAgent)
	consulServiceRepository := NewServiceRepository(consulAgent)

	consulAgent.On("ServiceRegister", &consul.AgentServiceRegistration{
		ID:   "redis1",
		Name: "redis",
		Port: 8000,
	}).Return(nil)

	err := consulServiceRepository.Register(&registry.Service{
		ID:      "redis1",
		Service: "redis",
		Port:    8000,
	})

	assert.Nil(t, err)
	consulAgent.AssertExpectations(t)
}

func TestThatDeregisterCallConsulApiDeregister(t *testing.T) {
	consulAgent := new(MockConsulAgent)
	consulServiceRepository := NewServiceRepository(consulAgent)

	consulAgent.On("ServiceDeregister", "redis1").Return(nil)

	err := consulServiceRepository.Deregister("redis1")
	assert.Nil(t, err)
	consulAgent.AssertExpectations(t)
}

type MockConsulAgent struct {
	mock.Mock
}

func TestThatGetAllIdsReturnArrayOfServicesIds(t *testing.T) {
	consulAgent := new(MockConsulAgent)
	consulServiceRepository := NewServiceRepository(consulAgent)

	consulAgent.On("Services").Return(map[string]*consul.AgentService{
		"redis": &consul.AgentService{
			ID:      "redis",
			Service: "redis",
			Address: "",
			Port:    8000,
		},
		"memcached": &consul.AgentService{
			ID:      "memcached",
			Service: "memcached",
			Address: "",
			Port:    9000,
		},
	}, nil)

	expectedArrayIds := []string{"memcached", "redis"}
	servicesIds := consulServiceRepository.GetAllIds()
	sort.Strings(servicesIds)
	assert.Equal(t, expectedArrayIds, servicesIds)

	consulAgent.AssertExpectations(t)
}

func (mca *MockConsulAgent) Services() (map[string]*consul.AgentService, error) {
	args := mca.Called()
	return args.Get(0).(map[string]*consul.AgentService), args.Error(1)
}

func (mca *MockConsulAgent) ServiceRegister(service *consul.AgentServiceRegistration) error {
	args := mca.Called(service)
	return args.Error(0)
}

func (mca *MockConsulAgent) ServiceDeregister(serviceID string) error {
	args := mca.Called(serviceID)
	return args.Error(0)
}
