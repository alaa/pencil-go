package main

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"time"
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

func (r *ConsulClient) NewConsulAgent() ConsulAgent {
	agent := r.client.Agent()
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

func (r *ConsulAgent) members() Members {
	list := Members{}
	use_wan := false
	members, _ := r.agent.Members(use_wan)
	for _, member := range members {
		list = append(list, buildMember(member.Name, member.Addr, member.Port))
	}
	return list
}

func buildService(id string, name string, port int, ip string) consul.AgentServiceRegistration {
	return consul.AgentServiceRegistration{ID: id, Name: name, Port: port, Address: ip}
}

func (r *ConsulAgent) registerService(id string, name string, port int, ip string) error {
	srv := buildService(id, name, port, ip)
	if err := r.agent.ServiceRegister(&srv); err != nil {
		return err
	}
	return nil
}

func (r *ConsulAgent) deregisterService(id string) error {
	if err := r.agent.ServiceDeregister(id); err != nil {
		return err
	}
	return nil
}

func (r *ConsulAgent) services() (map[string]*consul.AgentService, error) {
	services, err := r.agent.Services()
	if err != nil {
		return services, err
	}
	return services, nil
}

func main() {
	client, _ := NewConsulClient()
	agent := client.NewConsulAgent()

	fmt.Println(agent.members())
	fmt.Println(agent.services())

	agent.registerService("docker_id_here", "srv-search", 1234, "127.0.0.1")
	time.Sleep(20 * time.Second)

	agent.deregisterService("docker_id_here")
}
