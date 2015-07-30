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
func (r *ServiceRepository) Register(service *registry.Service) error {
	return r.consulAgent.ServiceRegister(
		buildAgentServiceRegistration(service),
	)
}

// Deregister removes service from consul
func (r *ServiceRepository) Deregister(serviceID string) error {
	return r.consulAgent.ServiceDeregister(serviceID)
}

// GetAllIds return array of services ids registered in consul
func (r *ServiceRepository) GetAllIds() []string {
	services, _ := r.consulAgent.Services()
	servicesIDs := []string{}
	for _, service := range services {
		servicesIDs = append(servicesIDs, service.ID)
	}
	return servicesIDs
}

func buildAgentServiceRegistration(service *registry.Service) *consul.AgentServiceRegistration {
	return &consul.AgentServiceRegistration{
		ID:   service.ID,
		Name: service.Service,
		Port: service.Port,
		Tags: service.Tags,
	}
}
