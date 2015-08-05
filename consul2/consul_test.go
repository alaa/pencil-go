package consul2

import (
	"github.com/alaa/pencil-go/container"
	consul "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	consulService = consul.AgentService{
		ID:      "consul",
		Service: "consul",
		Tags:    []string{},
		Port:    8300,
	}
)

func TestResyncEmpty(t *testing.T) {
	containers := []container.Container{}
	err := Resync(containers)
	assert.Nil(t, err)

	services, err := client.Agent().Services()
	assert.Nil(t, err)

	assertSameServices(t, map[string]consul.AgentService{
		"consul": consulService,
	}, services)
}

func TestResyncAddContainers(t *testing.T) {
	TestResyncEmpty(t)

	containers := []container.Container{
		container.Container{
			ID:   "57jkhgmv67",
			Name: "helloworld",
			Port: 3394,
			Tags: []string{"tag1", "tag2"},
		},
		container.Container{
			ID:   "57jkhgmv67",
			Name: "helloworld",
			Port: 22,
			Tags: []string{"tag1", "tag2"},
		},
		container.Container{
			ID:   "99nameg643b",
			Name: "some-test:latest",
			Port: 9931,
			Tags: []string{},
		},
	}
	err := Resync(containers)
	assert.Nil(t, err)

	services, err := client.Agent().Services()
	assert.Nil(t, err)

	assertSameServices(t, map[string]consul.AgentService{
		"consul": consulService,
		"57jkhgmv67-helloworld-3394": consul.AgentService{
			ID:      "57jkhgmv67-helloworld-3394",
			Service: "57jkhgmv67-helloworld-3394",
			Port:    3394,
			Tags:    []string{"tag1", "tag2"},
		},
		"57jkhgmv67-helloworld-22": consul.AgentService{
			ID:      "57jkhgmv67-helloworld-22",
			Service: "57jkhgmv67-helloworld-22",
			Port:    22,
			Tags:    []string{"tag1", "tag2"},
		},
		"99nameg643b-some-test-9931": consul.AgentService{
			ID:      "99nameg643b-some-test-9931",
			Service: "99nameg643b-some-test-9931",
			Port:    9931,
			Tags:    nil,
		},
	}, services)
}

func TestResyncRemoveAndAddContainers(t *testing.T) {
	TestResyncAddContainers(t)

	containers := []container.Container{
		container.Container{
			ID:   "57jkhgmv67",
			Name: "helloworld",
			Port: 22,
			Tags: []string{"tag1", "tag2"},
		},
		container.Container{
			ID:   "771hhjasfu34hu",
			Name: "landing-pages",
			Port: 9999,
			Tags: []string{"tag3", "tag2"},
		},
	}
	err := Resync(containers)
	assert.Nil(t, err)

	services, err := client.Agent().Services()
	assert.Nil(t, err)

	assertSameServices(t, map[string]consul.AgentService{
		"consul": consulService,
		"57jkhgmv67-helloworld-22": consul.AgentService{
			ID:      "57jkhgmv67-helloworld-22",
			Service: "57jkhgmv67-helloworld-22",
			Port:    22,
			Tags:    []string{"tag1", "tag2"},
		},
		"771hhjasfu34hu-landing-pages-9999": consul.AgentService{
			ID:      "771hhjasfu34hu-landing-pages-9999",
			Service: "771hhjasfu34hu-landing-pages-9999",
			Port:    9999,
			Tags:    []string{"tag3", "tag2"},
		},
	}, services)
}

func assertSameServices(t *testing.T, expected map[string]consul.AgentService, actual map[string]*consul.AgentService) {
	for key, service := range expected {
		if actual[key] == nil {
			t.Errorf("Expected to see service %s\n", key)
			continue
		}

		assert.Equal(t, service.ID, actual[key].ID)
		assert.Equal(t, service.Service, actual[key].Service)
		assert.Equal(t, service.Port, actual[key].Port)
		assert.Equal(t, service.Tags, actual[key].Tags)
	}

	for key, _ := range actual {
		service, found := expected[key]
		if !found {
			t.Errorf("Expected not to see service %s\n", key)
			continue
		}

		assert.Equal(t, service.ID, actual[key].ID)
		assert.Equal(t, service.Service, actual[key].Service)
		assert.Equal(t, service.Port, actual[key].Port)
		assert.Equal(t, service.Tags, actual[key].Tags)
	}
}
