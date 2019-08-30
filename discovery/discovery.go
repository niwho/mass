package discovery

import (
	"fmt"
	"github.com/niwho/mass/discovery/proto"
)

type Service struct {
	Name string
	Meta map[string]string
	Addr string
}

func (sv *Service) GetName() string {
	return sv.Name
}

func (sv *Service) GetTags() []string {
	return nil
}

func (sv *Service) GetMeta() map[string]string {
	return sv.Meta
}

func (sv *Service) GetMetaValue(key string) string {
	val, _ := sv.Meta[key]
	return val
}

func (sv *Service) GetAddr() string {
	return sv.Addr
}

//
type Discovery struct {
	proto.IService
	proto.IRegister
	LocalIp string
	Port    int
}

func NewDiscovery(name string, meta map[string]string, localIp string, port int) proto.IDiscovery {
	dc := &Discovery{
		LocalIp: localIp,
		Port:    port,
	}

	dc.IService = &Service{
		Name: name,
		Meta: meta,
		Addr: fmt.Sprintf("%s:%d", localIp, port),
	}
	dc.IRegister = NewRegistration(dc.GetName(), dc.GetTags(), dc.GetMeta(), localIp, port)

	return dc
}

func (dc *Discovery) RegisterService(service proto.IService) {
	dc.Register(dc.GetName(), dc.GetTags(), dc.GetMeta(), dc.LocalIp, dc.Port)
}

func (dc *Discovery) GetService() ([]proto.IService , error){
	return dc.IRegister.GetService()
}
