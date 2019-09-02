package discovery

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/niwho/logs"
	"github.com/niwho/mass/discovery/proto"
	"github.com/niwho/utils"
	"log"
	"sync"
)

const (
	DEFAULT_CONSUL = "consul.rightpaddle.cn:8500"
)

var onetime sync.Once

func Init(consulAddress string) *consulapi.Client {
	config := consulapi.DefaultConfig()
	config.Address = consulAddress
	if consulAddress == "" {
		config.Address = DEFAULT_CONSUL
	}
	clientx, err := consulapi.NewClient(config)
	if err != nil {
		logs.Log(logs.F{"err": err}).Error("init consul")
		panic(err)
	}
	return clientx
}

func NewRegistration(name, serviceName string, tags []string, meta map[string]string, advt string, port int, consulAddress string) proto.IRegister {
	client := Init(consulAddress)
	reg := &Registration{
		client: client,

		ServiceName: name,
		Advt:        advt,
		Port:        port,
		Meta:        meta,
		Tags:        tags,
	}
	reg.Register(name, serviceName, tags, meta, advt, port)

	return reg
}

type Registration struct {
	Id          string
	ServiceName string
	Advt        string
	Port        int

	Meta map[string]string
	Tags []string

	client *consulapi.Client
}

func (reg *Registration) Register(nodeName, serviceName string, tags []string, meta map[string]string, advt string, port int) {
	reg.Id = nodeName
	reg.ServiceName = serviceName
	reg.Tags = tags
	reg.Meta = meta

	if advt != "" {
		reg.Advt = advt
	} else {
		reg.Advt = utils.GetLocalIP()
	}

	reg.Port = port

	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = reg.Id
	registration.Name = reg.ServiceName
	registration.Port = reg.Port
	registration.Tags = reg.Tags
	registration.Meta = reg.Meta

	registration.Address = reg.Advt

	//增加check。
	check := new(consulapi.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("http://%s:%d%s", registration.Address, registration.Port, "/health/check/")
	//设置超时 5s。
	check.Timeout = "5s"
	//设置间隔 5s。
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "30s"
	//注册check服务。
	registration.Check = check
	//log.Println("get check.HTTP:", check)
	logs.Log(logs.F{"check": check}).Debug("")
	err := reg.client.Agent().ServiceRegister(registration)

	if err != nil {
		log.Fatal("register server error : ", err)
		logs.Log(logs.F{"err": err}).Error("Register")
	}
}

func (reg *Registration) GetService() ([]proto.IService, error) {
	e, _, err := reg.client.Health().Service(reg.ServiceName, "", true, nil)
	if err != nil {
		logs.Log(logs.F{"err": err}).Error("GetService")
		return nil, err
	}
	var iss []proto.IService
	for _, ei := range e {
		log.Println(ei.Service.Address, ei.Service.Port, ei.Service.Meta, ei.Service.Service, ei.Service.ID,
			ei.Node.Address, ei.Node.ID)
		iss = append(iss, &Registration{
			Id:          ei.Node.ID,
			ServiceName: ei.Service.Service,
			Meta:        ei.Service.Meta,
			Advt:        ei.Service.Address,
			Port:        ei.Service.Port,
		})
	}
	logs.Log(logs.F{"iss": iss}).Debug("GetService")
	return iss, nil
}

func (reg *Registration) Unregister() {
	err := reg.client.Agent().ServiceDeregister(reg.Id)
	if err != nil {
		//log.Fatal("register server error : ", err)
		logs.Log(logs.F{"err": err}).Error("Unregister")
	}
}

func (reg *Registration) GetServiceName() string {
	return reg.ServiceName
}
func (reg *Registration) GetName() string {
	return reg.Id
}

func (reg *Registration) GetAddressWithPort() string {
	return fmt.Sprintf("%s:%d", reg.Advt, reg.Port)
}

func (reg *Registration) GetMetaValue(key string) string {
	val, _ := reg.Meta[key]
	return val
}

func (reg *Registration) GetMeta() map[string]string {
	return reg.Meta
}

func (reg *Registration) GetTags() []string {
	return reg.Tags
}

func (reg *Registration) GetAddr() string {
	return reg.Advt
}

func (reg *Registration) GetHost() string {
	return reg.Advt
}

func (reg *Registration) GetPort() int {
	return reg.Port
}
