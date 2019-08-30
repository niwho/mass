package proto

type IMember interface {
	GetName() string
	GetTermId() int64
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

	// 同时会注册结点（服务）
	StartService() error

	GetMembers() []IMember
	GetLocal() IMember
	//Set(key string, val interface{})
	// 返回路由的IMember ,r如果没有则本地做，并广播其他已知节点 这个状态，其他节点收到要clear本地local记录
	GetMember(routerKey string) IMember
	UpateLocalRoute(routerKey string, member IMember)

	// sync 随机一个节点同步？
	SyncRoute(member IMember)

	// 广播 异步
	BroadCast(rpcname string, args interface{}, reply interface{}) error // 异步

	// 广播 同步
	CallAll(rpcname string, args interface{}, reply interface{}) ([]interface{},error) // 同步，有一个“正确”响应即可

}