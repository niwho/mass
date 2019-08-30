package simple_rpc

import (
	"github.com/niwho/mass/simple_rpc/proto"
	"net/rpc"
)

type SimpleRpc struct {
	Host string
	Port int
}

func NewSimpleRpc(host string, port int) proto.ISimpleRpc {
	return &SimpleRpc{
		Host: host,
		Port: port,
	}
}

func (sp *SimpleRpc) RegisterRpc(rcvrs ...interface{}) error {
	var err error
	for _, rcvr := range rcvrs {
		if err = rpc.Register(rcvr); err != nil {
			return err
		}
	}
	return nil
}

func (sp *SimpleRpc) ListenRPC() error {
	return ListenRPC(sp.Port)
}

func (sp *SimpleRpc) Call(rpcname string, args interface{}, reply interface{}) error {
	return Call(sp.Host, sp.Port, rpcname, args, reply)
}

func (sp *SimpleRpc) Pub(rpcname string, args interface{}, reply interface{}) error {
	return Pub(sp.Host, sp.Port, rpcname, args, reply)
}
