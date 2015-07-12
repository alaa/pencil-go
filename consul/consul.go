package consul

import (
	consul "github.com/hashicorp/consul/api"
	"log"
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

func (a *ConsulAgent) registerService(id string, name string, port int, ip string) error {
	srv := buildService(id, name, port, ip)
	if err := a.agent.ServiceRegister(&srv); err != nil {
		return err
	}
	return nil
}

func (a *ConsulAgent) deregisterService(id string) error {
	if err := a.agent.ServiceDeregister(id); err != nil {
		return err
	}
	return nil
}

func (a *ConsulAgent) services() (map[string]*consul.AgentService, error) {
	services, err := a.agent.Services()
	if err != nil {
		return services, err
	}
	return services, nil
}

func (a *ConsulAgent) deregisterServiceID(service_id string) error {
	log.Printf("deregistering service: %s \n", service_id)
	if err := a.deregisterService(service_id); err != nil {
		return err
	}
	return nil
}

func (a *ConsulAgent) deregisterAllServices() error {
	services, _ := a.services()
	for service := range services {
		a.deregisterServiceID(service)
	}
	return nil
}

// TODO
// Chain the creation of client().agent()
// func main() {
// 	client, err := NewConsulClient()
// 	if err != nil {
// 		log.Fatal("Could not connect to consul client")
// 	}
//
// 	agent := client.NewConsulAgent()
//
// 	fmt.Println(agent.members())
//
// 	agent.registerService("cid1", "srv-1", 1234, "127.0.0.1")
// 	agent.registerService("cid2", "srv-2", 2345, "127.0.0.1")
// 	agent.registerService("cid3", "srv-3", 3456, "127.0.0.1")
//
// 	time.Sleep(20 * time.Second)
//
// 	agent.deregisterAllServices()
// }
