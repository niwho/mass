package utils

import (
	"fmt"
	"net"
	"time"
)

var conn net.Conn

func InitHeartBeat(isMaster bool, port int, cb func(val string), addr string) {
	conn, _ = net.DialTimeout("udp", addr, time.Second*2)
	if isMaster {
		for {
			time.Sleep(time.Second)
			SendHeartBeat()
		}
	} else {
		go func() {
			UDPServer(port, cb)
		}()
	}
}

func UDPServer(port int, cb func(string)) {
	pc, err := net.ListenPacket("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("UDPServer", err)
	}
	buffer := make([]byte, 10240)
	for {
		n, addr, err := pc.ReadFrom(buffer)
		_ = addr
		if err == nil {
			if n > 0 {
				cb(string(buffer[:n]))
			}
		}
	}
}

func SendHeartBeat() {
	conn.Write([]byte("doki"))
}
