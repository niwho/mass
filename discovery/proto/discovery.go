package proto

type IRegister interface {
	Register(name, service string, tags []string, meta map[string]string, advt string, port int)
	GetService() ([]IService, error)
	Unregister()
}

type IService interface {
	GetName() string
	GetServiceName() string
	GetTags() []string
	GetMetaValue(string) string
	GetMeta() map[string]string
	GetAddr() string
	GetHost() string
	GetPort() int
}

type IDiscovery interface {
	RegisterService(service IService)
	Unregister()
	GetService() ([]IService, error)
}
