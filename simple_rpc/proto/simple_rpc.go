package proto

type ISimpleRpc interface {
	RegisterRpc(rcvrs ...interface{}) error
	ListenRPC() error

	Call(rpcname string, args interface{}, reply interface{}) error
	Pub(rpcname string, args interface{}, reply interface{}) error
}
