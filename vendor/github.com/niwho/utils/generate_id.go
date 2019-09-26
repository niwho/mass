package utils

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/satori/go.uuid"
)

var _LOCALIP string

func GetLocalIP() string {
	if _LOCALIP != "" {
		return _LOCALIP
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "0.0.0.0"
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				_LOCALIP = ipnet.IP.String()
				return _LOCALIP
			}
		}
	}
	return "0.0.0.0"
}

func GetLocalIPWithPrefix(prefix string) string {
	var localIP string
	if localIP  != "" {
		return localIP
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "0.0.0.0"
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIP = ipnet.IP.String()
				if strings.HasPrefix(localIP, prefix){
					return localIP
				}

			}
		}
	}
	return "0.0.0.0"
}

func GenerateLogID() string {
	iparr := strings.Split(GetLocalIP(), ".")
	fip := ""
	for _, v := range iparr {
		fip += fmt.Sprintf("%03s", v)
	}

	dt := time.Now().Format("20060102150405")

	u1, _ := uuid.NewV4()
	urand := fmt.Sprintf("%d", binary.BigEndian.Uint64(u1.Bytes()))[:6]
	return dt + fip + urand
}

func FormatLogID(logid string) string {
	if len(logid) < 32 {
		return logid
	}
	return fmt.Sprintf("%s-%s-%s-%s-%s", logid[:8], logid[8:12], logid[12:16], logid[16:20], logid[20:])
}
