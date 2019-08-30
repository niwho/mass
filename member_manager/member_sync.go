package member_manager

import "github.com/niwho/mass/member_manager/proto"

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

type ProxyRequest struct {
}

type ProxyResponse struct {
}

type SyncRequest struct {
	Node MemberSub
	Key  string
}

type SyncResponse struct {
	ErrorCode int
	Node      MemberSub
	Key       string
}

// 代理请求
//func (ms *MemberSync) ProxyCall(req ProxyRequest, resp *ProxyResponse) {
//
//}

// 同步或确认key的最新“分布”
func (ms *MemberSync) Probe(req SyncRequest, resp *SyncResponse) error{
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
	ms.manager.UpateLocalRoute(req.Key, &Member{MemberSub: req.Node})

	resp.ErrorCode = 0
	return nil
}
