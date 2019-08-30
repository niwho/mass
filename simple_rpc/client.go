package simple_rpc

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"net/rpc"
	"time"
)

func Call(srv string, port int, rpcname string, args interface{}, reply interface{}) error {
	// 短连接调用方式
	conn, err := net.DialTimeout("tcp", srv+fmt.Sprintf(":%d", port), time.Second*10)
	if err != nil {
		return fmt.Errorf("ConnectError: %s", err.Error())
	}
	defer conn.Close() //
	readAndWriteTimeout := 2 * time.Second
	err = conn.SetDeadline(time.Now().Add(readAndWriteTimeout))
	if err != nil {
		return err
	}

	encBuf := bufio.NewWriter(conn)
	codec := &gobClientCodec{conn, gob.NewDecoder(conn), gob.NewEncoder(encBuf), encBuf}
	c := rpc.NewClientWithCodec(codec)
	err = c.Call(rpcname, args, reply)
	errc := c.Close()
	if err != nil && errc != nil {
		return fmt.Errorf("%s %s", err, errc)
	}
	if err != nil {
		return err
	} else {
		return errc
	}
}

func Pub(srv string, port int, rpcname string, args interface{}, reply interface{}) error {
	// 短连接调用方式
	conn, err := net.DialTimeout("tcp", srv+fmt.Sprintf(":%d", port), time.Second*10)
	if err != nil {
		return fmt.Errorf("ConnectError: %s", err.Error())
	}
	defer conn.Close() //
	readAndWriteTimeout := 2 * time.Second
	err = conn.SetDeadline(time.Now().Add(readAndWriteTimeout))
	if err != nil {
		return err
	}
	encBuf := bufio.NewWriter(conn)
	codec := &gobClientCodec{conn, gob.NewDecoder(conn), gob.NewEncoder(encBuf), encBuf}
	c := rpc.NewClientWithCodec(codec)
	c.Go(rpcname, args, reply, nil)
	errc := c.Close()
	if err != nil && errc != nil {
		return fmt.Errorf("%s %s", err, errc)
	}
	if err != nil {
		return err
	} else {
		return errc
	}
}
