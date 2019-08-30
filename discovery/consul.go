package discovery

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/niwho/logs"
	"github.com/niwho/mass/discovery/proto"
	"github.com/niwho/utils"
	"log"
)

var client *consulapi.Client

func init() {
	config := consulapi.DefaultConfig()
	config.Address = "consul.rightpaddle.cn:8500"
	clientx, err := consulapi.NewClient(config)
	if err != nil {
		logs.Log(logs.F{"err": err}).Error("init consul")
		panic(err)
	}
	client = clientx
}

func NewRegistration(name string, tags []string, meta map[string]string, advt string, port int) proto.IRegister {

	reg := &Registration{
		client: client,
	}
	reg.Register(name, tags, meta, advt, port)

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

func (reg *Registration) Register(name string, tags []string, meta map[string]string, advt string, port int) {
	reg.Id = utils.GenerateLogID()
	reg.ServiceName = name
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
	err := client.Agent().ServiceRegister(registration)

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
			ServiceName: ei.Service.Service,
			Meta:        ei.Service.Meta,
			Advt:        ei.Service.Address,
			Port:        ei.Service.Port,
		})
	}
	return iss, nil
}

func (reg *Registration) Unregister() {
	err := client.Agent().ServiceDeregister(reg.Id)
	if err != nil {
		//log.Fatal("register server error : ", err)
		logs.Log(logs.F{"err": err}).Error("Unregister")
	}
}

func (reg *Registration) GetName() string {
	return reg.ServiceName
}

func (reg *Registration) GetMeta(key string) string {
	val, _ := reg.Meta[key]
	return val
}

func (reg *Registration) GetAddr() string {
	return reg.Advt
}

func GetService(serviceName string) ([]proto.IService, error) {
	e, _, err := client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		logs.Log(logs.F{"err": err}).Error("GetService")
		return nil, err
	}
	var iss []proto.IService
	for _, ei := range e {
		log.Println(ei.Service.Address, ei.Service.Port, ei.Service.Meta, ei.Service.Service, ei.Service.ID,
			ei.Node.Address, ei.Node.ID)
		iss = append(iss, &Registration{
			ServiceName: ei.Service.Service,
			Meta:        ei.Service.Meta,
			Advt:        ei.Service.Address,
			Port:        ei.Service.Port,
		})
	}
	return iss, nil
}
