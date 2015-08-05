package consul

import (
	"fmt"
	"github.com/alaa/pencil-go/container"
	consul "github.com/hashicorp/consul/api"
	"os"
	"strings"
)

var (
	client, clientErr = consul.NewClient(consul.DefaultConfig())
)

func Resync(containers []container.Container) error {
	services, err := currentServiceIDs()
	if err != nil {
		return err
	}

	toDeregister := servicesToDeregister(services, containers)

	for _, s := range toDeregister {
		err = client.Agent().ServiceDeregister(s)
		if err != nil {
			return err
		}
	}

	toRegister := containersToRegister(services, containers)

	for _, c := range toRegister {
		id := serviceID(c)
		err = client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
			ID:   id,
			Name: id,
			Tags: c.Tags,
			Port: int(c.Port),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func currentServiceIDs() (map[string]*consul.AgentService, error) {
	services, err := client.Agent().Services()
	if err != nil {
		return nil, err
	}

	delete(services, "consul")
	return services, nil
}

func containersToRegister(running map[string]*consul.AgentService, containers []container.Container) []container.Container {
	result := []container.Container{}
	for _, c := range containers {
		if _, exists := running[serviceID(c)]; !exists {
			result = append(result, c)
		}
	}
	return result
}

func servicesToDeregister(running map[string]*consul.AgentService, containers []container.Container) []string {
	containerMap := map[string]struct{}{}
	for _, c := range containers {
		containerMap[serviceID(c)] = struct{}{}
	}

	result := []string{}
	for service, _ := range running {
		if _, exists := containerMap[service]; !exists {
			result = append(result, service)
		}
	}
	return result
}

func serviceID(c container.Container) string {
	return fmt.Sprintf("%s-%s-%d", c.ID, withoutTag(c.Name), c.Port)
}

func withoutTag(name string) string {
	return strings.Split(name, ":")[0]
}

func init() {
	if clientErr != nil {
		fmt.Printf("Unable to connect to consul: %v\n", clientErr)
		os.Exit(1)
	}
}
