package registry

// ContainerRepository is responsible for keeping Containers
type ContainerRepository interface {
	GetAll() []*Container
}

// ServiceRepository is responsible for keeping Services
type ServiceRepository interface {
	GetAllIds() []string
	Register(service *Service) error
	Deregister(serviceID string) error
}

// Container entity
type Container struct {
	ID   string
	Name string
	Port int
}

// Service entity
type Service struct {
	ID      string
	Service string
	Tags    []string
	Address string
	Port    int
	Check   ServiceCheck
}

// ServiceCheck describes details of service health check
type ServiceCheck struct {
	Script   string
	HTTP     string
	Interval string
	TTL      string
}
