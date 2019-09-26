package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

var udpClientIns *UdpClient

type UdpClient struct {
	conn  net.Conn
	token string
	addr  string
}

type ApplicationMessage struct {
	Title   string   `json:"title"`
	MType   string   `json:"mtype"`   //dingidng, kafka, influx
	Content string   `json:"content"` // 具体数据信息，json格式
	At      []string `json:"at"`
	Token   string   `json:"token"`
}

func InitUdpClient(token string) {
	udpClientIns = &UdpClient{
		token: token,
		addr:  "127.0.0.1:5678",
	}
	udpClientIns.conn, _ = net.DialTimeout("udp", udpClientIns.addr, time.Second*60)
}

func SendDingDing(title, content string) error {
	am := ApplicationMessage{
		Title:   title,
		MType:   "dingding",
		Content: content,
		At:      []string{"17388935273"},
		Token:   udpClientIns.token,
	}
	b, _ := json.Marshal(am)
	_, err := udpClientIns.conn.Write(b)
	fmt.Println("send dingidng err:", err)
	return err
}
