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
func (r *Registry) Synchronize() {
	registeredServicesIDs := r.serviceRepository.GetAllIds()
	runningContainers := r.containerRepository.GetAll()

	r.registerServices(registeredServicesIDs, runningContainers)
	r.unregisterServices(registeredServicesIDs, runningContainers)
}

func (r *Registry) registerServices(registeredServicesIDs []string, runningContainers []*Container) {
	for _, service := range r.servicesToRegister(registeredServicesIDs, runningContainers) {
		r.serviceRepository.Register(service)
	}
}

func (r *Registry) unregisterServices(registeredServicesIDs []string, runningContainers []*Container) {
	for _, serviceID := range r.servicesIDsToUnregister(registeredServicesIDs, runningContainers) {
		r.serviceRepository.Unregister(serviceID)
	}
}

func (r *Registry) servicesToRegister(registeredServicesIDs []string, runningContainers []*Container) []*Service {
	servicesToRegister := []*Service{}
	registeredServicesIDsMap := r.sliceToMap(registeredServicesIDs)
	for _, container := range runningContainers {
		if _, ok := registeredServicesIDsMap[container.ID]; !ok {
			servicesToRegister = append(servicesToRegister, containerToService(container))
		}
	}
	return servicesToRegister
}

func (r *Registry) servicesIDsToUnregister(registeredServicesIDs []string, runningContainers []*Container) []string {
	servicesIdsToUnregister := []string{}
	runningContainersIDsSet := r.containersIDsMap(runningContainers)
	for _, serviceID := range registeredServicesIDs {
		if _, ok := runningContainersIDsSet[serviceID]; !ok {
			servicesIdsToUnregister = append(servicesIdsToUnregister, serviceID)
		}
	}
	return servicesIdsToUnregister
}

func containerToService(container *Container) *Service {
	return &Service{ID: container.ID, Service: container.Name, Port: container.Port}
}

func (r *Registry) sliceToMap(slice []string) map[string]bool {
	result := map[string]bool{}
	for _, item := range slice {
		result[item] = true
	}
	return result
}

func (r *Registry) containersIDsMap(containers []*Container) map[string]bool {
	result := map[string]bool{}
	for _, container := range containers {
		result[container.ID] = true
	}
	return result
}
