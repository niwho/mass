package simple_rpc

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"net/rpc"
)


func Regiseter(rcvrs ...interface{}) error {
	var err error
	for _, rcvr := range rcvrs {
		if err = rpc.Register(rcvr); err != nil {
			return err
		}
	}
	return nil
}

func ListenRPC(port int) error{
	l, e := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if e != nil {
		log.Print("Error: listen error:", e)
		return e
	}
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Print("Error: accept rpc connection", err.Error())
				continue
			}
			// 可以限制并发数，待优化
			go func(conn net.Conn) {
				buf := bufio.NewWriter(conn)
				srv := &gobServerCodec{
					rwc:    conn,
					dec:    gob.NewDecoder(conn),
					enc:    gob.NewEncoder(buf),
					encBuf: buf,
				}
				err = rpc.ServeRequest(srv)
				if err != nil {
					log.Print("Error: server rpc request", err.Error())
				}
				srv.Close()
			}(conn)
		}
	}()
	return nil
}
