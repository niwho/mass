package proto

type IRegister interface {
	Register(name string, tags []string, meta map[string]string, advt string, port int)
	GetService() ([]IService, error)
	Unregister()
}

type IService interface {
	GetName() string
	GetTags() []string
	GetMetaValue(string) string
	GetMeta() map[string]string
	GetAddr() string
}

type IDiscovery interface {
	RegisterService(service IService)
	GetService() ([]IService, error)
}
