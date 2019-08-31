package member_manager

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/niwho/logs"
	"github.com/niwho/mass/discovery"
	dc_proto "github.com/niwho/mass/discovery/proto"
	"github.com/niwho/mass/member_manager/proto"
	"github.com/niwho/mass/simple_rpc"
	rpc_proto "github.com/niwho/mass/simple_rpc/proto"
	"github.com/tidwall/buntdb"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type MemberSub struct {
	Name string `json:"name"`

	host   string `json:"host"`
	port   int    `json:"port"`
	TermId int64  `json:"term_id"`

	State      int   `json:"state"`
	UpdateTime int64 `json:"update_time"`
}

type Member struct {
	rpc_proto.ISimpleRpc `json:"-"`
	MemberSub
}

func NewMember(name string, host string, port int) proto.IMember {
	return &Member{
		ISimpleRpc: simple_rpc.NewSimpleRpc(host, port),
		MemberSub: MemberSub{
			Name:   name,
			host:   host,
			port:   port,
			State:  1,
			TermId: time.Now().Unix(),
		},
	}
}

func (m *Member) GetName() string {
	return m.Name
}

func (m *Member) GetState() int {
	return m.State
}

func (m *Member) GetTermId() int64 {
	return m.TermId
}

func (m *Member) SetInActive() {
	m.State = 0
}

func (m *Member) SetActive() {
	m.State = 1
}

func (m *Member) Call(rpcname string, args interface{}, reply interface{}) error {
	return m.ISimpleRpc.Call(rpcname, args, reply)
}

func (m *Member) Pub(rpcname string, args interface{}, reply interface{}) error {
	return m.ISimpleRpc.Pub(rpcname, args, reply)
}

func (m *Member) Marshal() ([]byte, error) {
	return jsonFastest.Marshal(m.MemberSub)
}

type MemberManager struct {
	ServiceName string

	LocalIp string
	Port    int
	//rpc_proto.ISimpleRpc
	localMember proto.IMember
	dc_proto.IDiscovery

	// 管理所有结点（状态）
	members   sync.Map
	routeInfo *buntdb.DB

	stoped bool
}

func NewMemberManager(localName, ServiceName string, localIp string, port int, meta map[string]string, consuleAddress string) (proto.IMemberManager, error) {
	imm := &MemberManager{
		ServiceName: ServiceName,
		LocalIp:     localIp,
		Port:        port,
		localMember: NewMember(localName, localIp, port),
	}
	/*
		local, _ := buntdb.Open(":memory:")
		local.Update(func(tx *buntdb.Tx) error {
			tx.CreateIndex("index_uid", "*", buntdb.IndexJSON("uid"))
			tx.CreateIndex("index_did", "*", buntdb.IndexJSON("did"))
			return nil
		})
	*/

	var err error
	imm.routeInfo, err = buntdb.Open(":memory:")
	imm.routeInfo.Update(func(tx *buntdb.Tx) error {
		tx.CreateIndex("index_name", "*", buntdb.IndexJSON("name"))
		return nil
	})

	if err != nil {
		return nil, err
	}
	//imm.ISimpleRpc = simple_rpc.NewSimpleRpc(localIp, port)
	imm.IDiscovery = discovery.NewDiscovery(ServiceName, meta, localIp, port, consuleAddress)

	imm.localMember.(*Member).RegisterRpc(&MemberSync{
		local: imm.localMember.(*Member),
	})

	// 同步节点
	go imm.backgroud()

	return imm, nil
}

func (mm *MemberManager) StartService() error {
	//err := mm.ISimpleRpc.ListenRPC()
	err := mm.localMember.(*Member).ListenRPC()
	return err
}

// 使用discovery获取
func (mm *MemberManager) GetMembers() []proto.IMember {
	var members []proto.IMember
	mm.members.Range(func(key, value interface{}) bool {
		members = append(members, value.(proto.IMember))
		return true
	})
	return members
}

func (mm *MemberManager) GetLocal() proto.IMember {
	return mm.localMember
}

func (mm *MemberManager) Set(key string, val interface{}) {
	panic("implement me")
}

// 查找local路由信息
// room_id -> {}
// key -> {Member的字段}
//
func (mm *MemberManager) GetMember(routerKey string) proto.IMember {
	var memsub MemberSub
	err := mm.routeInfo.View(func(tx *buntdb.Tx) error {
		if val, err := tx.Get(routerKey); err != nil {
			return jsonFastest.Unmarshal([]byte(val), &memsub)
		} else {
			return err
		}
	})
	if err != nil {
		logs.Log(logs.F{"err": err}).Error("GetMember")
		return nil
	}

	return &Member{MemberSub: memsub}
}

