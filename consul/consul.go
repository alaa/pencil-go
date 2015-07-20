package consul

import (
	consul "github.com/hashicorp/consul/api"
)

type ConsulClient struct {
	client *consul.Client
}

func NewConsulClient() (ConsulClient, error) {
	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		return ConsulClient{}, err
	}
	return ConsulClient{client: client}, nil
}

type ConsulAgent struct {
	agent *consul.Agent
}

func (c *ConsulClient) NewConsulAgent() ConsulAgent {
	agent := c.client.Agent()
	return ConsulAgent{agent: agent}
}

type Member struct {
	Name string
	IP   string
	Port uint16
}

type Members []Member

func buildMember(name string, ip string, port uint16) Member {
	return Member{Name: name, IP: ip, Port: port}
}

func (a *ConsulAgent) members() Members {
	list := Members{}
	use_wan := false
	members, _ := a.agent.Members(use_wan)
	for _, member := range members {
		list = append(list, buildMember(member.Name, member.Addr, member.Port))
	}
	return list
}

func buildService(id string, name string, port int, ip string) consul.AgentServiceRegistration {
	return consul.AgentServiceRegistration{ID: id, Name: name, Port: port, Address: ip}
}

func (a *ConsulAgent) Services() (map[string]*consul.AgentService, error) {
	services, err := a.agent.Services()
	if err != nil {
		return services, err
	}
	return services, nil
}

func (a *ConsulAgent) ServicesIDs() ([]string, error) {
	services, err := a.agent.Services()
	if err != nil {
		return []string{""}, err
	}

	list := []string{}
	for _, srv := range services {
		list = append(list, srv.ID)
	}

	return list, nil
}

func (a *ConsulAgent) RegisterService(id string, name string, port int, ip string) error {
	srv := buildService(id, name, port, ip)
	if err := a.agent.ServiceRegister(&srv); err != nil {
		return err
	}
	return nil
}

func (a *ConsulAgent) DeregisterService(id string) error {
	if err := a.agent.ServiceDeregister(id); err != nil {
		return err
	}
	return nil
}

func (a *ConsulAgent) DeregisterAllServices() error {
	services, _ := a.Services()
	for service := range services {
		a.DeregisterService(service)
	}
	return nil
}
