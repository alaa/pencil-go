package registry

// Registry understands how to synchronize registered Services with running Containers
type Registry struct {
	containerRepository ContainerRepository
	serviceRepository   ServiceRepository
}

// NewRegistry creates new instance of Registry
func NewRegistry(containerRepository ContainerRepository, serviceRepository ServiceRepository) *Registry {
	return &Registry{
		containerRepository,
		serviceRepository,
	}
}

// Synchronize synchronizes registered services according to running containers
func (r *Registry) Synchronize() error {
	registeredServicesIDs := r.serviceRepository.GetAllIds()
	runningContainers, err := r.containerRepository.GetAll()

	if err != nil {
		return err
	}

	r.registerServices(registeredServicesIDs, runningContainers)
	r.deregisterServices(registeredServicesIDs, runningContainers)

	return nil
}

func (r *Registry) registerServices(registeredServicesIDs []string, runningContainers []Container) {
	for _, service := range r.servicesToRegister(registeredServicesIDs, runningContainers) {
		r.serviceRepository.Register(service)
	}
}

func (r *Registry) deregisterServices(registeredServicesIDs []string, runningContainers []Container) {
	for _, serviceID := range r.servicesIDsToDeregister(registeredServicesIDs, runningContainers) {
		r.serviceRepository.Deregister(serviceID)
	}
}

func (r *Registry) servicesToRegister(registeredServicesIDs []string, runningContainers []Container) []*Service {
	servicesToRegister := []*Service{}
	registeredServicesIDsMap := r.sliceToMap(registeredServicesIDs)
	for _, container := range runningContainers {
		if _, ok := registeredServicesIDsMap[container.ID]; !ok {
			servicesToRegister = append(servicesToRegister, containerToService(&container))
		}
	}
	return servicesToRegister
}

func (r *Registry) servicesIDsToDeregister(registeredServicesIDs []string, runningContainers []Container) []string {
	servicesIdsToDeregister := []string{}
	runningContainersIDsSet := r.containersIDsMap(runningContainers)
	for _, serviceID := range registeredServicesIDs {
		if _, ok := runningContainersIDsSet[serviceID]; !ok {
			servicesIdsToDeregister = append(servicesIdsToDeregister, serviceID)
		}
	}
	return servicesIdsToDeregister
}

func containerToService(container *Container) *Service {
	return &Service{ID: container.ID, Service: container.Name, Port: container.Port, Tags: container.Tags}
}

func (r *Registry) sliceToMap(slice []string) map[string]bool {
	result := map[string]bool{}
	for _, item := range slice {
		result[item] = true
	}
	return result
}

func (r *Registry) containersIDsMap(containers []Container) map[string]bool {
	result := map[string]bool{}
	for _, container := range containers {
		result[container.ID] = true
	}
	return result
}
