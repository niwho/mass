package member_manager

import (
	"github.com/niwho/mass/member_manager/proto"
	"github.com/niwho/mass/simple_rpc"
)

type IProxyRequest interface {
}

type IProxyResponse interface {
}

type ISyncRequest interface {
}

type ISyncResponse interface {
}

type IMemberSync interface {
	//ProxyCall(IProxyRequest, IProxyResponse)
	SyncKey(ISyncRequest, ISyncResponse)
}

//
type MemberSync struct {
	local proto.IMember

	manager proto.IMemberManager
}

//type ProxyRequest struct {
//}
//
//type ProxyResponse struct {
//}

type SyncRequest struct {
	Node MemberSub
	Key  string
}

func (req SyncRequest) GetKeyNode() proto.IMember {
	return nil
}

func (req SyncRequest) GetKey() string {
	return req.Key
}

func (req SyncRequest) GetSourceNode() proto.IMember {
	return &Member{
		ISimpleRpc: simple_rpc.NewSimpleRpc(req.Node.Host, req.Node.Port),
		MemberSub:  req.Node,
	}
}

type SyncResponse struct {
	ErrorCode int
	Node      MemberSub
	SouceNode MemberSub
	Key       string
}

func (resp *SyncResponse) GetKeyNode() proto.IMember {
	if resp.Node.Name != "" {
		return &Member{
			ISimpleRpc: simple_rpc.NewSimpleRpc(resp.Node.Host, resp.Node.Port),
			MemberSub:  resp.Node,
		}
	}
	return nil
}

func (resp *SyncResponse) GetKey() string {
	return resp.Key
}

func (resp *SyncResponse) GetSourceNode() proto.IMember {
	return &Member{
		ISimpleRpc: simple_rpc.NewSimpleRpc(resp.SouceNode.Host, resp.SouceNode.Port),
		MemberSub:  resp.SouceNode,
	}
}

// 代理请求
//func (ms *MemberSync) ProxyCall(req ProxyRequest, resp *ProxyResponse) {
//
//}

// 同步或确认key的最新“分布”
func (ms *MemberSync) Probe(req SyncRequest, resp *SyncResponse) error {
	resp.SouceNode = ms.local.(*Member).MemberSub

	mem := ms.manager.GetMember(req.Key)
	if mem == nil {
		resp.ErrorCode = -1

		return nil
	}

	resp.ErrorCode = 0
	resp.Key = req.Key
	resp.Node = mem.(*Member).MemberSub
	return nil
}

// 同步或确认key的最新“分布”
func (ms *MemberSync) SyncKey(req SyncRequest, resp *SyncResponse) error {
	resp.SouceNode = ms.local.(*Member).MemberSub

	ms.manager.UpateLocalRoute(req.Key, &Member{MemberSub: req.Node})

	resp.ErrorCode = 0
	resp.Node = req.Node
	resp.Key = req.Key

	return nil
}

func (ms *MemberSync) DelKey(req SyncRequest, resp *SyncResponse) error {
	resp.SouceNode = ms.local.(*Member).MemberSub

	ms.manager.RemoveLocalRoute(req.Key)

	resp.ErrorCode = 0
	resp.Node = req.Node
	resp.Key = req.Key

	return nil
}