// 本地没有则探测
func (mm *MemberManager) GetMemberMust(routerKey string) proto.IMember {
	var memsub MemberSub
	err := mm.routeInfo.View(func(tx *buntdb.Tx) error {
		if val, err := tx.Get(routerKey); err != nil {
			return jsonFastest.Unmarshal([]byte(val), &memsub)
		} else {
			return err
		}
	})
	if err != nil {
		logs.Log(logs.F{"err": err}).Error("GetMember")
		// probe 探测
		var resp SyncResponse
		mm.CallAll("MemberSync.Probe", SyncRequest{Key: routerKey}, &resp)
		return nil
	}

	return &Member{MemberSub: memsub}
}

// 根据termid更新
func (mm *MemberManager) UpateLocalRoute(routerKey string, member proto.IMember) {
	err := mm.routeInfo.Update(func(tx *buntdb.Tx) error {

		if val, err := tx.Get(routerKey); err == nil {
			OldTermId := jsoniter.Get([]byte(val), "term_id").ToInt64()
			if member.GetTermId() > OldTermId {
				if val, err := member.Marshal(); err == nil {
					tx.Set(routerKey, string(val), nil)
				}
			} else {
				// 这种情况是不是触发一次广播当前的最新的routeKey ->member的映射，防止别的节点也有旧的数据映射
				// 或反馈 请求来源方
			}
		} else {
			if val, err := member.Marshal(); err == nil {
				tx.Set(routerKey, string(val), nil)
			} else {
				return err
			}
		}
		return nil
	})

	if err != nil {
		logs.Log(logs.F{"err": err}).Error("UpateLocalRoute")
	}

}

//和另一个结点交互路由信息
func (mm *MemberManager) SyncRoute(member proto.IMember) {
	//
	panic("implement me")
}

// 广播
func (mm *MemberManager) BroadCast(rpcname string, args interface{}, reply interface{}) error {
	for _, member := range mm.GetMembers() {
		// 是否显示过滤自己local
		if member.GetName() == mm.localMember.GetName() {
			continue
		}

		// 并发pub todo
		member.Call(rpcname, args, reply)
	}
	return nil
}

// 广播
// 只要有一个正常响应就行
// 这个逻辑是不是放到业务层去做更合适
// 可充当协调员的角色
func (mm *MemberManager) CallAll(rpcname string, args interface{}, reply interface{}) (callReplys []interface{}, err error) {
	// 有一个响应的即可
	var wg sync.WaitGroup

	for _, member := range mm.GetMembers() {
		// 是否显示过滤自己local
		wg.Add(1)
		// 并发pub
		go func() {
			innerReply := reflect.New(reflect.ValueOf(reply).Elem().Type()) //Elem 必须是指针类型
			member.Call(rpcname, args, innerReply)
			wg.Done()

			//innerReply 这次调用的返回结果
			// 怎么确定是正常成功的呢
			callReplys = append(callReplys, innerReply) // 和reply是同类型的集合
		}()
	}
	return
}

func (mm *MemberManager) routerInfoRemoveNode(nodeName string) error {
	var needDeleteKeys []string
	_ = mm.routeInfo.View(func(tx *buntdb.Tx) error {
		return tx.AscendEqual("index_name", fmt.Sprintf(`{"name":"%s"}`, nodeName), func(key, value string) bool {
			needDeleteKeys = append(needDeleteKeys, key)
			return true // 一直继续
		})
	})

	if len(needDeleteKeys) > 0 {
		_ = mm.routeInfo.Update(func(tx *buntdb.Tx) error {
			for _, key := range needDeleteKeys {
				tx.Delete(key)
			}
			return nil
		})
	}
	return nil
}

// 轮询节点状态变化信息
func (mm *MemberManager) backgroud() error {
	job := func() {
		defer func() {
			//atomic.AddInt32(&af.idleNum, -1)
			if err := recover(); err != nil {
				const size = 64 << 20
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				fmt.Printf("AsyncJob panic=%v\n%s\n", err, buf)
			}
		}()

		for !mm.stoped {
			updateTime := time.Now().Unix()
			// 更新节点
			nodes, _ := mm.GetService()
			for _, node := range nodes {
				// 相互通信的端口信息在meta里，最外层的port是服务的端口（check health）
				meta := node.GetMeta()
				if meta == nil {
					logs.Log(logs.F{"node": node}).Error("")
					continue
				}
				host := meta["c_host"]
				port, _ := strconv.ParseInt(meta["c_port"], 10, 64)
				localNode, ok := mm.members.Load(node.GetName())
				if !ok {
					localNode, _ = mm.members.LoadOrStore(node.GetName(), NewMember(node.GetName(), host, int(port)))
				}
				localNode.(*Member).UpdateTime = updateTime
			}

			// 如果是 初始化启动，是不是随机从几个节点获取信息填充自己的routerInfo todo

			for _, mem := range mm.GetMembers() {
				if mem.(*Member).UpdateTime < updateTime {
					// 该节点已经失效了
					mm.members.Delete(mem.GetName())
					//清理routerinfo
					mm.routerInfoRemoveNode(mem.GetName())
				}
			}

			time.Sleep(time.Second * 5) //阻塞
		}
	}

	for !mm.stoped {
		job() //阻塞
	}

	return nil
}
