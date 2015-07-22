package consul

import (
	"github.com/alaa/pencil-go/registry"
	consul "github.com/hashicorp/consul/api"
)

// ServiceRepository is consul-based implementation of registry.ServiceRepository
type ServiceRepository struct {
	consulAgent consulAgent
}

type consulAgent interface {
	Services() (map[string]*consul.AgentService, error)
	ServiceRegister(service *consul.AgentServiceRegistration) error
	ServiceDeregister(serviceID string) error
}

// NewServiceRepository creates new instance of ServiceRepository structure
func NewServiceRepository(consulAgent consulAgent) *ServiceRepository {
	return &ServiceRepository{consulAgent}
}

// Register adds service into consul
func (c *ServiceRepository) Register(service *registry.Service) error {
	return c.consulAgent.ServiceRegister(
		c.serviceToAgentServiceRegistration(service),
	)
}

// Unregister removes service from consul
func (c *ServiceRepository) Unregister(serviceID string) error {
	return c.consulAgent.ServiceDeregister(serviceID)
}

// GetAllIds return array of services ids registered in consul
func (c *ServiceRepository) GetAllIds() []string {
	services, _ := c.consulAgent.Services()
	servicesIDs := []string{}
	for _, service := range services {
		servicesIDs = append(servicesIDs, service.ID)
	}
	return servicesIDs
}

func (c *ServiceRepository) serviceToAgentServiceRegistration(service *registry.Service) *consul.AgentServiceRegistration {
	return &consul.AgentServiceRegistration{
		ID:   service.ID,
		Name: service.Service,
		Port: service.Port,
	}
}
