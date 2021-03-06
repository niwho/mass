package proto

import "github.com/tidwall/buntdb"

type IMember interface {
	GetName() string
	GetTermId() int64
	UpdateTermId(termId int64)
	GetState() int
	SetInActive()
	SetActive()

	// 服务接口
	Call(rpcname string, args interface{}, reply interface{}) error
	Pub(rpcname string, args interface{}, reply interface{}) error

	Marshal() ([]byte, error)
}

// 路由信息
// 维护一颗B+数的映射，key->IMember
// 启动时从一个节点同步过来
type IMemberManager interface {
	GetRemovedKeysChan(routerKey string) chan string

	// 同时会注册结点（服务）
	RegisteRpc(m ...interface{}) error
	StartService() error
	Unregister()
	GetServiceXX(blocking bool)

	GetMembers() []IMember
	GetLocal() IMember
	GetLocalCopy() IMember
	//Set(key string, val interface{})
	// 返回路由的IMember ,r如果没有则本地做，并广播其他已知节点 这个状态，其他节点收到要clear本地local记录
	GetMember(routerKey string) IMember           // 只取当前节点 只是判断取 全局的 还是 local的
	GetMemberWithRemote(routerKey string) IMember // 先判断本节点，再获取远程节点，并同步, 如果都没有否则业务自行处理，新建或some
	GetMemberWithTry(routerKey string, retryCount int) (wantMember IMember)

	UpateLocalRoute(routerKey string, member IMember) // 更新本地路由
	RemoveLocalRoute(routerKey string) error
	RemoveUpdateLocalRoute(routerKey string, member IMember) error

	BroadCastRoute(routerKey string, member IMember) // 发布key在哪个节点， 至于发布策略，业务自己搞 一般第一次广播就行了
	BroadCastDelRoute(routerKey string) error
	BroadCastDelUpdateRoute(routerKey string, member IMember) error //

	// 高度集成的接口，包括了查找及update router
	// 查找节点，找不到则自己处理
	CallWithSync(routerKey, rpcName string, req, resp interface{}, localHandle func(found bool, req, resp interface{}) error) error
	CallWithOnceSync(routerKey, rpcName string, req, resp interface{}, localHandle func(found bool, req, resp interface{}) error) error
	CallWithUpdateSync(routerKey, rpcName string, req, resp interface{}, localHandle func(found bool, req, resp interface{}) error) error
	CallWithDelSync(routerKey, rpcName string, req, resp interface{}, localHandle func(req, resp interface{}) error) error

	// 以下这些接口 暂时没有真正的策略逻辑
	// sync 随机一个节点同步？
	SyncRoute(member IMember)

	// 广播 异步
	BroadCast(rpcname string, args interface{}, reply interface{}) error // 异步

	// 广播 同步
	CallAll(rpcname string, args interface{}, reply interface{}) ([]interface{}, error) // 同步，有一个“正确”响应即可

	// debug router
	GetRouter() *buntdb.DB
}

// req， resp 必须 实现的接口
type ISynMessage interface {
	GetKeyNode() IMember
	GetKey() string
	GetSourceNode() IMember
}
