package main

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
)

type ConsulClient struct {
	client *consulapi.Client
}

func NewConsulClient() (ConsulClient, error) {
	client, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		return ConsulClient{}, err
	}
	return ConsulClient{client: client}, nil
}

type ConsulAgent struct {
	agent *consulapi.Agent
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
	stack := Members{}
	use_wan := false
	members, _ := r.agent.Members(use_wan)
	for _, member := range members {
		stack = append(stack, buildMember(member.Name, member.Addr, member.Port))
	}
	return stack
}

func main() {
	client, _ := NewConsulClient()
	agent := client.NewConsulAgent()
	fmt.Println(agent.members())
}
